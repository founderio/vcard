// See RFC 2425 Mime Content-Type for Directory Information
package vcard

import (
	"fmt"
	"os"
	"scanner"
)

type ContentLineFunc func(group, name string, params map[string]string, values []string)


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

func readValues(s *scanner.Scanner) (values []string) {
	lastChar := s.Next()
	c := lastChar
	var buf []int
	escape := false
	for c != scanner.EOF {
		if c == '\n' {
			la := s.Peek()
			if la != 32 && la != 9 {
				// call handler and return
				if len(buf) > 0 {
					values = append(values, string(buf))
				}
				return
			} else {
				// unfold
				lastChar = la
				c = s.Next()
				for c == 32 || c == 9 {
					c = s.Next()
				}
			}
		}

		if c == '\\' {
			escape = true
		} else if escape {
			if c == 'n' {
				c = '\n'
			}
			buf = append(buf, c)
			escape = false
		} else if c == ';' {
			if len(buf) > 0 {
				values = append(values, string(buf))
				buf = []int{}
			}
		} else if c != '\n' && c != '\r' {
			buf = append(buf, c)
		}
		lastChar = c
		c = s.Next()
	}
	return
}

func readParameters(s *scanner.Scanner) (params map[string]string) {
	lastChar := s.Peek()
	c := lastChar
	var buf []int
	var name string
	var value string
	params = make(map[string]string)
	for c != scanner.EOF {
		if c == ';' || c == ':' {
			if name == "" {
				name = string(buf)
			} else {
				value = string(buf)
			}
			if name != "" {
				params[name] = value
			}
			if c == ':' {
				return
			}
			buf = []int{}
			name = ""
			value = ""
		} else if c == '=' {
			name = string(buf)
			buf = []int{}
		} else {
			buf = append(buf, c)
		}
		s.Next()
		c = s.Peek()
	}
	return
}

func readContentLine(s *scanner.Scanner, handler ContentLineFunc) {
	group, name := readGroupName(s)
	var params map[string]string
	if s.Peek() == ';' {
		params = readParameters(s)
	}
	s.Next()
	values := readValues(s)
	handler(group, name, params, values)
}

func contentLine(group, name string, params map[string]string, values []string) {
	fmt.Println(group, name, params, values)
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
