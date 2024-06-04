package main

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/cliApp"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	name := "PasswordKeeper"
	baseURL := "http://localhost:8081"
	baseURLAuth := fmt.Sprintf("%s%s", baseURL, "/api/user/")
	baseURLData := fmt.Sprintf("%s%s", baseURL, "/api/data/")
	app := &cli.App{
		Name:  name,
		Usage: "Password Keeper CLI",
		Commands: []*cli.Command{
			cliApp.RegisterCommand(baseURLAuth),
			cliApp.LoginCommand(baseURLAuth),

			cliApp.AddCardCommand(baseURLData),
			cliApp.GetCardCommand(baseURLData),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
