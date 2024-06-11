package cliApp

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/security"
	"github.com/elina-chertova/auth-keeper.git/internal/sender"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
)

func getRegisterFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"u"},
			Usage:    "Username",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "Password",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "email",
			Aliases:  []string{"e"},
			Usage:    "Email",
			Required: false,
		},
	}
}

func registerUser(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		username := c.String("username")
		password := c.String("password")
		email := c.String("email")

		personalKey, err := security.GeneratePersonalKey()
		if err != nil {
			log.Fatalf("Error generating personal key: %v", err)
		}

		err = os.WriteFile("pkey.txt", personalKey, 0644)
		if err != nil {
			log.Fatalf("Error saving personal key to file: %v", err)
		}

		encryptedPersonalKey, err := security.EncryptPersonalKey(personalKey)
		if err != nil {
			log.Fatalf("Error encrypting personal key: %v", err)
		}

		user := &models.User{
			Username:    username,
			Password:    password,
			Email:       email,
			PersonalKey: encryptedPersonalKey,
		}

		client := sender.NewClient(baseURL)
		resp, err := client.SendRequest("POST", "register", user, "")
		if err != nil {
			log.Fatalf("Error registering user: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			log.Fatalf(
				"Failed to register user, status code: %d, response: %s",
				resp.StatusCode,
				resp.String(),
			)
		}

		fmt.Printf("User registered successfully: %s\n", resp.String())
		return nil
	}
}

func RegisterCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "register",
		Usage:  "Register a new user",
		Flags:  getRegisterFlags(),
		Action: registerUser(baseURL),
	}
}
