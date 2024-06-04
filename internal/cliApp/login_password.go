package cliApp

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/elina-chertova/auth-keeper.git/internal/sender"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
)

func getAddLoginPasswordFlags() []cli.Flag {
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

func AddLoginPassword(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		lpData := models.LoginPassword{
			Login:    c.String("username"),
			Password: c.String("password"),
			Metadata: c.String("metadata"),
		}
		token := c.String("token")
		client := sender.NewClient(baseURL)

		resp, err := client.SendRequest("POST", "add-login-password", lpData, token)
		if err != nil {
			log.Fatalf("Error adding login-password data: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			log.Fatalf(
				"Failed to add login-password data, status code: %d, response: %s",
				resp.StatusCode,
				resp.String(),
			)
		}
		fmt.Printf("Login-Password added successfully: %s\n", resp.String())
		return nil
	}
}
func AddLoginPasswordCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "add-login-password",
		Usage:  "Add a login-password data",
		Flags:  getAddLoginPasswordFlags(),
		Action: AddLoginPassword(baseURL),
	}
}

func getLoginPasswordFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "Token for Authorization",
			Required: true,
		},
	}
}

func GetLoginPassword(baseURL string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		token := c.String("token")
		client := sender.NewClient(baseURL)
		resp, err := client.SendRequest("GET", "get-login-password", nil, token)
		if err != nil {
			log.Fatalf("Error getting login-password data: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf(
				"Failed to get login-password data, status code: %d, response: %s",
				resp.StatusCode,
				resp.String(),
			)
		}

		fmt.Printf("Login-Password got successfully: %s\n", resp.String())
		return nil
	}
}
func GetLoginPasswordCommand(baseURL string) *cli.Command {
	return &cli.Command{
		Name:   "get-login-password",
		Usage:  "Get login-password data",
		Flags:  getLoginPasswordFlags(),
		Action: GetLoginPassword(baseURL),
	}
}
