package main

import (
	"os"
	//"flag"
	"log"
	"bitbucket.org/llg/vcard"
)


func main() {
	//for _, abpath := range flag.Args() {
	var addressBook vcard.AddressBook
	abpath := "contacts.vcf"
	f, err := os.Open(abpath, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		log.Printf("Can't read file %s\n", abpath)
		return
	}
	addressBook.Read(f)
	log.Printf("Read %s\n", abpath)
	//}
}
