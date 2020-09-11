package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `binding:"required,len=11" gorm:"type:varchar(110);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}
