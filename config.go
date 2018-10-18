package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const (
	errorMsgHomeDirIsNotFound = "Home directory is not found"
	errorMsgConfigIsNotFound  = ".miroirconfig is not found"
)

var (
	ErrorHomeDirIsNotFound = errors.New(errorMsgConfigIsNotFound)
	ErrorConfigIsNotFound  = errors.New(errorMsgHomeDirIsNotFound)
)

// Config configuration
type Config struct {
	Bucket       string
	BucketPrefix string
	Table        string
	RoleARN      string
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// CreateConfig creates configurations from .miroirconfig(toml)
func CreateConfig() (Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return Config{}, ErrorHomeDirIsNotFound
	}

	configPath := filepath.Join(home, ".miroirconfig")
	if isExists := exists(configPath); !isExists {
		return Config{}, ErrorConfigIsNotFound
	}

	var conf Config
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
