package main

import (
	"os"
	"flag"
	"fmt"
	"bitbucket.org/llg/govcard"
)

func contentLine(group, name string, params map[string]string, values []string) {
	fmt.Println(group, name, params, values)
}

func main() {
	for _, abpath := range flag.Args() {
		f, err := os.Open(abpath, os.O_RDONLY, 0666)
		defer f.Close()
		if err != nil {
			return
		}
		govcard.ReadDirectoryInformation(f, contentLine)
	}
}
