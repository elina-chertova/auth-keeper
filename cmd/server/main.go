package main

import (
	"github.com/elina-chertova/auth-keeper.git/internal/config"
	"github.com/elina-chertova/auth-keeper.git/internal/db/database"
	"github.com/elina-chertova/auth-keeper.git/internal/handlers"
	"github.com/elina-chertova/auth-keeper.git/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	router := gin.Default()
	dbConf, appConf := config.LoadEnv()

	db := database.InitDB(&dbConf)
	authRoutes(router, db)

	err := router.Run(appConf.Address)
	if err != nil {
		return err
	}
	return nil
}

func authRoutes(r *gin.Engine, db *gorm.DB) {
	u := repository.NewUserRepo(db)
	h := handlers.NewUserHandler(u)
	r.POST("register", h.Register())
	r.POST("login", h.Signup())
}
