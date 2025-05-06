package bot

import (
	"tg-bot/internal/memory"
	"tg-bot/internal/summarizer"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

type Router struct {
	botAPI       *tgbotapi.BotAPI
	memorySvc    memory.Memory
	chatClient   *openai.Client
	summarizeSvc summarizer.Summarizer
}

func NewRouter(
	botAPI *tgbotapi.BotAPI,
	mem memory.Memory,
	chatClient *openai.Client,
	sum summarizer.Summarizer,
) *Router {
	return &Router{
		botAPI:       botAPI,
		memorySvc:    mem,
		chatClient:   chatClient,
		summarizeSvc: sum,
	}
}

func (r *Router) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := r.botAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			go r.handleChat(update.Message)
		}
	}
}
