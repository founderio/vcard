package vcard

import "strings"

type ContentLine struct {
	Group, Name string
	Params      map[string]Value
	Value       StructuredValue
}

// values separated by ';' has a structural meaning
type StructuredValue []Value

// values seprated by ',' is a multi value
type Value []string

func (sv StructuredValue) GetTextList() []string {
	var textList []string
	for _, v := range sv {
		for _, s := range v {
			textList = append(textList, s)
		}
	}
	return textList
}

func (v StructuredValue) GetText() string {
	if len(v) > 0 && len(v[0]) > 0 {
		return v[0][0]
	}
	return ""
}

func (v StructuredValue) GetAllText() string {
	if len(v) > 0 {
		text := []string{}
		for _, v0 := range v {
			if len(v0) > 0 {
				text = append(text, v0[0])
			}
		}
		return strings.Join(text, ";")
	}
	return ""
}

func (v Value) GetText() string {
	if len(v) > 0 {
		return v[0]
	}
	return ""
}
