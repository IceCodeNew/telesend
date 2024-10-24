package main

import (
	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/internal/app/db"
	"github.com/IceCodeNew/telesend/internal/app/telebot"
	"github.com/IceCodeNew/telesend/pkg/bark"
	_ "go.uber.org/automaxprocs"
)

func init() {
	if err := config.ReadConfig(); err != nil {
		panic(err)
	}
}

func main() {
	// init db
	if err := db.CreateTable(&bark.BarkSender{}); err != nil {
		panic(err)
	}
	defer db.Close()

	// init telebot
	telebot.StartBot()
}
