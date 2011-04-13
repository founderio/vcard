package main

import (
	"os"
	"flag"
	"fmt"
	"bitbucket.org/llg/govcard"
)

const ( // Constant define address information index in directory information StructuredValue
	familyNames       = 0
	givenNames        = 1
	additionalNames   = 2
	honorificPrefixes = 3
	honorificSuffixes = 4
	nameSize          = honorificSuffixes + 1
)

type VCard struct {
	formattedName     string
	FamilyNames       []string
	GivenNames        []string
	AdditionalNames   []string
	HonorificNames    []string
	HonorificSuffixes []string
	NickNames         []string
	Photo             Photo
	Birthday          string
}

type AddressType int

const (
	Intl   AddressType = iota // International Delivery Address
	Postal                    // Postal Delivery Address
	Parcel                    // Parcel Delivery Address
	Work                      // Delivery Address for a place of work
	Dom                       // Domestic Delivery Address
	Pref                      // Prefered Address
)

const ( // Constant define address information index in directory information StructuredValue
	postOfficeBox   = 0
	extendedAddress = 1
	street          = 2
	locality        = 3
	region          = 4
	postalCode      = 5
	countryName     = 6
	addressSize     = countryName + 1
)

type Address struct {
	Type            []AddressType // default is Intl,Postal,Parcel,Work
	Label           string
	PostOfficeBox   string
	ExtendedAddress string
	Street          string
	Locality        string // e.g: city
	Region          string // e.g: state or province
	PostalCode      string
	CountryName     string
}

type Photo struct {
	Encoding string
	Type     string
	Value    string
	Data     []byte
}

func contentLine(group, name string, params map[string]govcard.Value, value govcard.StructuredValue) {
	fmt.Println(group, name, params, value)
}

func main() {
	for _, abpath := range flag.Args() {
		f, err := os.Open(abpath, os.O_RDONLY, 0666)
		defer f.Close()
		if err != nil {
			return
		}
		govcard.ReadDirectoryInformation(f, contentLine)
	}
}
