package main

import (
	"io"
	"log"
	"os"
	"unicode"

	"github.com/urfave/cli"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func main() {
	app := cli.NewApp()
	app.Name = "ciphers"
	app.Usage = "cipher text"

	var shift int
	var keyword string
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "decipher, d",
			Usage: "Decipher the input",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "caesar",
			Aliases: []string{"c"},
			Usage:   "Caesar cipher the input",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "shift, s",
					Value:       13,
					Usage:       "Right shift by `SHIFT`",
					Destination: &shift,
				},
			},
			Action: func(c *cli.Context) error {
				doIt(c, Caesar(shift))
				return nil
			},
		},
		{
			Name:    "rot13",
			Aliases: []string{"r"},
			Usage:   "Rot13 cipher the input",
			Action: func(c *cli.Context) error {
				doIt(c, Rot13())
				return nil
			},
		},
		{
			Name:    "keyword",
			Aliases: []string{"k"},
			Usage:   "Keyword cipher the input",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "keyword, k",
					Usage:       "Cipher using `KEYWORD`",
					Destination: &keyword,
				},
			},
			Action: func(c *cli.Context) error {
				doIt(c, Keyword(keyword))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func doIt(ctx *cli.Context, cipher Cipher) {
	var cipherTransform transform.Transformer
	if ctx.GlobalBool("decipher") {
		cipherTransform = cipher.Encipher()
	} else {
		cipherTransform = cipher.Decipher()
	}
	t := transform.Chain(
		norm.NFKD,
		runes.Remove(runes.In(unicode.Mark)),
		runes.Map(func(r rune) rune {
			return unicode.ToUpper(r)
		}),
		cipherTransform,
	)
	reader := transform.NewReader(os.Stdin, t)

	_, err := io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}
}
