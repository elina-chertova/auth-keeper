package config

import (
	"github.com/elina-chertova/auth-keeper.git/internal/db/database"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

type AppConf struct {
	Address string
}

var SecretKey string

func LoadEnv() (database.DBConfig, AppConf) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	viper.AutomaticEnv()

	dbConf := database.DBConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetInt64("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		DBName:   viper.GetString("DB_NAME"),
		Password: viper.GetString("DB_PASSWORD"),
	}
	appConf := AppConf{
		Address: viper.GetString("APP_ADDRESS"),
	}
	SecretKey = os.Getenv("SECRET_KEY")
	return dbConf, appConf
}
