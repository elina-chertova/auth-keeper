package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username    string `json:"username" gorm:"unique;not null"`
	Password    string `json:"password" gorm:"not null"`
	Email       string `json:"email" gorm:"unique;not null"`
	PersonalKey []byte `json:"personal_key" gorm:"not null"`
}

type LoginPassword struct {
	gorm.Model
	UserID   string `json:"user_id" gorm:"not null"`
	Login    string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Metadata string `json:"metadata" gorm:"type:text"`
}

type TextData struct {
	gorm.Model
	UserID   string `json:"user_id" gorm:"not null"`
	Content  string `json:"content" gorm:"not null"`
	Metadata string `json:"metadata" gorm:"type:text"`
}

type BinaryData struct {
	gorm.Model
	UserID   string `json:"user_id" gorm:"not null"`
	Content  []byte `json:"content" gorm:"not null"`
	Metadata string `json:"metadata" gorm:"type:text"`
}

type CreditCard struct {
	gorm.Model
	UserID     string `json:"user_id" gorm:"not null"`
	CardNumber string `json:"card_number" gorm:"not null"`
	ExpiryDate string `json:"expiry_date" gorm:"not null"`
	CVV        string `json:"cvv" gorm:"not null"`
	CardHolder string `json:"card_holder" gorm:"not null"`
	Metadata   string `json:"metadata" gorm:"type:text"`
}
