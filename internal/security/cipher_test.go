package security

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"log"
	"os"
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	originalData := models.CreditCard{
		CardNumber: "1234567890123456",
		ExpiryDate: "12/25",
		CVV:        "123",
		CardHolder: "John Doe",
		Metadata:   "",
	}
	personalKey, err := GeneratePersonalKey()
	if err != nil {
		log.Fatalf("Error generating personal key: %v", err)
	}

	err = os.WriteFile("pkey.txt", personalKey, 0644)
	if err != nil {
		log.Fatalf("Error saving personal key to file: %v", err)
	}

	encryptedPersonalKey, err := EncryptPersonalKey(personalKey)
	personalKey, err = DecryptPersonalKey(encryptedPersonalKey)

	jsonData, err := json.Marshal(originalData)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
	}

	encryptedData, err := EncryptData(jsonData, personalKey)
	if err != nil {
		log.Fatalf("Error encrypting data: %v", err)
	}
	encodedData := base64.StdEncoding.EncodeToString(encryptedData)

	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		log.Fatalf("Error decoding data: %v", err)
	}

	decryptedData, err := DecryptData(decodedData, personalKey)
	if err != nil {
		log.Fatalf("Error decrypting data: %v", err)
	}

	var decryptedCard models.CreditCard
	if err := json.Unmarshal(decryptedData, &decryptedCard); err != nil {
		log.Fatalf("Error unmarshalling decrypted data: %v", err)
	}

	fmt.Printf("Original Data: %+v\n", originalData)
	fmt.Printf("Decrypted Data: %+v\n", decryptedCard)
}
