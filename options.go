package tgvercel

import (
	"fmt"
	"strings"
)

type Options struct {
	WebhookRelativeUrl                string
	TelegramTokenEnvName              string
	TelegramWebhookSecretEnvName      string
	VercelEnvEnvName                  string
	VercelUrlEnvName                  string
	VercelProjectProductionUrlEnvName string
	KeyEnvName                        string
	KeyParamName                      string
}

func DefaultOptions() Options {
	return Options{
		WebhookRelativeUrl:                "/api/tg/webhook",
		TelegramTokenEnvName:              "TELEGRAM_TOKEN",
		TelegramWebhookSecretEnvName:      "TELEGRAM_WEBHOOK_SECRET",
		VercelEnvEnvName:                  "VERCEL_ENV",
		VercelUrlEnvName:                  "VERCEL_URL",
		VercelProjectProductionUrlEnvName: "VERCEL_PROJECT_PRODUCTION_URL",
		KeyEnvName:                        "TGVERCEL_KEY",
		KeyParamName:                      "key",
	}
}

func (o *Options) Validate() error {
	if o.WebhookRelativeUrl == "" {
		return fmt.Errorf("WebhookRelativeUrl must be set")
	}
	if !strings.HasPrefix(o.WebhookRelativeUrl, "/") {
		return fmt.Errorf("WebhookRelativeUrl must start with /")
	}

	if o.TelegramTokenEnvName == "" {
		return fmt.Errorf("TelegramTokenEnvName must be set")
	}

	if o.TelegramWebhookSecretEnvName == "" {
		return fmt.Errorf("TelegramWebhookSecretEnvName must be set")
	}

	if o.VercelEnvEnvName == "" {
		return fmt.Errorf("VercelEnvEnvName must be set")
	}

	if o.VercelUrlEnvName == "" {
		return fmt.Errorf("VercelUrlEnvName must be set")
	}

	if o.VercelProjectProductionUrlEnvName == "" {
		return fmt.Errorf("VercelProjectProductionUrlEnvName must be set")
	}

	if o.KeyEnvName == "" {
		return fmt.Errorf("KeyEnvName must be set")
	}

	if o.KeyParamName == "" {
		return fmt.Errorf("KeyParamName must be set")
	}
	return nil
}
