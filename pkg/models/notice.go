package models

type Notice struct {
	ID      int    `db:"id" goqu:"skipinsert"`
	Message string `db:"message"`
	Title   string `db:"title"`
}
