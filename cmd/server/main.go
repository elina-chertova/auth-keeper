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
	authorized := r.Group("/", middleware.JWTAuth())
	authorized.POST("/api/data/add-card", h.AddCreditCardHandler())
	authorized.GET("/api/data/get-card", h.GetCreditCardHandler())

	//r.POST("/api/data/add-card", middleware.JWTAuth(), h.AddCreditCardHandler())
	//r.POST("/api/data/get-card", middleware.JWTAuth(), h.GetCreditCardHandler())
}
