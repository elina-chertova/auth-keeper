package cliApp

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
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
		textData := models.BinaryData{
			Content:  content,
			Metadata: c.String("metadata"),
		}
		token := c.String("token")
		client := sender.NewClient(baseURL)

		resp, err := client.SendRequest("POST", "add-binary-data", textData, token)
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

		fmt.Printf("Binary Data got successfully: %s\n", resp.String())
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
