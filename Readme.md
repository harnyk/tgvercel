# tgvercel

`tgvercel` is a tool for easing the setup of Telegram bot in Vercel.

## What is it for?

So you want to create a telegram bot and host it on Vercel for free.
Ever feeld frustrated by the necessity of setting up a Telegram bot webhook?

`tgvercel` and its optional companion library (https://github.com/harnyk/tgvercelbot) automates this process for you making it as easy as possible.

For each fresh deployment of your Vercel project, `tgvercel` will:

-   get a URL of deployment
-   get a Telegram token and webhook secret from the environment variables stored in your Vercel project settings, corresponding to the target (production, or preview)
-   use Telegram API to set up the webhook to your Vercel project's function which will receive updates from Telegram at the specific deployment URL

Besides that, `tgvercel` also provides a CLI command to set the environment variables in the Vercel project.

On its turn, `tgvercelbot` allows to easily create a Vercel cloud function which would handle all the updates from Telegram. Feel free to browse its documentation [here](https://github.com/harnyk/tgvercelbot).

## Installation

If you have Go installed, you can call the following command to install `tgvercel`:

```bash
go install github.com/harnyk/tgvercel@latest
```

Or you can just run `tgvercel` directly:

```bash
go run github.com/harnyk/tgvercel@latest [arguments]
```

Another option is to install one of the pre-compiled binaries from the [Release](https://github.com/harnyk/tgvercel/releases) page.

If you have [`eget`](https://github.com/zyedidia/eget), the installation is as easy as:

```bash
eget harnyk/tgvercel
```

## Usage

### Step 1. Init the secrets

Suppose, you have created a Vercel project and a couple of Telegram bots (one for `preview` and one for `production`). You know the Telegram bot token and you are in the directory of the project (so `.vercel/project.json` is in the current directory). Also you are authenticated with Vercel CLI.

In order to initialize the `preview` environment, you can run:

```bash
tgvercel init --target=preview --telegram-token=YOUR_PREVIEW_TELEGRAM_TOKEN
```

Then, do the same for `production` environment:

```bash
tgvercel init --target=production --telegram-token=YOUR_PRODUCTION_TELEGRAM_TOKEN
```

This will create a pair of environment variables for each target environment. The variables created are:

-   `TELEGRAM_TOKEN` - Telegram bot token, the one that you provided.
-   `TELEGRAM_WEBHOOK_SECRET` - Telegram bot webhook secret, generated automatically. This is a `?secret=` part of the URL of the webhook, protecting the endpoint from unauthorized requests.

### Step 2. Create a Telegram bot endpoint

You can create an endpoint in any language, supported by Vercel. But if you want to write it in Go, you can use a conventional library [tgvercelbot](https://github.com/harnyk/tgvercelbot) which allows you to integrate telegram-bot-api with Vercel runtime easily.

Create a handler in `api/tg/webhook.go` with the following code:

```go
package handler

import (
	"net/http"

	"github.com/harnyk/tgvercelbot"
)

var tgv = tgvercelbot.New(tgvercelbot.DefaultOptions())

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	tgv.HandleWebhook(r, func(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
        if update.Message != nil {
            // do something, for example, echo:
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text))
        }
    })
}
```

### Step 3. Setup Telegram webhook

First of all, you have to deploy the project and get the URL of the deployment.

**Fun fact:** `vercel` command outputs the URL of the deployment to stdout, so you can save it to a variable.

```bash
# Deploy the project and get the URL of the deployment
DEPLOYMENT_URL=$(vercel)
```

Then, you can call `tgvercel hook` command with the following arguments:

```bash
# Tell Telegram to send updates to the URL of the deployment:
tgvercel hook $DEPLOYMENT_URL /api/tg/webhook
```

Here, `/api/tg/webhook` is the route that we created in the previous step. It may differ in your case, so you always have to specify it.

The following oneliner will do the same thing as the previous two commands:

```bash
tgvercel hook $(vercel) /api/tg/webhook
```

If you want to deploy to production, just add the `--prod` flag to the `vercel` command, as you would normally do:

```bash
tgvercel hook $(vercel --prod) /api/tg/webhook`
```

## Commands Reference

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
