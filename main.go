package main

import (
	"fmt"
	"os"

	"github.com/IceCodeNew/telesend/internal/app/config"
	bark "github.com/IceCodeNew/telesend/pkg/notificator"
	"github.com/IceCodeNew/telesend/pkg/uniqueID"
)

var testRcv = &bark.Sender{
	Creator:            666666,
	DeviceKey:          os.Getenv("test_DeviceKey"),
	ID:                 uniqueID.UniqueID(),
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
	fmt.Println(testRcv.ID)
}
