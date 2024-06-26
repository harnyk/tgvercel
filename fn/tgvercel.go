package fn

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgVercel struct {
	options Options
	bot     *tgbotapi.BotAPI
}

type UpdateHandlerFunc func(bot *tgbotapi.BotAPI, update *tgbotapi.Update)

func New(options Options) *TgVercel {
	err := options.Validate()
	if err != nil {
		log.Fatal(err)
	}

	return &TgVercel{
		options: options,
	}
}

func (w *TgVercel) Bot() (*tgbotapi.BotAPI, error) {
	if w.bot == nil {
		token := os.Getenv(w.options.TelegramTokenEnvName)
		if token == "" {
			return nil, fmt.Errorf("%s is not set", w.options.TelegramTokenEnvName)
		}

		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			return nil, err
		}
		w.bot = bot
	}
	return w.bot, nil
}

func (t *TgVercel) HandleWebhook(r *http.Request, onUpdate UpdateHandlerFunc) {
	o := t.options

	webhookSecret := os.Getenv(o.TelegramWebhookSecretEnvName)
	if webhookSecret == "" {
		log.Fatal(fmt.Errorf("%s is not set", o.TelegramWebhookSecretEnvName))
		return
	}

	secret := r.URL.Query().Get("secret")
	if secret != webhookSecret {
		log.Fatal(fmt.Errorf("invalid secret"))
		return
	}

	bot, err := t.Bot()
	if err != nil {
		log.Fatal(err)
	}
	update, err := bot.HandleUpdate(r)
	if err != nil {
		log.Fatal(err)
	}
	onUpdate(t.bot, update)
}
