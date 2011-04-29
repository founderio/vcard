package main

import (
	"os"
	"io"
	"bufio"
	"flag"
	"log"
	"bitbucket.org/llg/vcard"
	"fmt"
)

func ask(msg string, format ...interface{}) string {
	fmt.Printf(msg + "\n", format...)
	stdin := bufio.NewReader (os.Stdin)
	line, _, _  := stdin.ReadLine()
	if  line != nil {
		return string(line)	
	}
	return ""
}

func integrateAdditionalName(ab *vcard.AddressBook) {
	for _, contact := range ab.Contacts {
		if len(contact.AdditionalNames) > 0 {
			fmt.Println("---------------------------")
			fmt.Printf("%v\n", contact)
			msg := "Integrate Additional Name in \n"
			if len(contact.GivenNames) > 0 {
				msg +=  "given name (g)?\n"
			}
			if len(contact.FamilyNames) > 0 {
				msg +=  "family name (f) ?\n"
			}
			switch ask(msg) {
				case "g":
				contact.GivenNames[0] += " " + contact.AdditionalNames[0]
				fmt.Printf("result: %s\n", contact.GivenNames[0])
				case "f":
				contact.FamilyNames[0] += " " + contact.AdditionalNames[0]
				fmt.Printf("result: %s\n", contact.FamilyNames[0])
			}
			
		}
	}
}

func main() {
	var output io.Writer
	var err os.Error
	outputFilename := flag.String("o", "", "Output vcard file")
	flag.Parse()
	if *outputFilename == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(*outputFilename)
		output = bufio.NewWriter(output)
		if err != nil {
			log.Printf("Can't create %s\n", *outputFilename)
			return
		}
	}
	var args []string
	if len(flag.Args()) > 0 {
		args = flag.Args()
	} else {
		args = []string{"contacts.vcf"}
	}
	var addressBook vcard.AddressBook
	for _, abpath := range args {
		f, err := os.Open(abpath)
		defer f.Close()
		if err != nil {
			log.Printf("Can't read file %s\n", abpath)
			return
		}
		reader := vcard.NewDirectoryInfoReader(f)
		addressBook.ReadFrom(reader)
		log.Printf("Read %s\n", abpath)
	}

	integrateAdditionalName(&addressBook)
	writer := vcard.NewDirectoryInfoWriter(output)
	addressBook.WriteTo(writer)
	log.Printf("Write %s\n", *outputFilename)
}
