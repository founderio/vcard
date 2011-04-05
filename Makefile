OFILE_EXT=6
GC=$(OFILE_EXT)g
LINK=$(OFILE_EXT)l

all: flip_name


flip_name: vcards.go
	$(GC) vcards.go && $(LINK) -o flip_name vcards.$(OFILE_EXT)

clean:
	rm -f *.8 *.6 flip_name


