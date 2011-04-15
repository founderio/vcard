package main

import (
	"os"
	//"flag"
	"log"
	"bitbucket.org/llg/vcard"
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
	Version           string
	FormattedName     string
	FamilyNames       []string
	GivenNames        []string
	AdditionalNames   []string
	HonorificNames    []string
	HonorificSuffixes []string
	NickNames         []string
	Photo             Photo
	Birthday          string
	Addresses         []Address
	Telephones        []Telephone
	Emails            []Email
	Title             string
	Role              string
	Org               []string
	Categories        []string
	Note              string
	URL               string
	XJabbers    	  []XJabber
	// mac specific
	XABuid    string
	XABShowAs string
}


type Photo struct {
	Encoding string
	Type     string
	Value    string
	Data     string
}

func defaultAddressTypes() (types []string) {
	return []string{"Intl", "Postal", "Parcel", "Work"}
}

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
	Type            []string // default is Intl,Postal,Parcel,Work
	Label           string
	PostOfficeBox   string
	ExtendedAddress string
	Street          string
	Locality        string // e.g: city
	Region          string // e.g: state or province
	PostalCode      string
	CountryName     string
}

type Telephone struct {
	Type   []string
	Number string
}

type Email struct {
	Type    []string
	Address string
}

type XJabber struct {
	Type    []string
	Address string
}

type VCardReader struct {
	current *VCard
}

func (v *VCardReader) contentLine(group, name string, params map[string]vcard.Value, value vcard.StructuredValue) {
	switch name {
	case "BEGIN":
		if value.GetText() == "VCARD" {
			v.current = new(VCard)
		}
	case "VERSION":
		v.current.Version = value.GetText()
	case "END":
		if value.GetText() == "VCARD" {
			v.current = nil
		}
	case "FN":
		if v.current != nil {
			v.current.FormattedName = value.GetText()
		}
	case "N":
		if len(value) == nameSize {
			v.current.FamilyNames = value[familyNames]
			v.current.GivenNames = value[givenNames]
			v.current.AdditionalNames = value[additionalNames]
			v.current.HonorificNames = value[honorificPrefixes]
			v.current.HonorificSuffixes = value[honorificSuffixes]
		} else {
			log.Printf("Error structured data isn't appropriate: %d\n", len(value))
		}
	case "NICKNAME":
		v.current.NickNames = value.GetTextList()
	case "PHOTO":
		v.current.Photo.Encoding = params["ENCODING"].GetText()
		v.current.Photo.Type = params["TYPE"].GetText()
		v.current.Photo.Value = params["VALUE"].GetText()
		v.current.Photo.Data = value.GetText()
	case "BDAY":
		v.current.Birthday = value.GetText()
	case "ADR":
		if len(value) == addressSize {
			var address Address
			if param, ok := params["TYPE"]; ok {
				address.Type = param
			} else {
				address.Type = defaultAddressTypes()
			}
			address.PostOfficeBox = value[postOfficeBox].GetText()
			address.ExtendedAddress = value[extendedAddress].GetText()
			address.Street = value[street].GetText()
			address.Locality = value[locality].GetText()
			address.Region = value[region].GetText()
			address.PostalCode = value[postalCode].GetText()
			address.CountryName = value[countryName].GetText()
			v.current.Addresses = append(v.current.Addresses, address)
		} else {
			log.Printf("Error structured data isn't appropriate: %d\n", len(value))
		}
	case "X-ABUID":
		v.current.XABuid = value.GetText()
	case "TEL":
		var tel Telephone
		if param, ok := params["TYPE"]; ok {
			tel.Type = param
		} else {
			tel.Type = []string{"voice"}
		}
		tel.Number = value.GetText()
		v.current.Telephones = append(v.current.Telephones, tel)
	case "EMAIL":
		var email Email
		if param, ok := params["TYPE"]; ok {
			email.Type = param
		} else {
			email.Type = []string{"HOME"}
		}
		email.Address = value.GetText()
		v.current.Emails = append(v.current.Emails, email)
	case "TITLE":
		v.current.Title = value.GetText()
	case "ROLE":
		v.current.Role = value.GetText()
	case "ORG":
		v.current.Org = value.GetTextList()
	case "CATEGORIES":
		v.current.Categories = value.GetTextList()
	case "NOTE":
		v.current.Note = value.GetText()
	case "URL":
		v.current.URL = value.GetText()
	case "X-JABBER": case "X-GTALK":
		var jabber XJabber
		if param, ok := params["TYPE"]; ok {
			jabber.Type = param
		} else {
			jabber.Type = []string{"HOME"}
		}
		jabber.Address = value.GetText()
		v.current.XJabbers = append(v.current.XJabbers, jabber)
	case "X-ABShowAs":
		v.current.XABShowAs = value.GetText()
	/*case "X-ABLabel":
	case "X-ABADR":
		// ignore*/
	default:
		log.Printf("Not read %s, %s: %s\n", group, name, value)
	}
}
func (v *VCardReader) GetContentLineFunc() vcard.ContentLineFunc {
	return func(group, name string, params map[string]vcard.Value, value vcard.StructuredValue) {
		v.contentLine(group, name, params, value)
	}
}

func main() {
	//for _, abpath := range flag.Args() {
	var reader VCardReader
	abpath := "contacts.vcf"
	f, err := os.Open(abpath, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		return
	}
	vcard.ReadDirectoryInformation(f, reader.GetContentLineFunc())
	log.Printf("Read %s\n", abpath)
	//}
}
