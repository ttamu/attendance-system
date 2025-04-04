package services

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
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

func GetUserId(source webhook.SourceInterface) (string, bool) {
	switch s := source.(type) {
	case webhook.UserSource:
		return s.UserId, true
	case *webhook.UserSource:
		return s.UserId, true
	}
	return "", false
}

func Reply(bot *messaging_api.MessagingApiAPI, replyToken string, msg string) {
	_, err := bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: msg,
			},
		},
	})
	if err != nil {
		log.Printf(err.Error())
	}
}
