package main

import (
	"github.com/IceCodeNew/telesend/internal/app/aead"
	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/internal/app/db"
	"github.com/IceCodeNew/telesend/internal/app/telebot"
	"github.com/IceCodeNew/telesend/pkg/bark"
)

func init() {
	if err := config.ReadConfig(); err != nil {
		panic(err)
	}
}

func main() {
	// init aead cipher
	if err := aead.InitCipher(); err != nil {
		panic(err)
	}

	// init db
	if err := db.CreateTable(&bark.BarkSender{}); err != nil {
		panic(err)
	}
	defer db.Close()

	// init telebot
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
