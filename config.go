package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

const (
	errorMsgHomeDirIsNotFound = "Home directory is not found"
	errorMsgConfigIsNotFound  = ".miroirconfig is not found"
)

var (
	ErrorHomeDirIsNotFound = errors.New(errorMsgHomeDirIsNotFound)
	ErrorConfigIsNotFound  = errors.New(errorMsgConfigIsNotFound)
)

// Config configuration
type Config struct {
	Bucket           string `toml:"bucket"`
	BucketPrefix     string `toml:"bucket_prefix"`
	Table            string `toml:"table"`
	RoleARN          string `toml:"role_arn"`
	S3Endpoint       string `toml:"s3_endpoint"`
	DynamoDBEndpoint string `toml:"dynamodb_endpoint"`
	STSEndpoint      string `toml:"sts_endpoint"`
}

// CreateConfig creates configurations from .miroirconfig(toml)
func CreateConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, ErrorHomeDirIsNotFound
	}

	configPath := filepath.Join(home, ".miroirconfig")
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return Config{}, ErrorConfigIsNotFound
		}
		return Config{}, errors.Wrap(err, "Fail to inspect `.miroirconfig`.")
	}

	var conf Config
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
