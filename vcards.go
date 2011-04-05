package main

import (
	"fmt"
	"os"
	"bufio"
)

func nameValue(name, value string) {
	fmt.Printf("Name: %s, value: %s\n", name, value)
}

func main() {
	f, err := os.Open("../../data/addressBook.vcf", os.O_RDONLY, 0666)
	if err != nil {

		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	c, err := reader.ReadByte()
	var buf []byte
	escape := false
	var name string
	var value string
	for err == nil {
		if escape {
			switch c {
			case 'n':
				c = '\n'
			}
			buf = append(buf, c)
			escape = false
		} else if c == '\\' {
			escape = true
		} else if c == ':' {
			name = string(buf)
			buf = []byte{}
			value = ""
		} else if c == '\n' {
			value = string(buf)
			buf = []byte{}
			nameValue(name, value)
			name, value = "", ""
		} else {
			buf = append(buf, c)
		}
		c, err = reader.ReadByte()
	}
	if len(name) > 0 {
		nameValue(name, value)
	}
}
