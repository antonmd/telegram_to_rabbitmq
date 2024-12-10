package telegram

// Update represents an incoming update from Telegram.
type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

// Message represents a message sent to the bot.
type Message struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text,omitempty"`
	Chat      *Chat  `json:"chat,omitempty"`
}

// Chat represents a chat (group, private, etc.).
type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// getUpdatesResponse represents the top-level response from the Telegram API.
type getUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}
