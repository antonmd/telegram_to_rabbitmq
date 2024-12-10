package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "gt.devopspoint.io/devopspoint/telegram-bot/internal/config"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

// SendMessageRequest is the JSON payload we send to Telegram's sendMessage endpoint.
type SendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// sendMessageResponse represents the response from the sendMessage method.
type sendMessageResponse struct {
	OK     bool    `json:"ok"`
	Result Message `json:"result"`
}

func NewClient(TelegramBotToken string, TelegramAPIURL string) *Client {
	baseURL := TelegramAPIURL + "/bot" + TelegramBotToken

	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) GetUpdates(offset int, timeout int) ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=%d", c.baseURL, offset, timeout)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result getUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode getUpdates response: %w", err)
	}

	if !result.OK {
		return nil, fmt.Errorf("telegram API returned not OK")
	}
	return result.Result, nil
}

func (c *Client) sendMessage(chatID int64, text string) error {
	return nil
}

// SendMessage sends a text message to a given chat.
