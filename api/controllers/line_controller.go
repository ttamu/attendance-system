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

type lineCmd func(e webhook.MessageEvent, tokens []string)

var cmdMap = map[string]lineCmd{
	"登録":   register,
	"テスト": test,
}

func usage() string {
	type help struct {
		Cmd  string
		Desc string
	}

	helps := []help{
		{"登録 <社員ID> <名前>", "従業員とLINEアカウントを紐付けます。既に紐付いている場合は情報を更新します。"},
		{"出勤", "現在時刻で出勤打刻を登録します。"},
	}

	var sb strings.Builder
	sb.WriteString("==============================\n")
	sb.WriteString("ご利用ガイド\n")
	sb.WriteString("==============================\n\n")

	for _, c := range helps {
		sb.WriteString(c.Cmd + "\n")
		sb.WriteString("    └ " + c.Desc + "\n\n")
	}

	sb.WriteString("※ コマンドは半角スペースで区切って送信してください。")
	return sb.String()
}

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
				tokens := strings.Fields(text)

				if len(tokens) == 0 {
					services.Reply(e.ReplyToken, usage())
					return
				}

				cmd := tokens[0]
				if f, exists := cmdMap[cmd]; exists {
					f(e, tokens)
					return
				}

				services.Reply(e.ReplyToken, usage())
			}
		}
	})

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func test(e webhook.MessageEvent, _ []string) {
	services.Reply(e.ReplyToken, "接続OK!")
}

// register「登録 <社員ID> <名前>」形式のメッセージを処理
func register(e webhook.MessageEvent, tokens []string) {
	if len(tokens) != 3 || tokens[0] != "登録" {
		services.Reply(e.ReplyToken, "「登録 <社員ID> <名前>」形式のメッセージを送信してください。")
		return
	}
	empIdStr, name := tokens[1], tokens[2]

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
