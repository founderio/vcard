package vcard

import (
	"io"
	"log"
)


type AddressBook struct {
	Contacts []VCard
}

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
	XJabbers          []XJabber
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

func (ab *AddressBook) LastContact() *VCard {
	if len(ab.Contacts) > 0 {
		return &ab.Contacts[len(ab.Contacts)-1]
	}
	return nil
}

const ( // Constant define address information index in directory information StructuredValue
	familyNames       = 0
	givenNames        = 1
	additionalNames   = 2
	honorificPrefixes = 3
	honorificSuffixes = 4
	nameSize          = honorificSuffixes + 1
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

func (ab *AddressBook) ContentLine(group, name string, params map[string]Value, value StructuredValue) {
	current := ab.LastContact()
	switch name {
	case "BEGIN":
		if value.GetText() == "VCARD" {
			ab.Contacts = append(ab.Contacts, VCard{})
		}
	case "VERSION":
		current.Version = value.GetText()
	case "END":
		if value.GetText() == "VCARD" {

		}
	case "FN":
		if current != nil {
			current.FormattedName = value.GetText()
		}
	case "N":
		if len(value) == nameSize {
			current.FamilyNames = value[familyNames]
			current.GivenNames = value[givenNames]
			current.AdditionalNames = value[additionalNames]
			current.HonorificNames = value[honorificPrefixes]
			current.HonorificSuffixes = value[honorificSuffixes]
		} else {
			log.Printf("Error structured data isn't appropriate: %d\n", len(value))
		}
	case "NICKNAME":
		current.NickNames = value.GetTextList()
	case "PHOTO":
		current.Photo.Encoding = params["ENCODING"].GetText()
		current.Photo.Type = params["TYPE"].GetText()
		current.Photo.Value = params["VALUE"].GetText()
		current.Photo.Data = value.GetText()
	case "BDAY":
		current.Birthday = value.GetText()
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
			current.Addresses = append(current.Addresses, address)
		} else {
			log.Printf("Error structured data isn't appropriate: %d\n", len(value))
		}
	case "X-ABUID":
		current.XABuid = value.GetText()
	case "TEL":
		var tel Telephone
		if param, ok := params["TYPE"]; ok {
			tel.Type = param
		} else {
			tel.Type = []string{"voice"}
		}
		tel.Number = value.GetText()
		current.Telephones = append(current.Telephones, tel)
	case "EMAIL":
		var email Email
		if param, ok := params["TYPE"]; ok {
			email.Type = param
		} else {
			email.Type = []string{"HOME"}
		}
		email.Address = value.GetText()
		current.Emails = append(current.Emails, email)
	case "TITLE":
		current.Title = value.GetText()
	case "ROLE":
		current.Role = value.GetText()
	case "ORG":
		current.Org = value.GetTextList()
	case "CATEGORIES":
		current.Categories = value.GetTextList()
	case "NOTE":
		current.Note = value.GetText()
	case "URL":
		current.URL = value.GetText()
	case "X-JABBER":
	case "X-GTALK":
		var jabber XJabber
		if param, ok := params["TYPE"]; ok {
			jabber.Type = param
		} else {
			jabber.Type = []string{"HOME"}
		}
		jabber.Address = value.GetText()
		current.XJabbers = append(current.XJabbers, jabber)
	case "X-ABShowAs":
		current.XABShowAs = value.GetText()
	/*case "X-ABLabel":
	case "X-ABADR":
		// ignore*/
	default:
		log.Printf("Not read %s, %s: %s\n", group, name, value)
	}
}

func (v *AddressBook) Read(reader io.Reader) {
	di := NewDirectoryInformation(v)
	di.Read(reader)
}

func (v *AddressBook) Write(writer io.Writer) {
	//di := NewDirectoryInformation(v)
	//di.Write(writer)
}
