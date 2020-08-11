package main

import (
	"github.com/Skyhark-Projects/golang-bic-from-iban/bic"
    "github.com/Skyhark-Projects/golang-bic-from-iban/log"
	"github.com/urfave/cli"
	"strings"
	"errors"
	"fmt"
    "os"
)

var commands = []cli.Command{
    {
        Name: "server",
        Action: RunServer,
	},
	{
        Name: "iban",
        Action: CmdGetIban,
    },
    {
        Name: "parse-pdf",
        Action: TransformBePdfCmd,
    },
}

func main() {
    log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

    app := cli.NewApp()
	app.Name = "IBAN-tools"
	cli.OsExiter = func(int) {
		os.Exit(1)
	}

	app.After = func(c *cli.Context) error {
		return nil
	}

	app.Usage = "bit4you"
	app.Version = "1.0.0"
	app.Action = RunServer
	app.EnableBashCompletion = true
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Vandamme Sacha",
			Email: "sacha@bit4you.io",
		},
	}

	app.Commands = commands

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
    }
}

// ------

func CmdGetIban(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("No iban provided")
	}

	loadBicsDB()
	iban := c.Args()[0]
	log.Info("Getting bic from iban", "iban", iban)
	fmt.Println("")

	bank := bic.GetInfo(strings.ReplaceAll(iban, " ", ""))
	log.Warn(bank.IBAN, "valid", bank.Valid, "country", bank.Country)
	log.Warn(bank.BankName, "swift", bank.Swift, "code", bank.BankCode)

	return nil
}