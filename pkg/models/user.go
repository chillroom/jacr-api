package models

type User struct {
	ID int

	Username  string `valid:"stringlength(1|255),required"`
	Password  string `valid:"stringlength(5|100),required"`
	Email     string `valid:"email,stringlength(1|254),required"`
	Slug      string `valid:"stringlength(1|255),required"`
	Level     int
	Banned    bool
	Activated bool
}
