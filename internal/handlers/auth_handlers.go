package handlers

import (
	"errors"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/repository"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserHandler struct {
	userRep repository.UserRepo
}

func NewUserHandler(ur repository.UserRepo) *UserHandler {
	return &UserHandler{userRep: ur}
}

var (
	errBadRequest     = errors.New("email or password is incorrect")
	errTokenGenerated = errors.New("token is not generated")
)

func (h *UserHandler) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hashed, err := security.HashPassword(user.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Password = hashed
		err = h.userRep.CreateUser(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token, err := security.GenerateToken(user.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": errTokenGenerated})
			return
		}

		ctx.SetCookie("access_token", token, 3600, "/", "localhost", false, true)
		ctx.Writer.Header().Set("Authorization", "Bearer "+token)

		user.Password = "***"
		user.PersonalKey = []byte("***")
		ctx.IndentedJSON(
			http.StatusCreated, gin.H{
				"message": "User has been created",
				"token":   token,
				"user":    user,
			},
		)
	}
}

func (h *UserHandler) Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dbUser, err := h.userRep.GetUserByUsername(user.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": errBadRequest})
			return
		}
		hash := dbUser.Password
		isEqual := security.CheckPasswordHash(user.Password, hash)
		if !isEqual {
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": errBadRequest})
				return
			}
		}
		token, err := security.GenerateToken(user.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": errTokenGenerated})
			return
		}

		ctx.SetCookie("access_token", token, 3600, "/", "localhost", false, true)
		ctx.Writer.Header().Set("Authorization", "Bearer "+token)

		ctx.IndentedJSON(
			http.StatusOK, gin.H{
				"message": "User has logged in",
				"token":   token,
				"status":  http.StatusOK,
			},
		)
	}
}
