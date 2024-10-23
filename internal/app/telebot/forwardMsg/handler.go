package forwardMsg

import (
	"github.com/IceCodeNew/telesend/pkg/bark"
)

var (
	nextStep = 1

	verifyMsg = &bark.BarkMessage{
		Group: "telesend",
		Sound: "telegraph.caf",
		Title: "Forwarded from Telegram",
	}
)
