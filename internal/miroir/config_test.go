package miroir

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateConfig_LoadsSnakeCaseKeys(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, ".miroirconfig")
	content := `bucket = "mamansoft-miroir"
bucket_prefix = "production"
table = "miroir-summaries"
role_arn = "arn:aws:iam::123456789012:role/test"
s3_endpoint = "http://localhost:3456"
dynamodb_endpoint = "http://localhost:3456"
sts_endpoint = "http://localhost:3456"
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if oldHome == "" {
			_ = os.Unsetenv("HOME")
			return
		}
		_ = os.Setenv("HOME", oldHome)
	}()

	conf, err := CreateConfig()
	if err != nil {
		t.Fatalf("CreateConfig returned error: %v", err)
	}

	if conf.Bucket != "mamansoft-miroir" {
		t.Fatalf("unexpected bucket: %+v", conf)
	}
	if conf.BucketPrefix != "production" {
		t.Fatalf("unexpected bucket prefix: %+v", conf)
	}
	if conf.Table != "miroir-summaries" {
		t.Fatalf("unexpected table: %+v", conf)
	}
	if conf.RoleARN != "arn:aws:iam::123456789012:role/test" {
		t.Fatalf("unexpected role arn: %+v", conf)
	}
	if conf.S3Endpoint != "http://localhost:3456" {
		t.Fatalf("unexpected s3 endpoint: %+v", conf)
	}
	if conf.DynamoDBEndpoint != "http://localhost:3456" {
		t.Fatalf("unexpected dynamodb endpoint: %+v", conf)
	}
	if conf.STSEndpoint != "http://localhost:3456" {
		t.Fatalf("unexpected sts endpoint: %+v", conf)
	}
}

func TestCreateConfig_ReturnsConfigNotFoundWhenMissing(t *testing.T) {
	tmpDir := t.TempDir()

	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if oldHome == "" {
			_ = os.Unsetenv("HOME")
			return
		}
		_ = os.Setenv("HOME", oldHome)
	}()

	_, err := CreateConfig()
	if err != ErrorConfigIsNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfigErrorsHaveExpectedMessages(t *testing.T) {
	if ErrorHomeDirIsNotFound.Error() != errorMsgHomeDirIsNotFound {
		t.Fatalf("unexpected home dir error: %v", ErrorHomeDirIsNotFound)
	}

	if ErrorConfigIsNotFound.Error() != errorMsgConfigIsNotFound {
		t.Fatalf("unexpected config error: %v", ErrorConfigIsNotFound)
	}
}
