package main

import (
	"log"
	"os"
	"tg-bot/internal/bot"
	"tg-bot/internal/memory"
	"tg-bot/internal/summarizer"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	_ = godotenv.Load()

	tgToken := os.Getenv("TG_TOKEN")
	openaiKey := os.Getenv("OPENAI_KEY")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASSWORD")
	if tgToken == "" || openaiKey == "" || redisAddr == "" {
		log.Fatal("Не заданы TG_TOKEN, OPENAI_KEY или REDIS_ADDR")
	}

	botAPI, err := tgbotapi.NewBotAPI(tgToken)

	if err != nil {
		log.Fatalf("Ошибка создания Telegram-бота: %v", err)
	}

	log.Printf("Бот запущен: @%s", botAPI.Self.UserName)

	mem := memory.NewRedisMemory(redisAddr, redisPass, 0)

	summarizerSvc := summarizer.New(openaiKey)

	oaiclient := openai.NewClient(openaiKey)
	router := bot.NewRouter(botAPI, mem, oaiclient, summarizerSvc)

	router.Run()
}
