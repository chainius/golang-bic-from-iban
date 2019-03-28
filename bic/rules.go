package bic

import (
    "github.com/Skyhark-Projects/golang-bic-from-iban/log"
    "encoding/json"
    "encoding/csv"
    "io/ioutil"
    "strconv"
    "strings"
    "bufio"
    "os"
    "io"
)

type Format struct {
	CountryCode    string
	Country        string
	IbanLength     int
	BankCodeLength int
}

type Swifts struct {
    Id int
    Bank string
    City string
    SwiftCode string `json:"swift_code"`
}

type SwiftCountryFile struct {
    Country string
    CountryCode string `json:country_code`
    List []Swifts
}

var formats = []Format{}

func LoadIbanRules(path string) {
    csvFile, _ := os.Open(path)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Error("Error parsing iban rules", "error", error)
            return
		}

		ibanLength, err := strconv.Atoi(line[2])
        if err != nil {
            continue
        }

		bankCodeLength, err := strconv.Atoi(line[3])
		if err != nil {
            continue
        }

		format := Format{
			Country:        line[0],
			CountryCode:    line[1],
			IbanLength:     ibanLength,
			BankCodeLength: bankCodeLength,
		}

        formats = append(formats, format)
	}
}

func LoadCountrySwifts(path string) []Swifts {
    //https://github.com/PeterNotenboom/SwiftCodes/tree/master/AllCountries
    jsonFile, _ := os.Open(path)
    byteValue, _ := ioutil.ReadAll(jsonFile)
    country := SwiftCountryFile{}
    json.Unmarshal(byteValue, &country)
    return country.List
}

func GetCountryRules(country string) Format {
    country = strings.ToUpper(country)
    for _, format := range formats {
        if format.CountryCode == country {
            return format
        }
    }

    return Format{
        CountryCode: country,
    }
}