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

func (t *TgVercel) HandleSetup(w http.ResponseWriter, r *http.Request) {
	o := t.options

	tgBotServiceKey := os.Getenv(o.KeyEnvName)
	if tgBotServiceKey == "" {
		errorResponse(w, fmt.Errorf("%s is not set", o.KeyEnvName))
		return
	}

	vercelEnv := os.Getenv(o.VercelEnvEnvName)
	if vercelEnv == "" {
		errorResponse(w, fmt.Errorf("%s is not set", o.VercelEnvEnvName))
		return
	}

	vercelGeneratedDeploymentUrl := os.Getenv(o.VercelUrlEnvName)
	if vercelGeneratedDeploymentUrl == "" {
		errorResponse(w, fmt.Errorf("%s is not set", o.VercelUrlEnvName))
		return
	}

	vercelProjectProductionUrl := os.Getenv(o.VercelProjectProductionUrlEnvName)
	if vercelProjectProductionUrl == "" {
		errorResponse(w, fmt.Errorf("%s is not set", o.VercelProjectProductionUrlEnvName))
		return
	}

	webhookSecret := os.Getenv(o.TelegramWebhookSecretEnvName)
	if webhookSecret == "" {
		errorResponse(w, fmt.Errorf("%s is not set", o.TelegramWebhookSecretEnvName))
		return
	}

	key := r.URL.Query().Get(o.KeyParamName)
	if key != tgBotServiceKey {
		unauthorizedResponse(w, fmt.Errorf("invalid key"))
		return
	}

	domain := vercelGeneratedDeploymentUrl
	if vercelEnv == "production" {
		domain = vercelProjectProductionUrl
	}

	webhookUri := fmt.Sprintf("https://%s%s?secret=%s",
		domain,
		o.WebhookRelativeUrl,
		webhookSecret,
	)

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
