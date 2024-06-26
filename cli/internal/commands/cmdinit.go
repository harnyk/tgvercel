package commands

import (
	"log"

	"github.com/google/uuid"
	"github.com/harnyk/tgvercel/cli/internal/vapi"
	"github.com/harnyk/tgvercel/cli/internal/vconfig"
)

type CmdInitOptions struct {
	Target          string
	VercelToken     string
	TelegramToken   string
	TelegramSecret  string
	ProjectIDOrName string
}

type CmdInit struct {
	Options CmdInitOptions
}

func NewCmdInit(options CmdInitOptions) *CmdInit {
	return &CmdInit{
		Options: options,
	}
}

func (c *CmdInit) Execute() error {
	localConfig := vconfig.NewConfig()

	token := c.Options.VercelToken
	if token == "" {
		localToken, err := localConfig.GetAuthToken()
		if err != nil {
			return err
		}
		token = localToken
	}

	projectIDOrName := c.Options.ProjectIDOrName
	if projectIDOrName == "" {
		localProjectID, err := localConfig.GetProjectId()
		if err != nil {
			return err
		}
		projectIDOrName = localProjectID
	}

	target, err := vapi.NewTarget(c.Options.Target)
	if err != nil {
		return err
	}

	client := vapi.NewClientWithOptions(vapi.Options{
		Token: token,
	})

	telegramWebhookSecret := c.Options.TelegramSecret
	if telegramWebhookSecret == "" {
		telegramWebhookSecret = uuid.NewString()
	}

	log.Printf("Setting environment variable %s of project %s target %s", envVarNameTelegramToken, projectIDOrName, target)
	err = client.SetEnv(projectIDOrName, envVarNameTelegramToken, c.Options.TelegramToken, target)
	if err != nil {
		return err
	}

	log.Printf("Setting environment variable %s of project %s target %s", envVarNameTelegramSecret, projectIDOrName, target)
	err = client.SetEnv(projectIDOrName, envVarNameTelegramSecret, telegramWebhookSecret, target)
	if err != nil {
		return err
	}

	return nil
}
