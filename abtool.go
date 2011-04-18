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
	abpath := "addressBook.vcf"
	f, err := os.Open(abpath, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		log.Printf("Can't read file %s\n", abpath)
		return
	}
	reader := vcard.NewDirectoryInfoReader(f)
	addressBook.ReadFrom(reader)
	
	writer := vcard.NewDirectoryInfoWriter(os.Stdout)
	addressBook.WriteTo(writer)
	log.Printf("Read %s\n", abpath)
	//}
}
