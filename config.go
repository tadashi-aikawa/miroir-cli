package main

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// Config configuration
type Config struct {
	Bucket       string
	BucketPrefix string
	Table        string
}

// CreateConfig creates configurations from .miroirconfig(toml)
func CreateConfig() (Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return Config{}, errors.Wrap(err, "Home directory is not found.")
	}

	configPath := filepath.Join(home, ".miroirconfig")

	var conf Config
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
