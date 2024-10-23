package telebot

import (
	"time"

	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/internal/app/telebot/genBarkSender"
	tele "gopkg.in/telebot.v4"
)

func NewBot() (*tele.Bot, error) {
	bot, err := tele.NewBot(tele.Settings{
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Token:   config.TSConfig.BotToken,
		Verbose: config.TSConfig.Verbose,
	})
	if err != nil {
		return nil, err
	}

	// Memo: command using '-', U+002D, HYPHEN-MINUS will cause problem!
	bot.Handle("/start", startHandler)
	bot.Handle("/deviceKey", genBarkSender.DeviceKeyHandler)
	bot.Handle("/barkServer", genBarkSender.ServerAddrInputHandler)
	bot.Handle("/finish_gen_bark_sender", genBarkSender.SendVerifyMsgHandler)
	bot.Handle("/verify_bark_sender", genBarkSender.VerifyBarkSenderHandler)
	return bot, nil
}

func startHandler(c tele.Context) error {
	return c.Reply(`
您好，我可以将所有收到的信息通过 Bark 转发给 iPhone。如果您已经准备好，请发送：
Hi there, I can forward all the messages I received to your iPhone through the Bark service\. When you are ready, please send me the DeviceKey:

`+
		"`/deviceKey < Your DeviceKey >`",
		tele.ModeMarkdownV2,
	)
}
