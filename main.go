package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

func main() {
	args, err := CreateArgs(usage, os.Args[1:], version)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to create arguments."))
	}

	switch true {
	case args.CmdGet:
		switch true {
		case args.CmdSummaries:
			if err := CmdGetSummaries(args); err != nil {
				log.Fatal(errors.Wrap(err, "Fail to command `get summaries`"))
			}
		case args.CmdReport:
			if err := CmdGetReport(args); err != nil {
				log.Fatal(errors.Wrap(err, "Fail to command `get report`"))
			}
		}
	}
}
