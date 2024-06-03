package cliApp

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/sender"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
)

func getLoginFlags() []cli.Flag {
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

func loginUser(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		username := c.String("username")
		password := c.String("password")
		email := c.String("email")

		user := &models.User{
			Username: username,
			Password: password,
			Email:    email,
		}

		client := sender.NewClient(baseURL)
		resp, err := client.SendRequest("POST", "login", user)
		if err != nil {
			log.Fatalf("Error signing up user: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf(
				"Failed to log in user, status code: %d, response: %s",
				resp.StatusCode,
				resp.String(),
			)
		}

		fmt.Printf("Login successful: %s\n", resp.String())
		return nil
	}

}

func LoginCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "login",
		Usage:  "User sign up",
		Flags:  getLoginFlags(),
		Action: loginUser(baseURL),
	}
}
