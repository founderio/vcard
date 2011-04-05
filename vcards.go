// See RFC 2425 Mime Content-Type for Directory Information
package main

import (
	"fmt"
	"os"
	"scanner"
)


type ValueType int

const (
	textListType ValueType = iota
	GenericurlType
	DateListType
	TimeListType
	DateTimeListType
	BooleanType
	IntegerListType
	FloatListType
	IanaValueSpecType
)

type Value struct {
	Data string
	Type ValueType
}

type ContentLineFunc func(group, name string, params map[string]string, values []Value)

func contentLine(group, name string, params map[string]string, values []Value) {
	fmt.Println(group, name, params, values)
}

func readGroupName(s *scanner.Scanner) (group, name string) {
	c := s.Next()
	var buf []int
	for c != scanner.EOF {
		if c == '.' {
			group = string(buf)
			buf = []int{}
		} else if c == ';' || c == ':' {
			name = string(buf)
			return
		} else {
			buf = append(buf, c)
		}
		c = s.Next()
	}
	return
}

func readValues(s *scanner.Scanner) (values []Value) {
	lastChar := s.Peek()
	c := lastChar
	isCRLF := false
	for c != scanner.EOF {
		if lastChar == '\r' && c == '\n' {
			isCRLF = true
		}
		if isCRLF {
			if c == 32 || c == 9 {
				// unfold
				isCRLF = false
			} else {
				// call handler and return
				return
			}
		}
		lastChar = c
		c = s.Next()
	}
	return
}


func readContentLine(s *scanner.Scanner, handler ContentLineFunc) {
	group, name := readGroupName(s)
	var params map[string]string
	if s.Peek() == ';' {
		//params = readParameters(s)
	}
	values := readValues(s)
	handler(group, name, params, values)
}

func main() {
	f, err := os.Open("../../data/addressBook.vcf", os.O_RDONLY, 0666)
	if err != nil {

		return
	}
	defer f.Close()
	var s scanner.Scanner
	s.Init(f)
	for s.Peek() != scanner.EOF {
		readContentLine(&s, contentLine)
	}
}
