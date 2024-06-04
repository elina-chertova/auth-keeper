package handlers

import (
	"errors"
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/repository"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var errTokenNotFound = errors.New("token not found")

type DataHandler struct {
	lp repository.LoginPasswordRepo
	bd repository.BinaryDataRepo
	cc repository.CreditCardRepo
	td repository.TextDataRepo
}

func NewDataHandler(
	lp repository.LoginPasswordRepo,
	bd repository.BinaryDataRepo,
	cc repository.CreditCardRepo,
	td repository.TextDataRepo,
) *DataHandler {
	return &DataHandler{
		lp: lp,
		bd: bd,
		cc: cc,
		td: td,
	}
}

func extractUserFromRequest(ctx *gin.Context) (string, error) {
	token, exists := ctx.Get("token")

	if !exists {
		return "", errTokenNotFound
	}

	tokenStr := fmt.Sprintf("%v", token)
	userID, err := security.GetUserFromToken(tokenStr)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (dh *DataHandler) AddCreditCardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := extractUserFromRequest(ctx)
		if err != nil {
			if errors.Is(err, errTokenNotFound) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			log.Printf("Error extracting user from token: %v", err)
			return
		}

		var creditCard models.CreditCard
		if err := ctx.ShouldBindJSON(&creditCard); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		creditCard.UserID = userID
		err = dh.cc.SaveNewCreditCard(&creditCard)
		if err != nil {
			log.Printf("Error adding credit card: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.IndentedJSON(
			http.StatusCreated, gin.H{
				"message": "Credit Card has been added",
				"status":  http.StatusCreated,
			},
		)
	}
}

func (dh *DataHandler) GetCreditCardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := extractUserFromRequest(ctx)
		if err != nil {
			if errors.Is(err, errTokenNotFound) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			log.Printf("Error extracting user from token: %v", err)
			return
		}

		creditCard, err := dh.cc.GetCreditCardList(userID)
		if err != nil {
			log.Printf("Error getting credit card: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.IndentedJSON(
			http.StatusCreated, gin.H{
				"message": "User has been created",
				"body":    creditCard,
			},
		)

	}
}
