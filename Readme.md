# tgvercel

An easy way to deploy Telegram bots to Vercel

## Install

```bash
go get -u github.com/harnyk/tgvercel
```

## Usage

1. Prepare your environment variables on Vercel:

-   `TGVERCEL_KEY` - internal service key, used to set up the webhook.
    Put any random string here
-   `TGVERCEL_TOKEN` - Telegram bot token

2. Create the following files in your project:

-   `api/tg/setup.go`:

```go
package handler

import (
	"net/http"

	"github.com/harnyk/tgvercel"
)

var tgv = tgvercel.New(tgvercel.DefaultOptions())

func SetupHandler(w http.ResponseWriter, r *http.Request) {
	tgv.HandleSetup(w, r)
}
```

-   `api/tg/webhook.go`:

```go
package handler

import (
	"net/http"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OnUpdate(
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update) {
	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "echo: "+update.Message.Text)
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	tgv.HandleWebhook(r, botlogic.OnUpdate)
}
```

3. Deploy your project to Vercel

4. Once deployed, you have to set up the webhook:

```bash
curl 'https://{YOUR_DEPLOYMENT_SUBDOMAIN}.vercel.app/api/tg/setup?key={TGVERCEL_KEY}'
```

, where:

-   YOUR_DEPLOYMENT_SUBDOMAIN - your Vercel deployment subdomain
-   TGVERCEL_KEY - internal service key, that you set up in the first step

5. You can now send messages to your Telegram bot!

## Full example

For the full example, go here: https://github.com/harnyk/tgvercel-example

The example covers all above plus the ability to organize your code to run bot locally.
