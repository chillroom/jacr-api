package models

type Notice struct {
	ID      int    `db:"id"`
	Message string `db:"message"`
	Title   string `db:"title"`
}
