package services

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/t2469/attendance-system.git/config"
	"log"
)

func InitLineBot(cfg *config.Config) *messaging_api.MessagingApiAPI {
	bot, err := messaging_api.NewMessagingApiAPI(cfg.LineChannelToken)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}
