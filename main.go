package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ciphers"
	app.Usage = "cipher text"
	app.Action = func(c *cli.Context) error {
		fmt.Println("placeholder")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
