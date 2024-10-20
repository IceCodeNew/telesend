package main

import (
	"os"

	"github.com/IceCodeNew/telesend/internal/app/config"
	bark "github.com/IceCodeNew/telesend/pkg/notificator"
)

var testRcv = &bark.Receiver{
	DeviceKey:          os.Getenv("test_DeviceKey"),
	PreSharedSHA256IV:  os.Getenv("test_SHA256IV"),
	PreSharedSHA256Key: os.Getenv("test_SHA256Key"),
	Server:             "https://api.day.app/",
}

var testMsg = &bark.Message{
	Badge: 1,
	Body:  "Test Bark Server",
	Copy:  "foo",
	Group: "test",
	Icon:  "https://day.app/assets/images/avatar.jpg",
	Sound: "minuet.caf",
	Title: "Test Title",
	URL:   "https://mritd.com",
}

func init() {
	if err := config.ReadConfig(); err != nil {
		panic(err)
	}
}

func main() {
	if err := testRcv.Send(testMsg, config.TSConfig.Verbose); err != nil {
		panic(err)
	}
}
