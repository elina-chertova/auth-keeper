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

func getAddBinaryDataFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "file_path",
			Aliases:  []string{"f"},
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

func AddBinaryData(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filePath := c.String("file_path")
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
		binaryData := models.BinaryData{
			Content:  content,
			Metadata: c.String("metadata"),
		}
		token := c.String("token")
		client := sender.NewClient(baseURL)

		jsonData, err := json.Marshal(binaryData)
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
			"add-binary-data",
			map[string]string{"data": encodedData},
			token,
		)
		if err != nil {
			log.Fatalf("Error adding binary data: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			log.Fatalf(
				"Failed to add binary data, status code: %d, response: %s",
				resp.StatusCode,
				resp.String(),
			)
		}
		fmt.Printf("Binary Data added successfully: %s\n", resp.String())
		return nil
	}
}

func AddBinaryDataCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "add-binary-data",
		Usage:  "Add a binary data",
		Flags:  getAddBinaryDataFlags(),
		Action: AddBinaryData(baseURL),
	}
}

func getBinaryDataFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "Token for Authorization",
			Required: true,
		},
	}
}

func GetBinaryData(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		token := c.String("token")
		client := sender.NewClient(baseURL)
		resp, err := client.SendRequest("GET", "get-binary-data", nil, token)
		if err != nil {
			log.Fatalf("Error getting binary data: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf(
				"Failed to get binary data, status code: %d, response: %s",
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

func GetBinaryDataCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "get-binary-data",
		Usage:  "Get binary data",
		Flags:  getBinaryDataFlags(),
		Action: GetBinaryData(baseURL),
	}
}
