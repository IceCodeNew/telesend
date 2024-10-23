package forwardMsg

import (
	bark "github.com/IceCodeNew/telesend/pkg/notificator"
)

var (
	nextStep = 1

	verifyMsg = &bark.Message{
		Group: "telesend",
		Sound: "telegraph.caf",
		Title: "Forwarded from Telegram",
	}
)
