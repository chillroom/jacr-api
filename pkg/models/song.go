package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Song struct {
	ID         int        `db:"id"`
	Fkid       string     `db:"fkid"`
	Name       string     `db:"name"`
	LastPlay   time.Time  `db:"last_play"`
	SkipReason SkipReason `db:"skip_reason"`
	Song       int        `db:"song"`
	User       int        `db:"user"`
	Time       time.Time  `db:"time"`
}

type SkipReason string

const (
	ForbiddenSkip   SkipReason = "forbidden"
	NsfwSkip        SkipReason = "nsfw"
	ThemeSkip       SkipReason = "theme"
	UnavailableSkip SkipReason = "unavailable"
)

func (s SkipReason) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *SkipReason) Scan(src interface{}) error {
	//   out, ok := src.([]byte)
	//   if !ok {
	// 	  return errors.New("")
	//   }
	_, err := fmt.Sscanf(string(src.([]byte)), "%s", &s)

	// if err != nil {
	// 	if string(*s) == "" {
	// 		s = nil
	// 	}
	// }
	//   return SkipReason(out)
	return err
}

type SongType string

const (
	YouTubeSong    SongType = "youtube"
	SoundCloudSong SongType = "soundcloud"
)

func (s SongType) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *SongType) Scan(src interface{}) error {
	//   out, ok := src.([]byte)
	//   if !ok {
	// 	  return errors.New("")
	//   }
	_, err := fmt.Sscanf(string(src.([]byte)), "%s", s)

	//   return SkipReason(out)
	return err
}
