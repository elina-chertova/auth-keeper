package cliApp

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/sender"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
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

		resp, err := client.SendRequest("POST", "add-text-data", textData, token)
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

		fmt.Printf("Text Data got successfully: %s\n", resp.String())
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
