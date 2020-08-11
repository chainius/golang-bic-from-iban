package bic

import (
    "github.com/Skyhark-Projects/golang-bic-from-iban/log"
    "code.sajari.com/docconv"
    "path/filepath"
    "encoding/json"
    "strings"
    "strconv"
    "io/ioutil"
    "os"
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

func BeListToJSON(pdfPath string) error {
    b, err := ParseBelgiumList(pdfPath)
    if err != nil {
        return err
    }

    swifts := LoadCountrySwifts(filepath.Dir(pdfPath) + "/AllCountries/BE.json");
    for key, bank := range b {
        for _, swift := range swifts {
            if swift.SwiftCode == bank.Swift {
                bank.City = swift.City
                b[key] = bank
                break
            }
        }
    }

    bytes, err := json.Marshal(b)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(filepath.Dir(pdfPath) + "/be-swift-codes.json", bytes, 0755)
}

func LoadBelgiumListFromPDF(pdfPath string) ([]Bank, error) {
    b, err := ParseBelgiumList(pdfPath)
    swifts := LoadCountrySwifts(filepath.Dir(pdfPath) + "/AllCountries/BE.json");

    if err != nil {
        log.Error("error while loading belgium list", "error", err)
    }

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

func LoadBelgiumList(jsonPath string) ([]Bank, error) {
    /*if err := BeListToJSON(pdfPath); err != nil {
        log.Error("error while creating belgium json list", "error", err)
    }*/

    jsonFile, err := os.Open(jsonPath)
    if err != nil {
        log.Error("error while loading belgium list", "error", err)
        return []Bank{}, err
    }

    byteValue, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        log.Error("error while reading belgium list", "error", err)
        return []Bank{}, err
    }

    b := []Bank{}
    if err = json.Unmarshal(byteValue, &b); err != nil {
        log.Error("error while parsing belgium list", "error", err)
        return b, err
    }

    log.Info("Loaded Belgian banks", "count", len(b))
    for _, bank := range b {
        banks = append(banks, bank)
    }

    return banks, err
}

func LoadAllCountries() {
    dir := "./data/AllCountries/"
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Error("Could not load AllCountries", "err", err)
        return
    }

    for _, f := range files {
        if f.Name() == "BE.json" || len(f.Name()) < 2 {
            continue
        }

        country := f.Name()[:2]
        swifts := LoadCountrySwifts(dir + f.Name())
        for _, swift := range swifts {
            banks = append(banks, Bank{
                Country:    country,
                City:       swift.City,
                Start:      swift.Id,
                End:        swift.Id,
                Name:       swift.Bank,
                Swift:      swift.SwiftCode,
            })
        }
    }
}
