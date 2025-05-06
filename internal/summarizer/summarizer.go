package summarizer

import "tg-bot/internal/memory"

type Summarizer interface {
	Summarize(userID int64, history []memory.ChatMessage) (string, error)
}
