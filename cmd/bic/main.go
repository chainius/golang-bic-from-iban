package main

import (
    "github.com/Skyhark-Projects/golang-bic-from-iban/http"
    "github.com/Skyhark-Projects/golang-bic-from-iban/bic"
    "github.com/Skyhark-Projects/golang-bic-from-iban/log"
    "os"
)

func main() {
    log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

    bic.LoadBelgiumList("data/be-swift-codes.json")
    bic.LoadIbanRules("data/ibanrules.csv")

    http.AddRoute("/iban/:iban", func(req http.Request) (interface{}, error) {
        bank := bic.GetInfo(req.Variables["iban"])
        return bank, nil
    })

    /*http.AddRoute("/swift/:swift", func(req http.Request) (interface{}, error) {
        bank := bic.GetInfo(req.Variables["iban"])
        return bank, nil
    })*/

    //ToDo test othe countries outside belgium, starting with DE, Nl, FR
    //log.Info("test", "iban", bank)

    http.Start("8080")
}