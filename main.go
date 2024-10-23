package main

import (
	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/internal/app/telebot"
)

func init() {
	if err := config.ReadConfig(); err != nil {
		panic(err)
	}
}

func main() {
	telebot, err := telebot.NewBot()
	if err != nil {
		panic(err)
	}
	defer func() {
		_, err = telebot.Close()
		if err != nil {
			panic(err)
		}
	}()
	defer telebot.Stop()

	telebot.Start()
}
