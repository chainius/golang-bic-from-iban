package main

import (
	"github.com/urfave/cli"
	"github.com/Skyhark-Projects/golang-bic-from-iban/http"
	"github.com/Skyhark-Projects/golang-bic-from-iban/log"
	"github.com/Skyhark-Projects/golang-bic-from-iban/bic"
	"errors"
)

func loadBicsDB() {
	bic.LoadBelgiumList("data/be-swift-codes.json")
	bic.LoadIbanRules("data/ibanrules.csv")
	bic.LoadAllCountries()
}

func RunServer(c *cli.Context) error {
	loadBicsDB()

    http.AddRoute("/iban/:iban", func(req http.Request) (interface{}, error) {
        bank := bic.GetInfo(req.Variables["iban"])
        return bank, nil
    })

    /*http.AddRoute("/swift/:swift", func(req http.Request) (interface{}, error) {
        bank := bic.GetInfo(req.Variables["iban"])
        return bank, nil
    })*/

    //ToDo test othe countries outside belgium, starting with DE, Nl, FR
    //log.Info("test", "iban", bank)

	http.Start("8080")
	return nil
}

func TransformBePdfCmd(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("No path provided")
	}

	if err := bic.BeListToJSON(c.Args()[0]); err != nil {
		return err
	}

	log.Info("PDF Banks lists converted to json")
	return nil
}