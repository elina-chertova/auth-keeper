package cliApp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/sender"
	"github.com/urfave/cli/v2"
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

		resp, err := client.SendRequest("POST", "add-card", creditCard, token)
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

		fmt.Printf("Credit Card got successfully: %s\n", resp.String())
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
