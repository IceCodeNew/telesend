package telebot

import (
	"os"
	"time"

	tele "gopkg.in/telebot.v4"
)

func botStart() {
	botToken, verbose := os.Getenv("TELESEND_BOT_TOKEN"), false
	if os.Getenv("VERBOSE") == "true" {
		verbose = true
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:   botToken,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: verbose,
	})
	if err != nil {
		panic(err)
	}

	bot.Handle("/start", func(c tele.Context) error {
		return c.Reply("Hello world!")
	})

	bot.Start()
}
