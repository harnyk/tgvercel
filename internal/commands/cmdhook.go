package commands

import (
	"fmt"
	"log"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/harnyk/tgvercel/internal/vapi"
)

type CmdHookOptions struct {
	VercelToken       string
	DeploymentIDOrURL string
	TelegramBotRoute  string
}

type CmdHook struct {
	Options CmdHookOptions
}

func NewCmdHook(options CmdHookOptions) *CmdHook {
	return &CmdHook{
		Options: options,
	}
}

func (c *CmdHook) Execute() error {
	client := vapi.NewClientWithOptions(vapi.Options{
		Token: c.Options.VercelToken,
	})

	deployment, err := client.GetDeployment(c.Options.DeploymentIDOrURL)
	if err != nil {
		return err
	}

	deploymentDomain := deployment.Url
	deploymentTarget := deployment.Target()

	envTelegramSecret, err := client.GetEnv(deployment.ProjectID, envVarNameTelegramSecret, deployment.Target())
	if err != nil {
		return err
	}

	envTelegramToken, err := client.GetEnv(deployment.ProjectID, envVarNameTelegramToken, deployment.Target())
	if err != nil {
		return err
	}

	log.Printf("Telegram Webhook secret: %s", envTelegramSecret)
	log.Printf("Telegram Webhook token:  %s", envTelegramToken)
	log.Printf("Deployment Domain:       %s", deploymentDomain)
	log.Printf("Bot Route:               %s", c.Options.TelegramBotRoute)
	log.Printf("Target:                  %s", deploymentTarget)

	webhookUrlQuery := url.Values{}
	webhookUrlQuery.Set("secret", envTelegramSecret)
	webhookUrl := url.URL{
		Scheme:   "https",
		Host:     deploymentDomain,
		Path:     c.Options.TelegramBotRoute,
		RawQuery: webhookUrlQuery.Encode(),
	}

	log.Printf("Telegram Webhook URL:    %s", webhookUrl.String())

	bot, err := tgbotapi.NewBotAPI(envTelegramToken)
	if err != nil {
		return err
	}

	wh, err := tgbotapi.NewWebhook(webhookUrl.String())
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}
	resp, err := bot.Request(wh)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}
	if !resp.Ok {
		return fmt.Errorf("failed to set webhook: %s", resp.Description)
	}

	log.Printf("Telegram Webhook set successfully with message: %s", resp.Description)

	return nil
}
