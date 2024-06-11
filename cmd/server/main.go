package main

import (
	"github.com/elina-chertova/auth-keeper.git/internal/config"
	"github.com/elina-chertova/auth-keeper.git/internal/db/database"
	"github.com/elina-chertova/auth-keeper.git/internal/handlers"
	"github.com/elina-chertova/auth-keeper.git/internal/middleware"
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
	dataRoutes(router, db)

	err := router.Run(appConf.Address)
	if err != nil {
		return err
	}
	return nil
}

func authRoutes(r *gin.Engine, db *gorm.DB) {
	u := repository.NewUserRepo(db)
	h := handlers.NewUserHandler(u)
	r.POST("/api/user/register", h.Register())
	r.POST("/api/user/login", h.Signup())
}

func dataRoutes(r *gin.Engine, db *gorm.DB) {
	lp := repository.NewLPRepo(db)
	bd := repository.NewBDRepo(db)
	cc := repository.NewCCRepo(db)
	td := repository.NewTDRepo(db)

	h := handlers.NewDataHandler(lp, bd, cc, td)
	r.Use(middleware.JWTAuth())
	r.Use(middleware.ExtractUserID())
	userRepo := repository.NewUserRepo(db)

	r.Use(middleware.LoadPersonalKey(userRepo))

	r.POST("/api/data/add-card", h.AddCreditCardHandler())
	r.GET("/api/data/get-card", h.GetCreditCardHandler())

	r.POST("/api/data/add-text-data", h.AddTextDataHandler())
	r.GET("/api/data/get-text-data", h.GetTextDataHandler())

	r.POST("/api/data/add-binary-data", h.AddBinaryDataHandler())
	r.GET("/api/data/get-binary-data", h.GetBinaryDataHandler())

	r.POST("/api/data/add-login-password", h.AddLoginPasswordHandler())
	r.GET("/api/data/get-login-password", h.GetLoginPasswordHandler())

}
