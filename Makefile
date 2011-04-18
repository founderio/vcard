include $(GOROOT)/src/Make.inc

TARG=bitbucket.org/llg/vcard
GOFILES=\
	vcard.go\
	directory_reader.go\
	directory_writer.go\
	address_book.go\
	content_line.go\


include $(GOROOT)/src/Make.pkg
