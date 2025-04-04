package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"github.com/t2469/attendance-system.git/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleLineWebhook(channelSecret string) gin.HandlerFunc {
	handler, err := webhook.NewWebhookHandler(channelSecret)
	if err != nil {
		log.Fatalf("Error: Failed to create webhook handler: %v", err)
	}

	handler.HandleEvents(func(req *webhook.CallbackRequest, r *http.Request) {
		for _, event := range req.Events {
			switch e := event.(type) {
			case webhook.MessageEvent:
				msg, ok := e.Message.(webhook.TextMessageContent)
				if !ok {
					continue
				}
				text := strings.TrimSpace(msg.Text)

				// 「登録テスト」メッセージの場合は接続確認用の返信
				if text == "登録テスト" {
					services.Reply(e.ReplyToken, "接続OK!")
					return
				}

				// 「登録 <社員ID> <名前>」形式のメッセージを処理
				regMsg := strings.Fields(text)
				if len(regMsg) != 3 || regMsg[0] != "登録" {
					services.Reply(e.ReplyToken, "「登録 <社員ID> <名前>」形式のメッセージを送信してください。")
					return
				}
				empIdStr, name := regMsg[1], regMsg[2]

				empId, convErr := strconv.Atoi(empIdStr)
				if convErr != nil {
					log.Println(convErr)
					services.Reply(e.ReplyToken, "社員IDの形式が正しくありません。")
					return
				}

				lineUserId, ok := services.GetUserId(e.Source)
				if !ok {
					log.Println(e.Source)
					services.Reply(e.ReplyToken, "UserIDの取得に失敗しました。")
					return
				}

				var emp models.Employee
				err := db.DB.Where("id = ? AND name = ?", empId, name).First(&emp).Error
				if err != nil {
					log.Println(err)
					services.Reply(e.ReplyToken, "登録できませんでした。IDまたは名前を確認してください。")
					return
				}

				isLinked := emp.LineUserID != nil
				if err := db.DB.Model(&emp).Updates(models.Employee{LineUserID: &lineUserId}).Error; err != nil {
					log.Println(err)
					services.Reply(e.ReplyToken, "登録中にエラーが発生しました。")
					return
				}

				if isLinked {
					services.Reply(e.ReplyToken, "登録情報を更新しました。")
				} else {
					services.Reply(e.ReplyToken, "登録が完了しました。")
				}
			}
		}
	})

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
