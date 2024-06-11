package cliApp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/elina-chertova/auth-keeper.git/internal/sender"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
)

func getAddCardFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "card_number",
			Aliases:  []string{"cn"},
			Usage:    "Credit Card Number",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "expiry_date",
			Aliases:  []string{"ed"},
			Usage:    "Expiry Date",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "cvv",
			Aliases:  []string{"cv"},
			Usage:    "CVV",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "card_holder",
			Aliases:  []string{"ch"},
			Usage:    "Card Holder",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "metadata",
			Aliases:  []string{"m"},
			Usage:    "Metadata",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "Token for Authorization",
			Required: true,
		},
	}
}

func AddCard(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		creditCard := models.CreditCard{
			CardNumber: c.String("card_number"),
			ExpiryDate: c.String("expiry_date"),
			CVV:        c.String("cvv"),
			CardHolder: c.String("card_holder"),
			Metadata:   c.String("metadata"),
		}
		token := c.String("token")
		client := sender.NewClient(baseURL)

		jsonData, err := json.Marshal(creditCard)
		if err != nil {
			log.Fatalf("Error marshalling data: %v", err)
		}

		personalKey, err := os.ReadFile("pkey.txt")
		if err != nil {
			log.Fatalf("Error reading personal key: %v", err)
		}

		encryptedData, err := security.EncryptData(jsonData, personalKey)
		if err != nil {
			log.Fatalf("Error encrypting data: %v", err)
		}

		encodedData := base64.StdEncoding.EncodeToString(encryptedData)
		fmt.Printf("Encrypted Data: %s\n", encodedData)

		resp, err := client.SendRequest(
			"POST",
			"add-card",
			map[string]string{"data": encodedData},
			token,
		)
		if err != nil {
			log.Fatalf("Error adding credit card: %v", err)
		}
		if err != nil {
			log.Fatalf("Error adding credit card: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			log.Fatalf(
				"Failed to add credit card, status code: %d, response: %s",
				resp.StatusCode,
				resp.String(),
			)
		}

		fmt.Printf("Credit Card added successfully: %s\n", resp.String())
		return nil
	}
}

func AddCardCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "add-card",
		Usage:  "Add a credit card",
		Flags:  getAddCardFlags(),
		Action: AddCard(baseURL),
	}
}

func getCardFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "Token for Authorization",
			Required: true,
		},
	}
}

func GetCard(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		token := c.String("token")
		client := sender.NewClient(baseURL)
		resp, err := client.SendRequest("GET", "get-card", nil, token)
		if err != nil {
			log.Fatalf("Error getting credit card: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf(
				"Failed to get credit card, status code: %d, response: %s",
				resp.StatusCode,
				resp.String(),
			)
		}

		var responseData struct {
			Body    string `json:"body"`
			Message string `json:"message"`
		}

		if err := json.Unmarshal(resp.Bytes(), &responseData); err != nil {
			log.Fatalf("Error unmarshalling response data: %v", err)
		}

		encodedData := responseData.Body
		fmt.Printf("Encoded Data: %s\n", encodedData)

		personalKey, err := os.ReadFile("pkey.txt")
		if err != nil {
			log.Fatalf("Error reading personal key: %v", err)
		}

		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			log.Fatalf("Error decoding base64 data: %v", err)
		}

		decryptedData, err := security.DecryptData(decodedData, personalKey)
		if err != nil {
			log.Fatalf("Error decrypting data: %v", err)
		}

		fmt.Printf("Decrypted Data: %s\n", string(decryptedData))
		return nil
	}
}

func GetCardCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "get-card",
		Usage:  "Get credit cards",
		Flags:  getCardFlags(),
		Action: GetCard(baseURL),
	}
}
