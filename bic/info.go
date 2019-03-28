package bic

import (
    "github.com/almerlucke/go-iban/iban"
    "strconv"
    "strings"
)

type IbanInfo struct {
    Country string
    Valid bool
    IBAN string
    BBAN string
    BankCode int
    BankName string
    BankCity string
    Swift string
}

func GetInfo(ibanStr string) IbanInfo {
    iban, err := iban.NewIBAN(ibanStr)
    if err != nil {
        return IbanInfo{
            Valid: false,
        }
    }

    rule := GetCountryRules(iban.CountryCode)
    BankCode, _ := strconv.Atoi(iban.BBAN[:rule.BankCodeLength])
    bank := IbanInfo{
        Country: strings.ToUpper(iban.CountryCode),
        Valid: true,
        IBAN: iban.PrintCode,
        BBAN: iban.BBAN,
        BankCode: BankCode,
        BankCity: "",
    }

    for _, b := range banks {
        if b.Country != bank.Country || b.Start > bank.BankCode || b.End < bank.BankCode {
            continue
        }

        bank.BankName = b.Name
        bank.Swift = b.Swift
        bank.BankCity = b.City
        break;
    }

    return bank
}
