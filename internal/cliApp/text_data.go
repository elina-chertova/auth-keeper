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

func getAddTextDataFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "content",
			Aliases:  []string{"c"},
			Usage:    "Content",
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

func AddTextData(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		textData := models.TextData{
			Content:  c.String("content"),
			Metadata: c.String("metadata"),
		}
		token := c.String("token")
		client := sender.NewClient(baseURL)

		jsonData, err := json.Marshal(textData)
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
			"add-text-data",
			map[string]string{"data": encodedData},
			token,
		)
		if err != nil {
			log.Fatalf("Error adding text data: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			log.Fatalf(
				"Failed to add text data, status code: %d, response: %s",
				resp.StatusCode,
				resp.String(),
			)
		}

		fmt.Printf("Text Data added successfully: %s\n", resp.String())
		return nil
	}
}

func AddTextDataCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "add-text-data",
		Usage:  "Add a text data",
		Flags:  getAddTextDataFlags(),
		Action: AddTextData(baseURL),
	}
}

func getTextDataFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "Token for Authorization",
			Required: true,
		},
	}
}

func GetTextData(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		token := c.String("token")
		client := sender.NewClient(baseURL)
		resp, err := client.SendRequest("GET", "get-text-data", nil, token)
		if err != nil {
			log.Fatalf("Error getting text data: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf(
				"Failed to get text data, status code: %d, response: %s",
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

func GetTextDataCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "get-text-data",
		Usage:  "Get text data",
		Flags:  getTextDataFlags(),
		Action: GetTextData(baseURL),
	}
}
