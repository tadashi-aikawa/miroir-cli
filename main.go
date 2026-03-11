package main

import (
	"log"
	"os"
)

func main() {
	_, parser, ctx, err := ParseCLI(os.Args[1:])
	if err != nil {
		if parser == nil {
			log.Fatal(err)
		}
		parser.FatalIfErrorf(err)
	}

	config, err := CreateConfig()
	if err != nil {
		switch err {
		case ErrorHomeDirIsNotFound:
			log.Printf("[WARN] Home directory is not found and can't load .miroirconfig.... continue...")
		case ErrorConfigIsNotFound:
			log.Printf("[WARN] .miroirconfig is not found.... continue...")
		default:
			log.Fatal(err)
		}
	}

	if err := ctx.Run(&appContext{Config: config}); err != nil {
		log.Fatal(err)
	}
}
