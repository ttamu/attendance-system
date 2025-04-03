package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func HandleLineWebhook(channelSecret string, bot *messaging_api.MessagingApiAPI) gin.HandlerFunc {
	handler, err := webhook.NewWebhookHandler(channelSecret)
	if err != nil {
		log.Fatalf("failed to create webhook handler: %v", err)
	}

	handler.HandleEvents(func(req *webhook.CallbackRequest, r *http.Request) {
		for _, event := range req.Events {
			switch e := event.(type) {
			case webhook.MessageEvent:
				if msg, ok := e.Message.(webhook.TextMessageContent); ok {
					if msg.Text == "登録テスト" {
						_, err := bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								&messaging_api.TextMessage{Text: "接続OK!"},
							},
						})
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	})

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
