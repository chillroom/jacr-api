package models

type User struct {
	ID int `db:"id" goqu:"skipinsert"`

	Username  string `db:"username" valid:"stringlength(1|255),required"`
	Password  string `db:"password" valid:"stringlength(5|100),required"`
	Email     string `db:"email" valid:"email,stringlength(1|254),required"`
	Slug      string `db:"slug" valid:"stringlength(1|255),required"`
	Level     int    `db:"level"`
	Banned    bool   `db:"banned"`
	Activated bool   `db:"activated"`
}
