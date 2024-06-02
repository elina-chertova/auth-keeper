package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type LoginPassword struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	Login    string `gorm:"not null"`
	Password string `gorm:"not null"`
	Metadata string `gorm:"type:text"`
}

type TextData struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	Content  string `gorm:"not null"`
	Metadata string `gorm:"type:text"`
}

type BinaryData struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	Content  []byte `gorm:"not null"`
	Metadata string `gorm:"type:text"`
}

type CreditCard struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	CardNumber string `gorm:"not null"`
	ExpiryDate string `gorm:"not null"`
	CVV        string `gorm:"not null"`
	CardHolder string `gorm:"not null"`
	Metadata   string `gorm:"type:text"`
}
