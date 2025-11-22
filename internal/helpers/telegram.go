package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	SuccessAndUpdate = "[SUCCESS] [UPDATE] The Website Update run was successful: the website was updated"
	SuccessNoUpdate  = "[SUCCESS] [NO UPDATE] The Website Update run was successful: the website did not need to be updated"
	errorMessage     = "[ERROR] The website was not updated."
	StartUpdate      = "[START] The Website Update run was started."
)

type Telegram struct {
	token  string
	chatId string
}

type messageRequest struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

func NewTelegramClient(token, chatId string) (*Telegram, error) {
	return &Telegram{token: token, chatId: chatId}, nil
}

func (t *Telegram) SendMessage(message string) (string, error) {
	message = "From Web Site Update APP: " + message
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)
	body, _ := json.Marshal(messageRequest{
		ChatId: t.chatId,
		Text:   message,
	})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("could not read response body: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("could not send message to : " + string(bodyBytes))
	}

	return string(bodyBytes), err
}

func (t *Telegram) SendError(errMessage string) (string, error) {

	message := fmt.Sprintf("%s\n%s", errorMessage, errMessage)

	resp, err := t.SendMessage(message)
	if err != nil {
		return "", err
	}
	return resp, nil
}
