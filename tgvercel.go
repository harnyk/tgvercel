package tgvercel

import (
	"encoding/json"
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

func setJsonType(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func errorResponse(w http.ResponseWriter, err error) {
	setJsonType(w)
	w.WriteHeader(http.StatusInternalServerError)
	errorJson := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	jsonData, err := json.Marshal(errorJson)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonData)
}

func unauthorizedResponse(w http.ResponseWriter, err error) {
	setJsonType(w)
	w.WriteHeader(http.StatusUnauthorized)
	errorJson := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	jsonData, err := json.Marshal(errorJson)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonData)
}

func okResponse(w http.ResponseWriter, data interface{}) {
	setJsonType(w)
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonData)
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

func (w *TgVercel) HandleWebhook(r *http.Request, onUpdate UpdateHandlerFunc) {
	bot, err := w.Bot()
	if err != nil {
		log.Fatal(err)
	}
	update, err := bot.HandleUpdate(r)
	if err != nil {
		log.Fatal(err)
	}
	onUpdate(w.bot, update)
}

func (t *TgVercel) HandleSetup(w http.ResponseWriter, r *http.Request) {
	o := t.options

	tgBotServiceKey := os.Getenv(o.KeyEnvName)
	if tgBotServiceKey == "" {
		errorResponse(w, fmt.Errorf("%s is not set", o.KeyEnvName))
		return
	}

	vercelUrl := os.Getenv(o.VercelUrlEnvName)
	if vercelUrl == "" {
		errorResponse(w, fmt.Errorf("%s is not set", o.VercelUrlEnvName))
		return
	}

	key := r.URL.Query().Get(o.KeyParamName)
	if key != tgBotServiceKey {
		unauthorizedResponse(w, fmt.Errorf("invalid key"))
		return
	}

	webhookUri := fmt.Sprintf("https://%s%s", vercelUrl, o.WebhookRelativeUrl)

	bot, err := t.Bot()
	if err != nil {
		errorResponse(w, err)
		return
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(webhookUri)
	if err != nil {
		log.Fatal(err)
	}

	apiResponse, err := bot.Request(wh)
	if err != nil {
		errorResponse(w, fmt.Errorf("failed to set webhook: %w", err))
		return
	}

	if !apiResponse.Ok {
		errorResponse(w, fmt.Errorf("failed to set webhook: %s", apiResponse.Description))
		return
	}

	okResponse(w, apiResponse.Description)
}