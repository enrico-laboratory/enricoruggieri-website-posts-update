package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("successfully loads Notion token", func(t *testing.T) {
		const testToken = "test-secret"
		os.Setenv("NOTION_TOKEN", testToken)
		defer os.Unsetenv("NOTION_TOKEN")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cfg.NotionAPIKey != testToken {
			t.Errorf("expected NotionAPIKey=%q, got %q", testToken, cfg.NotionAPIKey)
		}
	})

	t.Run("fails when NOTION_TOKEN is not set", func(t *testing.T) {
		err := os.Unsetenv("NOTION_TOKEN")

		_, err = Load()
		if err == nil {
			t.Fatal("expected error when NOTION_TOKEN is not set, got nil")
		}
	})
}
