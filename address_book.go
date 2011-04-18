package vcard

import (
		"log"
		"io"
)

type AddressBook struct {
	Contacts []VCard
}



func (ab *AddressBook) LastContact() *VCard {
	if len(ab.Contacts) > 0 {
		return &ab.Contacts[len(ab.Contacts)-1]
	}
	return nil
}

func (ab *AddressBook) Read(di *DirectoryInfoReader) {
	contentLine := di.NextContentLine()
	for contentLine != nil {
		switch (contentLine.Name) {
		case "BEGIN":
			if contentLine.Value.GetText() == "VCARD" {
				var vcard VCard
				vcard.Read(di)
				ab.Contacts = append(ab.Contacts, vcard)
			}
		default:
			log.Printf("Not read %s, %s: %s\n", contentLine.Group, contentLine.Name, contentLine.Value)
		}
		contentLine = di.NextContentLine()
	}
}

func (v *AddressBook) Write(writer io.Writer) {
	//di := NewDirectoryInformation(v)
	//di.Write(writer)
}
