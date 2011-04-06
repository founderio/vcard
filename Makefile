OFILE_EXT=6
GC=$(OFILE_EXT)g
LINK=$(OFILE_EXT)l

all: readab


readab: vcards.go
	$(GC) vcards.go && $(LINK) -o readab vcards.$(OFILE_EXT)

clean:
	rm -f *.8 *.6 readab


