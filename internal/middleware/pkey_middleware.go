package middleware

import (
	"github.com/elina-chertova/auth-keeper.git/internal/repository"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func LoadPersonalKey(userRepo repository.UserRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			ctx.Abort()
			return
		}

		user, err := userRepo.GetUserByUsername(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			ctx.Abort()
			return
		}

		personalKey, err := security.DecryptPersonalKey(user.PersonalKey)
		if err != nil {
			log.Printf("Error decrypting personal key: %v", err)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to decrypt personal key"},
			)
			ctx.Abort()
			return
		}

		ctx.Set("personalKey", personalKey)
		ctx.Next()
	}
}
