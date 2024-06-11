package middleware

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		accessTokenBearer := c.GetHeader("Authorization")
		if accessTokenBearer != "" {
			token := strings.TrimPrefix(accessTokenBearer, "Bearer ")
			if token == accessTokenBearer {
				c.JSON(
					http.StatusUnauthorized,
					response{
						Message: "Check token",
						Status:  "Invalid Authorization Token format",
					},
				)
				c.Abort()
				return
			}

			err := security.ValidateToken(token)
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response{
						Message: err.Error(),
						Status:  "Unauthorized",
					},
				)
				return
			}

			c.Set("token", token)
			c.Next()
			return
		}

		accessTokenCookie, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				response{
					Message: err.Error(),
					Status:  "Unauthorized",
				},
			)
			c.Abort()
			return
		}

		err = security.ValidateToken(accessTokenCookie)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				response{
					Message: err.Error(),
					Status:  "Unauthorized",
				},
			)
			c.Abort()
			return
		}

		c.Set("token", accessTokenCookie)
		c.Next()
	}
}

func ExtractUserID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, exists := ctx.Get("token")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			ctx.Abort()
			return
		}

		tokenStr := fmt.Sprintf("%v", token)
		userID, err := security.GetUserFromToken(tokenStr)
		if err != nil {
			log.Printf("Error extracting user from token: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}
