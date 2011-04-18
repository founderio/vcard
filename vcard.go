package vcard

import (
	"log"
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

type DataType interface {
	GetType() []string
	HasType(t string) bool
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

func (vcard *VCard) Read(di *DirectoryInfoReader) {
	contentLine := di.ReadContentLine()
	for contentLine != nil {
		switch contentLine.Name {
		case "VERSION":
			vcard.Version = contentLine.Value.GetText()
		case "END":
			if contentLine.Value.GetText() == "VCARD" {
				return
			}
		case "FN":
			if vcard != nil {
				vcard.FormattedName = contentLine.Value.GetText()
			}
		case "N":
			if len(contentLine.Value) == nameSize {
				vcard.FamilyNames = contentLine.Value[familyNames]
				vcard.GivenNames = contentLine.Value[givenNames]
				vcard.AdditionalNames = contentLine.Value[additionalNames]
				vcard.HonorificNames = contentLine.Value[honorificPrefixes]
				vcard.HonorificSuffixes = contentLine.Value[honorificSuffixes]
			} else {
				log.Printf("Error structured data isn't appropriate: %d\n", len(contentLine.Value))
			}
		case "NICKNAME":
			vcard.NickNames = contentLine.Value.GetTextList()
		case "PHOTO":
			vcard.Photo.Encoding = contentLine.Params["ENCODING"].GetText()
			vcard.Photo.Type = contentLine.Params["TYPE"].GetText()
			vcard.Photo.Value = contentLine.Params["VALUE"].GetText()
			vcard.Photo.Data = contentLine.Value.GetText()
		case "BDAY":
			vcard.Birthday = contentLine.Value.GetText()
		case "ADR":
			if len(contentLine.Value) == addressSize {
				var address Address
				if param, ok := contentLine.Params["TYPE"]; ok {
					address.Type = param
				} else {
					address.Type = defaultAddressTypes()
				}
				address.PostOfficeBox = contentLine.Value[postOfficeBox].GetText()
				address.ExtendedAddress = contentLine.Value[extendedAddress].GetText()
				address.Street = contentLine.Value[street].GetText()
				address.Locality = contentLine.Value[locality].GetText()
				address.Region = contentLine.Value[region].GetText()
				address.PostalCode = contentLine.Value[postalCode].GetText()
				address.CountryName = contentLine.Value[countryName].GetText()
				vcard.Addresses = append(vcard.Addresses, address)
			} else {
				log.Printf("Error structured data isn't appropriate: %d\n", len(contentLine.Value))
			}
		case "X-ABUID":
			vcard.XABuid = contentLine.Value.GetText()
		case "TEL":
			var tel Telephone
			if param, ok := contentLine.Params["TYPE"]; ok {
				tel.Type = param
			} else {
				tel.Type = []string{"voice"}
			}
			tel.Number = contentLine.Value.GetText()
			vcard.Telephones = append(vcard.Telephones, tel)
		case "EMAIL":
			var email Email
			if param, ok := contentLine.Params["TYPE"]; ok {
				email.Type = param
			} else {
				email.Type = []string{"HOME"}
			}
			email.Address = contentLine.Value.GetText()
			vcard.Emails = append(vcard.Emails, email)
		case "TITLE":
			vcard.Title = contentLine.Value.GetText()
		case "ROLE":
			vcard.Role = contentLine.Value.GetText()
		case "ORG":
			vcard.Org = contentLine.Value.GetTextList()
		case "CATEGORIES":
			vcard.Categories = contentLine.Value.GetTextList()
		case "NOTE":
			vcard.Note = contentLine.Value.GetText()
		case "URL":
			vcard.URL = contentLine.Value.GetText()
		case "X-JABBER":
		case "X-GTALK":
			var jabber XJabber
			if param, ok := contentLine.Params["TYPE"]; ok {
				jabber.Type = param
			} else {
				jabber.Type = []string{"HOME"}
			}
			jabber.Address = contentLine.Value.GetText()
			vcard.XJabbers = append(vcard.XJabbers, jabber)
		case "X-ABShowAs":
			vcard.XABShowAs = contentLine.Value.GetText()
		/*case "X-ABLabel":
		case "X-ABADR":
			// ignore*/
		default:
			log.Printf("Not read %s, %s: %s\n", contentLine.Group, contentLine.Name, contentLine.Value)
		}
		contentLine = di.ReadContentLine()
	}
}

func (vcard *VCard) Write(di *DirectoryInfoWriter) {
	di.WriteContentLine(&ContentLine{"", "BEGIN", nil, StructuredValue{Value{"VCARD"}}})
	di.WriteContentLine(&ContentLine{"", "VERSION", nil, StructuredValue{Value{"3.0"}}})
	di.WriteContentLine(&ContentLine{"", "FN", nil, StructuredValue{Value{vcard.FormattedName}}})
	di.WriteContentLine(&ContentLine{"", "N", nil, StructuredValue{vcard.FamilyNames, vcard.GivenNames, vcard.AdditionalNames, vcard.HonorificNames, vcard.HonorificSuffixes}})
	di.WriteContentLine(&ContentLine{"", "NICKNAME", nil, StructuredValue{vcard.NickNames}})
	vcard.Photo.Write(di)
	di.WriteContentLine(&ContentLine{"", "BDAY", nil, StructuredValue{Value{vcard.Birthday}}})		
	for _, addr := range vcard.Addresses {
		addr.Write(di)
	}
	di.WriteContentLine(&ContentLine{"", "X-ABUID", nil, StructuredValue{Value{vcard.XABuid}}})
	for _, tel := range vcard.Telephones {
		tel.Write(di)
	}
	for _, email := range vcard.Emails {
		email.Write(di)
	}
	di.WriteContentLine(&ContentLine{"", "TITLE", nil, StructuredValue{Value{vcard.Title}}})
	di.WriteContentLine(&ContentLine{"", "ROLE", nil, StructuredValue{Value{vcard.Role}}})
	di.WriteContentLine(&ContentLine{"", "ORG", nil, StructuredValue{vcard.Org}})
	di.WriteContentLine(&ContentLine{"", "CATEGORIES", nil, StructuredValue{vcard.Categories}})
	di.WriteContentLine(&ContentLine{"", "NOTE", nil, StructuredValue{Value{vcard.Note}}})
	di.WriteContentLine(&ContentLine{"", "URL", nil, StructuredValue{Value{vcard.URL}}})
	for _, jab := range vcard.XJabbers {
		jab.Write(di)
	}
	di.WriteContentLine(&ContentLine{"", "X-ABShowAs", nil, StructuredValue{Value{vcard.XABShowAs}}})
	di.WriteContentLine(&ContentLine{"", "END", nil, StructuredValue{Value{"VCARD"}}})
}

func (photo *Photo) Write(di *DirectoryInfoWriter) {
	params := make(map[string]Value)
	if photo.Encoding != "" {
		params["ENCODING"] = Value{photo.Encoding}
	}
	if photo.Type != "" {
		params["TYPE"] = Value{photo.Type}
	}
	if photo.Value != "" {
		params["VALUE"] = Value{photo.Value}
	}
	di.WriteContentLine(&ContentLine{"", "PHOTO", params, StructuredValue{Value{photo.Data}}})
}

func (addr *Address) Write(di *DirectoryInfoWriter) {
	params := make(map[string]Value)
	params["TYPE"] = addr.Type
	di.WriteContentLine(&ContentLine{"", "ADR", params, StructuredValue{Value{addr.PostOfficeBox}, Value{addr.ExtendedAddress}, Value{addr.Street}, Value{addr.Locality}, Value{addr.Region}, Value{addr.PostalCode}, Value{addr.CountryName}}})
}

func (tel *Telephone) Write(di *DirectoryInfoWriter) {
	params := make(map[string]Value)
	params["TYPE"] = tel.Type
	di.WriteContentLine(&ContentLine{"", "TEL", params, StructuredValue{Value{tel.Number}}})
}

func (email *Email) Write(di *DirectoryInfoWriter) {
	params := make(map[string]Value)
	params["TYPE"] = email.Type
	di.WriteContentLine(&ContentLine{"", "EMAIL", params, StructuredValue{Value{email.Address}}})
}

func (jab *XJabber) Write(di *DirectoryInfoWriter) {
	params := make(map[string]Value)
	params["TYPE"] = jab.Type
	di.WriteContentLine(&ContentLine{"", "EMAIL", params, StructuredValue{Value{jab.Address}}})
}
