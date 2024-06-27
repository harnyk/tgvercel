package main

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/harnyk/tgvercel/internal/commands"
	"github.com/harnyk/tgvercel/internal/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tgvercel",
	Short: "CLI tool for setting up Telegram bot in Vercel",
}

var initCmd = &cobra.Command{
	Use:   "init <projectIdOrName>",
	Short: "Set the environment variables in the Vercel project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		telegramToken := viper.GetString("telegram-token")
		target := viper.GetString("target")
		vercelToken := viper.GetString("token")
		telegramWebhookSecret := viper.GetString("telegram-webhook-secret")

		if telegramToken == "" {
			log.Println(color.RedString("Error: --telegram-token or TELEGRAM_TOKEN environment variable is required"))
			cmd.Help()
			os.Exit(1)
		}

		projectIdOrName := ""
		if len(args) > 0 {
			projectIdOrName = args[0]
		}

		// Implement the logic for setting environment variables in Vercel project here
		log.Println("Initializing Vercel project with the following settings:")
		log.Printf("Target: %s\n", target)
		log.Printf("Telegram Token: %s\n", telegramToken)
		if vercelToken != "" {
			log.Printf("Vercel Token: %s\n", vercelToken)
		}
		if telegramWebhookSecret != "" {
			log.Printf("Telegram Webhook Secret: %s\n", telegramWebhookSecret)
		} else {
			log.Println("Telegram Webhook Secret will be generated randomly")
		}

		command := commands.NewCmdInit(commands.CmdInitOptions{
			Target:          target,
			VercelToken:     vercelToken,
			TelegramToken:   telegramToken,
			TelegramSecret:  telegramWebhookSecret,
			ProjectIDOrName: projectIdOrName,
		})
		return command.Execute()
	},
}

var setupWebhookCmd = &cobra.Command{
	Use:   "hook <deploymentIdOrUrl> <telegramBotRoute>",
	Short: "Setup Telegram webhook",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deploymentIdOrUrl := args[0]
		telegramBotRoute := args[1]
		vercelToken := viper.GetString("token")

		// Implement the logic for setting up the Telegram webhook here
		log.Println("Setting up Telegram webhook with the following settings:")
		log.Printf("Deployment ID or URL: %s\n", deploymentIdOrUrl)
		log.Printf("Telegram Bot Route: %s\n", telegramBotRoute)
		if vercelToken != "" {
			log.Printf("Vercel Token: %s\n", helpers.Redact(vercelToken, 5))
		} else {
			log.Println("Vercel Token will be taken from the file")
		}

		command := commands.NewCmdHook(commands.CmdHookOptions{
			VercelToken:       vercelToken,
			DeploymentIDOrURL: deploymentIdOrUrl,
			TelegramBotRoute:  telegramBotRoute,
		})

		return command.Execute()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	initCmd.Flags().String("target", "preview", "target environment (production, preview, development)")
	initCmd.Flags().String("telegram-token", "", "Telegram bot token")
	initCmd.Flags().String("token", "", "Vercel token. If not specified, the token will be taken from the file")
	initCmd.Flags().String("telegram-webhook-secret", "", "Telegram bot webhook secret. If not specified, the secret will be generated randomly")

	setupWebhookCmd.Flags().String("token", "", "Vercel token. If not specified, the token will be taken from the file")

	viper.BindPFlag("target", initCmd.Flags().Lookup("target"))
	viper.BindPFlag("telegram-token", initCmd.Flags().Lookup("telegram-token"))
	viper.BindPFlag("token", initCmd.Flags().Lookup("token"))
	viper.BindPFlag("telegram-webhook-secret", initCmd.Flags().Lookup("telegram-webhook-secret"))

	viper.BindPFlag("token", setupWebhookCmd.Flags().Lookup("token"))

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(setupWebhookCmd)
}

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	log.Println(banner)

	if err := rootCmd.Execute(); err != nil {
		log.Println(color.RedString(err.Error()))
		os.Exit(1)
	}
}
