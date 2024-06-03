package database

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DBConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	DBName   string
}

func InitDB(conf *DBConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error during database initialization: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.CreditCard{},
		&models.BinaryData{},
		&models.TextData{},
		&models.LoginPassword{},
	)
	if err != nil {
		log.Fatalf("Error during migration: %v", err)
	}
	return db
}
