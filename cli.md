# CLI documentsation

## Commands:

### Init

-   `tgvercel init <projectIdOrName>`: set the environment variables in the Vercel project

Arguments:

-   `<projectIdOrName>, optional`: Project ID or name, if not specified, the project ID will be taken from the file
-   `--target, env:TARGET, optional`: target environment (production, preview, development) (default: preview)
-   `--telegram-token, env:TELEGRAM_TOKEN, required`: Telegram bot token
-   `--token, env:VERCEL_TOKEN, optional`: Vercel token. If not specified, the token will be taken from the file
-   `--telegram-webhook-secret, env:TELEGRAM_WEBHOOK_SECRET, optional`: Telegram bot webhook secret. If not specified, the secret will be generated randomly

### Setup Telegram webhook

-   `tgvercel hook <deploymentIdOrUrl> <telegramBotRoute>`: setup Telegram webhook

Arguments:

-   `<deploymentIdOrUrl>`: Deployment ID or URL
-   `<telegramBotRoute>`: Telegram bot route
-   `--token, env:VERCEL_TOKEN, optional`: Vercel token. If not specified, the token will be taken from the file
