package telebot

import (
	"github.com/IceCodeNew/telesend/pkg/bark"
)

type barkSenderGenerator struct {
	nextStep      int
	newBarkSender *bark.BarkSender
	verifyMsg     *bark.BarkMessage
}
type msgForwarder struct {
	nextStep   int
	barkSender *bark.BarkSender
	message    *bark.BarkMessage
}

func newBarkSenderGenerator() *barkSenderGenerator {
	return &barkSenderGenerator{
		nextStep: 1,
	}
}
func newMsgForwarder() *msgForwarder {
	return &msgForwarder{
		nextStep: 1,
	}
}

var (
	stepsGenBarkSender = []string{
		"set the DeviceKey",
		"set the Bark Server address",
		"request a string for verifying the Bark Sender",
		"finish the setup",
	}
)
