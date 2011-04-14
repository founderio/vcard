// See RFC 2425 Mime Content-Type for Directory Information
package vcard

import (
	"io"
	"scanner"
)

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

// values separated by ';' has a structural meaning
type StructuredValue []Value

// values seprated by ',' is a multi value
type Value []string

func readValues(s *scanner.Scanner) (value StructuredValue) {
	lastChar := s.Next()
	c := lastChar
	var buf []int
	escape := false
	var val Value
	for c != scanner.EOF {
		if c == '\n' {
			la := s.Peek()
			if la != 32 && la != 9 {
				// return
				if len(buf) > 0 {
					val = append(val, string(buf))
				}
				value = append(value, val)
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
			if c == 'n' || c == 'N' {
				c = '\n'
			}
			buf = append(buf, c)
			escape = false
		} else if c == ',' {
			if len(buf) > 0 {
				val = append(val, string(buf))
				buf = []int{}
			}
		} else if c == ';' {
			if len(buf) > 0 {
				val = append(val, string(buf))
				buf = []int{}
			}
			value = append(value, val)
			val = Value{}
		} else if c != '\n' && c != '\r' {
			buf = append(buf, c)
		}
		lastChar = c
		c = s.Next()
	}
	return
}

func readParameters(s *scanner.Scanner) (params map[string]Value) {
	lastChar := s.Peek()
	c := lastChar
	var buf []int
	var name string
	var value string
	params = make(map[string]Value)
	var values Value
	for c != scanner.EOF {
		if c == ',' {
			values = append(values, string(buf))
			buf = []int{}
		} else if c == ';' || c == ':' {
			if name == "" {
				name = string(buf)
			} else {
				value = string(buf)
			}
			if name != "" {
				values = append(values, value)
				if _, ok := params[name]; ok {
					params[name] = append(params[name], values...)
				} else {
					params[name] = values
				}
			}
			if c == ':' {
				return
			}
			buf = []int{}
			values = Value{}
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

type ContentLineFunc func(group, name string, params map[string]Value, value StructuredValue)

func ReadContentLine(s *scanner.Scanner, handler ContentLineFunc) {
	group, name := readGroupName(s)
	var params map[string]Value
	if s.Peek() == ';' {
		params = readParameters(s)
	}
	s.Next()
	value := readValues(s)
	handler(group, name, params, value)
}

func ReadDirectoryInformation(reader io.Reader, contentLine ContentLineFunc) {
	var s scanner.Scanner
	s.Init(reader)
	for s.Peek() != scanner.EOF {
		ReadContentLine(&s, contentLine)
	}
}
