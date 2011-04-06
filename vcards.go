// See RFC 2425 Mime Content-Type for Directory Information
package main

import (
	"fmt"
	"os"
	"scanner"
)


type ValueType int

const (
	TextListType ValueType = iota
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
	c := s.Peek()
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
		s.Next()
		c = s.Peek()
	}
	return
}

func readValues(s *scanner.Scanner) (values []Value) {
	lastChar := s.Next()
	c := lastChar
	var buf []int
	for c != scanner.EOF {
		if lastChar == '\r' && c == '\n' {
			la := s.Peek()
			if la != 32 && la != 9 {
				// call handler and return
				if len(buf) > 0 {
					value := Value{string(buf), TextListType}
					values = append(values, value)
				}
				return
			} else {
				// unfold
				c = s.Next()
			}
		}
		if c != '\n' && c != '\r' && c!= 32 && c != 9 {
			buf = append(buf, c)
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
	s.Next()
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
