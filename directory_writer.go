package vcard

import (
	"io"
)

type DirectoryInfoWriter struct {
	writer io.Writer
}

func NewDirectoryInfoWriter(writer io.Writer) *DirectoryInfoWriter {
	return &DirectoryInfoWriter{writer}
}

func (di *DirectoryInfoWriter) WriteContentLine(contentLine *ContentLine) {
	if len(contentLine.Value) == 0 {
		return
	}
	if contentLine.Group != "" {
		io.WriteString(di.writer, contentLine.Group)
		io.WriteString(di.writer, ".")
	}
	io.WriteString(di.writer, contentLine.Name)
	if contentLine.Params != nil {
		for key, values := range contentLine.Params {
			io.WriteString(di.writer, key)
			io.WriteString(di.writer, "=")
			for vi := 0; vi < len(values); vi++ {
				io.WriteString(di.writer, values[vi])
				if vi+1 < len(values) {
					io.WriteString(di.writer, ",")
				}
			}
		}
	}
	io.WriteString(di.writer, ":")
	for si := 0; si < len(contentLine.Value); si++ {
		for vi := 0; vi < len(contentLine.Value[si]); vi++ {
			di.WriteValue(contentLine.Value[si][vi])
			if vi+1 < len(contentLine.Value[si]) {
				io.WriteString(di.writer, ",")
			}
		}
		if si+1 < len(contentLine.Value) {
			io.WriteString(di.writer, ";")
		}
	}
	io.WriteString(di.writer, "\r\n")
}

func (di *DirectoryInfoWriter) WriteValue(value string) {
	// TODO escape characters '\n' ';' ',' and fold long lines
	io.WriteString(di.writer, value)
}
