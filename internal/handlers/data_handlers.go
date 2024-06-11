package handlers

import (
	"encoding/base64"
	"encoding/json"
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
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var encryptedData struct {
			Data string `json:"data"`
		}
		if err := ctx.ShouldBindJSON(&encryptedData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		decodedData, err := base64.StdEncoding.DecodeString(encryptedData.Data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode base64 data"})
			return
		}
		personalKey, exists := ctx.Get("personalKey")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Personal key not found"})
			return
		}

		decryptedData, err := security.DecryptData(decodedData, personalKey.([]byte))
		if err != nil {
			log.Printf("Failed to decrypt data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt data"})
			return
		}

		var creditCard models.CreditCard
		if err := json.Unmarshal(decryptedData, &creditCard); err != nil {
			log.Printf("Failed to unmarshal decrypted data: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unmarshal decrypted data"})
			return
		}

		creditCard.UserID = userID.(string)
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

		creditCardsJSON, err := json.Marshal(creditCard)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to marshal credit cards"},
			)
			return
		}

		personalKey, exists := ctx.Get("personalKey")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Personal key not found"})
			return
		}

		encryptedData, err := security.EncryptData(creditCardsJSON, personalKey.([]byte))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt data"})
			return
		}

		encodedData := base64.StdEncoding.EncodeToString(encryptedData)

		ctx.IndentedJSON(
			http.StatusOK, gin.H{
				"message": "Credit card list",
				"body":    encodedData,
			},
		)
	}
}

func (dh *DataHandler) AddBinaryDataHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var encryptedData struct {
			Data string `json:"data"`
		}
		if err := ctx.ShouldBindJSON(&encryptedData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		decodedData, err := base64.StdEncoding.DecodeString(encryptedData.Data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode base64 data"})
			return
		}

		personalKey, exists := ctx.Get("personalKey")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Personal key not found"})
			return
		}

		decryptedData, err := security.DecryptData(decodedData, personalKey.([]byte))
		if err != nil {
			log.Printf("Failed to decrypt data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt data"})
			return
		}

		var binaryData models.BinaryData
		if err := json.Unmarshal(decryptedData, &binaryData); err != nil {
			log.Printf("Failed to unmarshal decrypted data: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unmarshal decrypted data"})
			return
		}

		binaryData.UserID = userID.(string)
		err = dh.bd.SaveNewBinaryData(&binaryData)
		if err != nil {
			log.Printf("Error adding binary data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.IndentedJSON(
			http.StatusCreated, gin.H{
				"message": "Binary data have been added",
				"status":  http.StatusCreated,
			},
		)
	}
}

func (dh *DataHandler) GetBinaryDataHandler() gin.HandlerFunc {
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

		binaryData, err := dh.bd.GetBinaryData(userID)
		if err != nil {
			log.Printf("Error getting binary data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		binaryDataJSON, err := json.Marshal(binaryData)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to marshal binary data"},
			)
			return
		}

		personalKey, exists := ctx.Get("personalKey")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Personal key not found"})
			return
		}

		encryptedData, err := security.EncryptData(binaryDataJSON, personalKey.([]byte))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt data"})
			return
		}

		encodedData := base64.StdEncoding.EncodeToString(encryptedData)

		ctx.IndentedJSON(
			http.StatusOK, gin.H{
				"message": "Binary data list",
				"body":    encodedData,
			},
		)
	}
}

func (dh *DataHandler) AddTextDataHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var encryptedData struct {
			Data string `json:"data"`
		}
		if err := ctx.ShouldBindJSON(&encryptedData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		decodedData, err := base64.StdEncoding.DecodeString(encryptedData.Data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode base64 data"})
			return
		}

		personalKey, exists := ctx.Get("personalKey")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Personal key not found"})
			return
		}

		decryptedData, err := security.DecryptData(decodedData, personalKey.([]byte))
		if err != nil {
			log.Printf("Failed to decrypt data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt data"})
			return
		}

		var textData models.TextData
		if err := json.Unmarshal(decryptedData, &textData); err != nil {
			log.Printf("Failed to unmarshal decrypted data: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unmarshal decrypted data"})
			return
		}

		textData.UserID = userID.(string)
		err = dh.td.SaveNewTextData(&textData)
		if err != nil {
			log.Printf("Error adding text data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.IndentedJSON(
			http.StatusCreated, gin.H{
				"message": "Text data have been added",
				"status":  http.StatusCreated,
			},
		)
	}
}

func (dh *DataHandler) GetTextDataHandler() gin.HandlerFunc {
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

		textData, err := dh.td.GetTextData(userID)
		if err != nil {
			log.Printf("Error getting text data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		textDataJSON, err := json.Marshal(textData)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to marshal text data"},
			)
			return
		}

		personalKey, exists := ctx.Get("personalKey")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Personal key not found"})
			return
		}

		encryptedData, err := security.EncryptData(textDataJSON, personalKey.([]byte))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt data"})
			return
		}

		encodedData := base64.StdEncoding.EncodeToString(encryptedData)

		ctx.IndentedJSON(
			http.StatusOK, gin.H{
				"message": "Text data list",
				"body":    encodedData,
			},
		)
	}
}

func (dh *DataHandler) AddLoginPasswordHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var encryptedData struct {
			Data string `json:"data"`
		}
		if err := ctx.ShouldBindJSON(&encryptedData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		decodedData, err := base64.StdEncoding.DecodeString(encryptedData.Data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode base64 data"})
			return
		}

		personalKey, exists := ctx.Get("personalKey")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Personal key not found"})
			return
		}

		decryptedData, err := security.DecryptData(decodedData, personalKey.([]byte))
		if err != nil {
			log.Printf("Failed to decrypt data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt data"})
			return
		}

		var loginPassword models.LoginPassword
		if err := json.Unmarshal(decryptedData, &loginPassword); err != nil {
			log.Printf("Failed to unmarshal decrypted data: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unmarshal decrypted data"})
			return
		}

		hashed, err := security.HashPassword(loginPassword.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		loginPassword.Password = hashed
		loginPassword.UserID = userID.(string)
		err = dh.lp.SaveNewLoginPassword(&loginPassword)
		if err != nil {
			log.Printf("Error adding new login and password: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.IndentedJSON(
			http.StatusCreated, gin.H{
				"message": "Login-Password data have been added",
				"status":  http.StatusCreated,
			},
		)
	}
}

func (dh *DataHandler) GetLoginPasswordHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := extractUserFromRequest(ctx)
		if err != nil {
			if errors.Is(err, errTokenNotFound) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			log.Printf("error extracting user from token: %v", err)
			return
		}

		lpData, err := dh.lp.GetLoginPasswordData(userID)
		if err != nil {
			log.Printf("error getting login-password data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lpDataJSON, err := json.Marshal(lpData)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to marshal login-password data"},
			)
			return
		}

		personalKey, exists := ctx.Get("personalKey")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Personal key not found"})
			return
		}

		encryptedData, err := security.EncryptData(lpDataJSON, personalKey.([]byte))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt data"})
			return
		}

		encodedData := base64.StdEncoding.EncodeToString(encryptedData)

		ctx.IndentedJSON(
			http.StatusOK, gin.H{
				"message": "Login-Password data list",
				"body":    encodedData,
			},
		)
	}
}
