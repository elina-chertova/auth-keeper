package main

import (
	"github.com/elina-chertova/auth-keeper.git/internal/cliApp"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	name := "PasswordKeeper"
	baseURL := "http://localhost:8080/"
	app := &cli.App{
		Name:  name,
		Usage: "Password Keeper CLI",
		Commands: []*cli.Command{
			cliApp.RegisterCommand(baseURL),
			cliApp.LoginCommand(baseURL),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
