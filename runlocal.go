package tgvercel

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RunLocal(token string, onUpdate UpdateHandlerFunc) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	w, err := tgbotapi.NewWebhook("")
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}
	_, err = bot.Request(w)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		onUpdate(bot, &update)
	}

	return nil
}
