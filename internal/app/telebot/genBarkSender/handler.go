package genBarkSender

import (
	"crypto/aes"
	"fmt"

	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/internal/app/db"
	"github.com/IceCodeNew/telesend/pkg/bark"
	"github.com/IceCodeNew/telesend/pkg/crypto"
	"github.com/IceCodeNew/telesend/pkg/uniqueID"
	tele "gopkg.in/telebot.v4"
)

var (
	newBarkSender = &bark.BarkSender{
		Server: "https://api.day.app/",
	}

	verifyMsg = &bark.BarkMessage{
		Group: "telesend",
		Sound: "telegraph.caf",
		Title: "Forwarded from Telegram",
	}
)

func DeviceKeyHandler(c tele.Context) error {
	// currStep := 1

	// does not have to check whether the previous step successed or not
	// no op

	newBarkSender.Creator = c.Sender().ID
	deviceKey := c.Message().Payload
	if len(deviceKey) == 0 {
		newBarkSender, verifyMsg = nil, nil
		return fmt.Errorf("ERROR: [User Input] DeviceKey is empty")
	}

	newBarkSender.DeviceKey = []byte(deviceKey)

	nextStep++

	return c.Reply(`
请发送 Bark 服务器地址，如：`+
		`
Please send me the Bark Server address, e\.g\.
`+fmt.Sprintf("\n`/barkServer %s`", newBarkSender.Server),
		tele.ModeMarkdownV2,
	)
}

func ServerAddrInputHandler(c tele.Context) error {
	// currStep := 2

	// does not have to check whether the previous step successed or not
	// no op

	serverAddr := c.Message().Payload
	if len(serverAddr) == 0 {
		newBarkSender, verifyMsg = nil, nil
		return fmt.Errorf("ERROR: [User Input] Bark Server URL is empty")
	}

	// TODO: check whether the < serverAddr + deviceKey > is already existed

	newBarkSender.Server = serverAddr

	key, err := crypto.RandAsciiBytes(crypto.KeySizeAES256)
	if err != nil {
		newBarkSender, verifyMsg = nil, nil
		return fmt.Errorf("ERROR: [Internal] Failed to generate AES key")
	}
	iv, err := crypto.RandAsciiBytes(aes.BlockSize)
	if err != nil {
		newBarkSender, verifyMsg = nil, nil
		return fmt.Errorf("ERROR: [Internal] Failed to generate AES IV")
	}
	newBarkSender.PreSharedSHA256Key,
		newBarkSender.PreSharedSHA256IV = key, iv

	nextStep++
	quotedKey, quotedIv :=
		fmt.Sprintf("Key: `%s`\n", string(newBarkSender.PreSharedSHA256Key)),
		fmt.Sprintf("Iv: `%s`\n", string(newBarkSender.PreSharedSHA256IV))

	return c.Reply(`
为避免不慎发送敏感信息，强制进行 Bark 加密推送。请在 iOS 客户端的推送加密设置中填入以下信息：
`+
		"当您完成客户端配置后。请发送 `/finish_gen_bark_sender` 来验证 Bark Server 是否正常工作。"+
		`
In case of accidentally sending sensitive messages, it is MANDATORY to send encrypted messages\.
Please configure the encryption settings of the iOS client according to the following info:
`+
		"Once you have completed the client configuration, please send `/finish_gen_bark_sender` to verify whether the Bark Server is working properly\\."+
		"\n\n"+
		quotedKey+
		quotedIv,
		tele.ModeMarkdownV2,
	)
}

func SendVerifyMsgHandler(c tele.Context) error {
	currStep := 3
	// there is no turning back
	if nextStep != currStep ||
		newBarkSender == nil ||
		verifyMsg == nil {
		reply := previousStepsNotComplete(currStep)
		_ = c.Reply(reply)
		return fmt.Errorf(reply)
	}

	senderID := uniqueID.UniqueID()
	newBarkSender.ID, verifyMsg.Title = senderID, "Verify Bark Sender"

	verifyMsg.Copy = "/verify_bark_sender " + senderID
	verifyMsg.Body = fmt.Sprintf(`
请使用通知自带的复制功能，并将复制到的字符串发给 telegram bot 进行验证。验证字符串应为以下内容：
Please utilize the copy message function of the received notification and send the copied message to this telegram bot for verification.
The verification message should be as follows:

%s`,
		verifyMsg.Copy)
	if err := newBarkSender.Send(verifyMsg, config.TSConfig.Verbose); err != nil {
		// allowing retries
		// newBarkSender, verifyMsg = nil, nil

		// no details in error message returned to users
		reply := "ERROR: [Internal] Failed to send verify message"
		_ = c.Reply(reply)

		return fmt.Errorf("%s: %v", reply, err)
	}

	nextStep++

	return c.Reply(`
验证消息已成功发送。请检查 iPhone 上的系统通知，并按要求回复验证字符串。
The initiate message has been sent successfully. Please check the system notifications on your iPhone and reply with the random string for verification.`,
	)
}

// Step 4
func VerifyBarkSenderHandler(c tele.Context) error {
	currStep := 4
	if nextStep != currStep ||
		newBarkSender == nil ||
		verifyMsg == nil {
		reply := previousStepsNotComplete(currStep)
		_ = c.Reply(reply)
		return fmt.Errorf(reply)
	}

	if c.Message().Payload != newBarkSender.ID {
		// allowing retries
		// newBarkSender, verifyMsg = nil, nil

		return c.Reply("ERROR: [User Input] Verification failed")
	}

	newBarkSender.SelfEncrypt()
	if err := db.StoreSender(newBarkSender); err != nil {
		// retry is not possible since the sender info was already encrypted
		newBarkSender, verifyMsg = nil, nil
		reply := "ERROR: [Internal] Failed to save new Bark Sender, please try again"
		_ = c.Reply(reply)

		return fmt.Errorf(`%s
DEBUG: original error was:
%v`,
			reply, err)
	}

	nextStep++

	return c.Reply("SUCCESS: Bark Server is working properly")
}
