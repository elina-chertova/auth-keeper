package handlers

import (
	"errors"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/repository"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
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
		expirationTime := time.Now().Add(72 * time.Hour)
		cookie := http.Cookie{
			Name:     "access_token",
			Value:    token,
			Expires:  expirationTime,
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(ctx.Writer, &cookie)
		ctx.Writer.Header().Set("Authorization", "Bearer "+token)

		user.Password = "***"
		ctx.IndentedJSON(
			http.StatusCreated, gin.H{
				"message": "User has been created",
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

		expirationTime := time.Now().Add(72 * time.Hour)
		cookie := http.Cookie{
			Name:     "access_token",
			Value:    token,
			Expires:  expirationTime,
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(ctx.Writer, &cookie)
		ctx.Writer.Header().Set("Authorization", "Bearer "+token)

		ctx.IndentedJSON(
			http.StatusOK, gin.H{
				"message": "User has been created",
				"status":  http.StatusOK,
			},
		)
	}
}
