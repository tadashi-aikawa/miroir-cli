package miroir

import (
	"errors"
	"log"

	"github.com/alecthomas/kong"
)

func Execute(argv []string) (*kong.Kong, error) {
	_, parser, ctx, err := ParseCLI(argv)
	if err != nil {
		return parser, err
	}

	config, err := CreateConfig()
	if err != nil {
		switch err {
		case ErrorHomeDirIsNotFound:
			log.Printf("[WARN] Home directory is not found and can't load .miroirconfig.... continue...")
		case ErrorConfigIsNotFound:
			log.Printf("[WARN] .miroirconfig is not found.... continue...")
		default:
			return parser, err
		}
	}

	if err := ctx.Run(&appContext{Config: config}); err != nil {
		return parser, err
	}

	return parser, nil
}

func IsParseError(err error) bool {
	var parseErr *kong.ParseError
	return errors.As(err, &parseErr)
}
