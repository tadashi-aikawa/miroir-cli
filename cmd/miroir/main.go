package main

import (
	"log"
	"os"

	"github.com/tadashi-aikawa/miroir-cli/internal/miroir"
)

func main() {
	parser, err := miroir.Execute(os.Args[1:])
	if err != nil {
		if parser != nil && miroir.IsParseError(err) {
			parser.FatalIfErrorf(err)
		}
		log.Fatal(err)
	}
}
