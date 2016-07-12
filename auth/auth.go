package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	r "gopkg.in/dancannon/gorethink.v2"

	"errors"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"net/http"
	"strings"
	"time"
)

type AdminUser struct {
	ID       string `gorethink:"id"`
	Level    int    `gorethink:"level"`
	Username string `gorethink:"username"`
	Password string `gorethink:"password"`
}

// GinJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userId is made available as
// c.Get("userId").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX#!/usr/bin/env
type GinJWTMiddleware struct {
	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Rethink to use
	Rethink *r.Session
}

// MiddlewareInit initialize jwt configs.
func (mw *GinJWTMiddleware) MiddlewareInit() error {

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.Key == nil {
		return errors.New("secret key is required")
	}

	return nil
}

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *GinJWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			mw.unauthorized(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	return func(c *gin.Context) {
		mw.middlewareImpl(c)
		return
	}
}

func (mw *GinJWTMiddleware) middlewareImpl(c *gin.Context) {
	token, err, code := mw.parseToken(c)
	if err != nil {
		mw.unauthorized(c, code, err.Error())
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Set("JWT_PAYLOAD", claims)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	res, err := r.
		Table("admins").
		Get(claims["id"].(string)).
		Run(mw.Rethink)

	if err != nil {
		mw.unauthorized(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer res.Close()

	var user AdminUser
	err = res.One(&user)

	// TODO: Differentiate between "rows = 0" & "a real cursor error"
	if err != nil {
		mw.unauthorized(c, http.StatusForbidden, "You're not allowed to do this")
		return
	}

	c.Set("user", user)
	c.Next()
}

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context) {

	// Initial middleware default setting.
	mw.MiddlewareInit()

	var loginVals struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	if c.BindJSON(&loginVals) != nil {
		mw.unauthorized(c, http.StatusBadRequest, "Username/password missing")
		return
	}

	// Callback function that should perform the authentication of the user based on userId and
	// password. Must return true on success, false on failure. Required.
	// Option return user id, if so, user id will be stored in Claim Array.
	res, err := r.Table("admins").GetAllByIndex("username", loginVals.Username).Pluck("id", "password").Run(mw.Rethink)
	if err != nil {
		mw.unauthorized(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer res.Close()

	var doc map[string]string
	err = res.One(&doc)

	if (err != nil) || (bcrypt.CompareHashAndPassword([]byte(doc["password"]), []byte(loginVals.Password)) != nil) {
		mw.unauthorized(c, http.StatusUnauthorized, "Username/password mismatch")
		return
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	expiry := time.Now().Add(mw.Timeout)
	claims["id"] = doc["id"]
	claims["exp"] = expiry.Unix()
	claims["orig_iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(mw.Key)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Failed to create JWT token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"token":  tokenString,
			"expire": expiry.Format(time.RFC3339),
		},
	})
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the GinJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {
	oldToken, err, code := mw.parseToken(c)
	if err != nil {
		mw.unauthorized(c, code, err.Error())
		return
	}

	claims := oldToken.Claims.(jwt.MapClaims)

	// Checks if this token can still be refreshed
	origIat := int64(claims["orig_iat"].(float64))
	if origIat < time.Now().Add(-mw.MaxRefresh).Unix() {
		mw.unauthorized(c, http.StatusUnauthorized, "Token has expired")
		return
	}

	// Create a new token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims) // the claims for the new token

	// copy claims from old to new
	for key := range claims {
		newClaims[key] = claims[key]
	}

	// update the expiry time for this new token
	expiry := time.Now().Add(mw.Timeout)
	newClaims["exp"] = expiry.Unix()

	tokenString, err := newToken.SignedString(mw.Key)
	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Failed to create JWT token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"token":  tokenString,
			"expire": expiry.Format(time.RFC3339),
		},
	})
}

func (mw *GinJWTMiddleware) parseToken(c *gin.Context) (*jwt.Token, error, int) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization required"), http.StatusForbidden
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("invalid auth header"), http.StatusBadRequest
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		return mw.Key, nil
	})

	code := http.StatusOK
	if err != nil {
		code = http.StatusBadRequest
	}
	return token, err, code
}

func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm=jacr-api")
	c.Abort()

	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})

	return
}
