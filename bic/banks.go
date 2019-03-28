package bic

type Bank struct {
    Country string
    City string
    Start int
    End int
    Name string
    Swift string
}

var banks = []Bank{}

func GetSwiftBank(swift string) *Bank {
    for _, bank := range banks {
        if bank.Swift == swift {
            return &bank
        }
    }

    return nil
}