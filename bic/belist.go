package bic

import (
    "github.com/Skyhark-Projects/golang-bic-from-iban/log"
    "code.sajari.com/docconv"
    "path/filepath"
    "strings"
    "strconv"
)

func ParseBelgiumList(pdfPath string) ([]Bank, error) {
    res, err := docconv.ConvertPath(pdfPath)
	if err != nil {
        return []Bank{}, err
    }

    s := strings.Split(res.Body, "\n")
    lastBanks := []Bank{}
    allBanks := []Bank{}

    for _, line := range s {
        if len(line) > 29 && line[:20] == "Numbers Institution " {
            start, _ := strconv.Atoi(line[20:23])
            end, _ := strconv.Atoi(line[26:29])

            bank := Bank{
                Country: "BE",
                Start: start,
                End: end,
                Name: line[30:],
                Swift: "",
            }

            lastBanks = append(lastBanks, bank)
        } else if len(line) > 5 && line[:5] == "BIC: " {
            bics := strings.Split(line, "BIC: ")[1:]
            for _, bic := range bics {
                index := strings.Index(bic, " Bank identification")
                if index != -1 {
                    bic = bic[:index]
                } else if strings.Index(bic, "Bank identification") == 0 {
                    bic = ""
                }

                if len(lastBanks) > 0 {
                    lastBanks[0].Swift = strings.Replace(bic, " ", "", -1)
                    allBanks = append(allBanks, lastBanks[0])
                    lastBanks = lastBanks[1:]
                } else {
                    log.Error("could not assing bic to bank", "bic", bic)
                }
            }
        } else if len(line) >= 9 {
            start, err1 := strconv.Atoi(line[:3])
            end, err2 := strconv.Atoi(line[6:9])

            if err1 != nil || err2 != nil {
                continue
            }

            name := ""
            if len(line) > 10 {
                name = line[10:]
            }

            bank := Bank{
                Country: "BE",
                Start: start,
                End: end,
                Name: name,
                Swift: "",
            }

            lastBanks = append(lastBanks, bank)
        }
    }

    return allBanks, nil
}

func LoadBelgiumList(pdfPath string) ([]Bank, error) {
    b, err := ParseBelgiumList(pdfPath)
    swifts := LoadCountrySwifts(filepath.Dir(pdfPath) + "/AllCountries/BE.json");

    for _, bank := range b {
        for _, swift := range swifts {
            if swift.SwiftCode == bank.Swift {
                bank.City = swift.City
                break
            }
        }

        banks = append(banks, bank)
    }

    return banks, err
}