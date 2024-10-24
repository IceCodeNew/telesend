package telebot

import (
	"crypto/aes"
	"fmt"

	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/internal/app/db"
	"github.com/IceCodeNew/telesend/pkg/bark"
	"github.com/IceCodeNew/telesend/pkg/crypto"
	"github.com/IceCodeNew/telesend/pkg/uniqueID"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// NOTE: Make sure to use a pointer receiver to avoid copying the client data!
func (bsg *barkSenderGenerator) deviceKeyHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	var (
		// currStep = 1
		reply string
	)

	payloads := payloadsOfMessage(ctx)
	if payloads == nil || len(payloads) != 1 {
		reply = "ERROR: [User Input] Only takes one arg as the DeviceKey, please check your input"
		return replyNoDetailInternalErr(bot, ctx, nil, reply)
	}

	bsg.newBarkSender = &bark.BarkSender{
		Creator:   ctx.EffectiveSender.User.Id,
		DeviceKey: []byte(payloads[0]),
		Server:    "https://api.day.app/",
	}
	bsg.verifyMsg = &bark.BarkMessage{
		Group: "telesend",
		Sound: "telegraph.caf",
		Title: "Forwarded from Telegram",
	}

	reply = `
请发送 Bark 服务器地址，如：
Please send me the Bark Server address, e\.g\.
` +
		fmt.Sprintf("\n`/barkServer %s`\n", bsg.newBarkSender.Server)

	bsg.nextStep++

	_, err := ctx.EffectiveMessage.Reply(
		bot, reply, &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		},
	)
	return err
}

func (bsg *barkSenderGenerator) serverAddrInputHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	var (
		// currStep = 2
		reply string
	)
	// does not have to check whether the previous step successed or not
	// no op

	key, err := crypto.RandAsciiBytes(crypto.KeySizeAES256)
	if err != nil {
		reply := "ERROR: [Internal] Failed to generate AES key"
		return replyNoDetailInternalErr(bot, ctx, nil, reply)
	}
	iv, err := crypto.RandAsciiBytes(aes.BlockSize)
	if err != nil {
		reply := "ERROR: [Internal] Failed to generate AES IV"
		return replyNoDetailInternalErr(bot, ctx, nil, reply)
	}

	payloads := payloadsOfMessage(ctx)
	if payloads == nil || len(payloads) != 1 {
		reply = "ERROR: [User Input] Only takes one arg as the Bark server address, please check your input"
		return replyNoDetailInternalErr(bot, ctx, nil, reply)
	}
	// TODO: check whether the < serverAddr + deviceKey > is already existed

	bsg.newBarkSender.Server,
		bsg.newBarkSender.PreSharedSHA256Key,
		bsg.newBarkSender.PreSharedSHA256IV =
		payloads[0], key, iv

	quotedKey, quotedIv :=
		fmt.Sprintf("Key: `%s`\n", string(key)),
		fmt.Sprintf("Iv: `%s`\n", string(iv))

	reply = `
为避免不慎发送敏感信息，强制进行 Bark 加密推送。请在 iOS 客户端的推送加密设置中填入以下信息：
` +
		"当您完成客户端配置后。请发送 `/finish_gen_bark_sender` 来验证 Bark Server 是否正常工作。" +
		`
In case of accidentally sending sensitive messages, it is MANDATORY to send encrypted messages\.
Please configure the encryption settings of the iOS client according to the following info:
` +
		"Once you have completed the client configuration, please send `/finish_gen_bark_sender` to verify whether the Bark Server is working properly\\." +
		"\n\n" +
		quotedKey +
		quotedIv

	bsg.nextStep++

	_, err = ctx.EffectiveMessage.Reply(
		bot, reply, &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		},
	)
	return err
}

func (bsg *barkSenderGenerator) sendVerifyMsgHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	var (
		currStep = 3
		reply    string
	)
	// there is no turning back
	if bsg.nextStep != currStep {
		reply := previousStepsNotComplete(currStep, stepsGenBarkSender)
		return replyNoDetailInternalErr(bot, ctx, nil, reply)
	}

	senderID := uniqueID.UniqueID()
	bsg.newBarkSender.ID, bsg.verifyMsg.Title = senderID, "Verify Bark Sender"

	bsg.verifyMsg.Copy = "/verify_bark_sender " + senderID
	bsg.verifyMsg.Body = fmt.Sprintf(`
请使用通知自带的复制功能，并将复制到的字符串发给 telegram bot 进行验证。验证字符串应为以下内容：
Please utilize the copy message function of the received notification and send the copied message to this telegram bot for verification.
The verification message should be as follows:

%s`,
		bsg.verifyMsg.Copy)

	if err := bsg.newBarkSender.Send(bsg.verifyMsg, config.TSConfig.Verbose); err != nil {
		// allowing retries
		// newBarkSender, verifyMsg = nil, nil

		// no details in error message returned to users
		reply = "ERROR: [Internal] Failed to send verify message"
		return replyNoDetailInternalErr(bot, ctx, err, reply)
	}

	reply = `
验证消息已成功发送。请检查 iPhone 上的系统通知，并按要求回复验证字符串。
The initiate message has been sent successfully. Please check the system notifications on your iPhone and reply with the random string for verification.`

	bsg.nextStep++

	_, err := ctx.EffectiveMessage.Reply(bot, reply, nil)
	return err
}

func (bsg *barkSenderGenerator) verifyBarkSenderHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	var (
		currStep = 4
		reply    string
	)

	if bsg.nextStep != currStep {
		reply := previousStepsNotComplete(currStep, stepsGenBarkSender)
		return replyNoDetailInternalErr(bot, ctx, nil, reply)
	}

	payloads := payloadsOfMessage(ctx)
	if payloads == nil || len(payloads) != 1 {
		reply = "ERROR: [User Input] Only takes one arg as the verification message, please check your input"
		return replyNoDetailInternalErr(bot, ctx, nil, reply)
	}

	if payloads[0] != bsg.newBarkSender.ID {
		// allowing retries
		// newBarkSender, verifyMsg = nil, nil

		reply = "ERROR: [User Input] Bark sender verification failed"
		return replyNoDetailInternalErr(bot, ctx, nil, reply)
	}

	if err := bsg.newBarkSender.SelfEncrypt(); err != nil {
		// retry is not possible since the sender info might have been partially encrypted
		reply := "ERROR: [Internal] Failed to self-encrypt the new Bark Sender, please start over and give another try"
		return replyNoDetailInternalErr(bot, ctx, err, reply)
	}

	if err := db.StoreSender(bsg.newBarkSender); err != nil {
		// retry is not possible since the sender info was already encrypted
		reply := "ERROR: [Internal] Failed to store the new Bark Sender, please start over and give another try"
		return replyNoDetailInternalErr(bot, ctx, err, reply)
	}

	bsg.nextStep++

	_, err := ctx.EffectiveMessage.Reply(
		bot, "SUCCESS: Bark Server is working properly", nil,
	)
	return err
}
