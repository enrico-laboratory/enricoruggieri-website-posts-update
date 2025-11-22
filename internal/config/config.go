package config

import (
	"github.com/enrico-laboratory/website-update/internal/helpers"
)

// Config holds all environment variables
type Config struct {
	NotionAPIKey   string
	GitPAT         string
	TelegramToken  string
	TelegramChatID string
}

// Load reads environment variables and returns a Config
func Load() (*Config, error) {
	notionTokenEnvName := "NOTION_TOKEN"
	notionTokenEnvPath := "NOTION_TOKEN_PATH"
	notionToken, err := helpers.SetToken(notionTokenEnvName, notionTokenEnvPath)
	if err != nil {
		return nil, err
	}

	gitTokenEnvName := "GIT_PAT"
	gitTokenEnvPath := "GIT_PAT_PATH"
	gitToken, err := helpers.SetToken(gitTokenEnvName, gitTokenEnvPath)
	if err != nil {
		return nil, err
	}

	telegramTokenEnvName := "TELEGRAM_TOKEN"
	telegramTokenEnvPath := "TELEGRAM_TOKEN_PATH"
	telegramToken, err := helpers.SetToken(telegramTokenEnvName, telegramTokenEnvPath)
	if err != nil {
		return nil, err
	}

	telegramChatIDEnvName := "TELEGRAM_CHAT_ID"
	telegramChatIDEnvPath := "TELEGRAM_CHAT_ID_PATH"
	telegramChatID, err := helpers.SetToken(telegramChatIDEnvName, telegramChatIDEnvPath)
	if err != nil {
		return nil, err
	}
	return &Config{
		NotionAPIKey:   notionToken,
		GitPAT:         gitToken,
		TelegramToken:  telegramToken,
		TelegramChatID: telegramChatID,
	}, nil
}
