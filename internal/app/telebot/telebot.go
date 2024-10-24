package telebot

import (
	"fmt"
	"time"

	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func StartBot() {
	// Create bot from environment value.
	bot, err := gotgbot.NewBot(config.TSConfig.BotToken, nil)
	// bot, err := gotgbot.NewBot("7224021452:AAHiJMav78A2SA896JtQyxcqardpp2gy6No", nil)

	if err != nil {
		panic(fmt.Errorf(
			"ERROR: [Internal] Failed to create bot: %v", err,
		))
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// We create our stateful bot client here.
	bsg := newBarkSenderGenerator()

	dispatcher.AddHandler(handlers.NewCommand("start", startHandler))
	dispatcher.AddHandler(handlers.NewCommand("deviceKey", bsg.deviceKeyHandler))
	dispatcher.AddHandler(handlers.NewCommand("barkServer", bsg.serverAddrInputHandler))
	dispatcher.AddHandler(handlers.NewCommand("finish_gen_bark_sender", bsg.sendVerifyMsgHandler))
	dispatcher.AddHandler(handlers.NewCommand("verify_bark_sender", bsg.verifyBarkSenderHandler))

	// Start receiving updates.
	if err := updater.StartPolling(
		bot, &ext.PollingOpts{
			DropPendingUpdates: true,
			GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
				Timeout: 9,
				RequestOpts: &gotgbot.RequestOpts{
					Timeout: time.Second * 10,
				},
			},
		}); err != nil {
		panic(fmt.Errorf(
			"ERROR: [Internal] Failed to start polling: %v", err,
		))
	}

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func startHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	reply := `
您好，我可以将所有收到的信息通过 Bark 转发给 iPhone。如果您已经准备好，请发送：
Hi there, I can forward all the messages I received to your iPhone through the Bark service\. When you are ready, please send me the DeviceKey:

` +
		"`/deviceKey < Your DeviceKey >`"

	if _, err := ctx.EffectiveMessage.Reply(
		bot, reply,
		&gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		},
	); err != nil {
		return fmt.Errorf("ERROR: [Internal] Failed to send message: %v", err)
	}
	return nil
}
