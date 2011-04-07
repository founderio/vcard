package main

import (
	"os"
	"fmt"
	"bitbucket.org/llg/govcard"
)

func contentLine(group, name string, params map[string]string, values []string) {
	fmt.Println(group, name, params, values)
}

func main() {
	f, err := os.Open("../../data/addressBook.vcf", os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	govcard.ReadDirectoryInformation(f, contentLine)
}
