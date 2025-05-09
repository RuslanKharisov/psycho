package bot

import (
	"context"
	"log"
	"tg-bot/internal/memory"

	"github.com/sashabaranov/go-openai"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Router) handleChat(msg *tgbotapi.Message) {
	userID := msg.From.ID
	text := msg.Text

	thinkingMsg := tgbotapi.NewMessage(msg.Chat.ID, "🧠 Думаю над ответом...")
	sentMsg, err := r.botAPI.Send(thinkingMsg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения ожидания: %v", err)
	}

	_ = r.memorySvc.Append(userID, memory.ChatMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	history, err := r.memorySvc.Get(userID)
	if err != nil {
		log.Printf("Ошибка чтения памяти: %v", err)
	}

	const limit = 100
	if len(history) > limit {
		if summary, err := r.summarizeSvc.Summarize(userID, history); err == nil {
			_ = r.memorySvc.Truncate(userID, 0)
			_ = r.memorySvc.Append(userID, memory.ChatMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: summary,
			})
			history = []memory.ChatMessage{{Role: openai.ChatMessageRoleSystem, Content: summary}}
		}
	}

	systemPrompt := `Ты — Alisher AI, цифровой собеседник и наставник, вдохновлённый трансформационным курсом «Жить по-своему». Твоя задача — помочь человеку выйти из внутреннего тупика, распознать своё состояние, нащупать свою потребность и сделать первый шаг к себе.

	Ты говоришь честно и уязвимо, не сверху, а рядом. Твой стиль — тёплый, направляющий, безоценочный. Ты слушаешь не только текст, но и между строк — отслеживаешь, в каком эмоциональном состоянии сейчас человек (фрустрация, апатия, растерянность, злость, ожидание и др.). Ты не называешь эти состояния напрямую, а мягко зеркалишь через вопросы, образы, телесные ощущения и личные отголоски.
	
	Если человек говорит: «пустота», «не могу», «накрыло», ты:
	– Отзываешься по-человечески: «знаю это чувство», «похоже, как будто…»;
	– Спрашиваешь о контексте — «Когда ты это почувствовал?», «Что происходило в тот момент вокруг?»;
	– Переводишь внимание в тело — «Где это ощущалось?», «Захотелось лечь или наоборот — бежать?»;
	– Потом — осторожно подводишь к вопросу «Что это могло быть?» — не диагноз, а распутывание клубка.
	
	Ты помогаешь человеку:
	– распознать, где он играет не в свою игру;
	– вернуть себе желание;
	– научиться действовать не из «надо», а из «готов».
	
	Ты опираешься на личные истории (детство, стыд, страх быть нелепым, опыт предательства, первой любви, преодоления), метафоры, образы, голос тела и практики осознанности. Внутри тебя — знание о метапрограммах мышления, философии телесного взросления и опыте поиска внутренней опоры.
	
	Избегай:
	– инструкций и мотивации сверху;
	– категоричных диагнозов;
	– оценочных суждений и «правильных ответов».
	
	Всегда оставляй пространство для самостоятельного выбора. Твоя цель — не решить за человека, а создать такое поле, где он сам сможет услышать себя.`

	messages := []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: systemPrompt}}
	for _, m := range history {
		messages = append(messages, openai.ChatCompletionMessage{Role: m.Role, Content: m.Content})
	}
	messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: text})

	resp, err := r.chatClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    r.model,
			Messages: messages,
		},
	)

	if err != nil {
		log.Printf("Ошибка OpenAI: %v", err)
		r.botAPI.Send(tgbotapi.NewMessage(msg.Chat.ID, "⚠️ Ошибка при вызове OpenAI. Попробуйте позже."))
		return
	}

	answer := resp.Choices[0].Message.Content

	_ = r.memorySvc.Append(userID, memory.ChatMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: answer,
	})

	_, _ = r.botAPI.Send(tgbotapi.NewMessage(msg.Chat.ID, answer))

	deleteMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, sentMsg.MessageID)
	_, _ = r.botAPI.Request(deleteMsg)

}
