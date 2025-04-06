package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"github.com/t2469/attendance-system.git/services"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type lineCmd func(e webhook.MessageEvent, tokens []string)

var cmdMap = map[string]lineCmd{
	"登録":     register,
	"テスト":   test,
	"出勤":     clockIn,
	"退勤":     clockOut,
	"休憩開始": breakBegin,
	"休憩終了": breakEnd,
}

func usage() string {
	type help struct {
		Cmd  string
		Desc string
	}

	helps := []help{
		{"登録 <社員ID> <名前>", "従業員とLINEアカウントを紐付けます。既に紐付いている場合は情報を更新します。"},
		{"出勤", "現在時刻で出勤打刻を登録します。"},
		{"退勤", "現在時刻で退勤打刻を登録します。"},
		{"休憩開始", "現在時刻で休憩開始の打刻を登録します。"},
		{"休憩終了", "現在時刻で休憩終了の打刻を登録します。"},
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

	lineUserId, ok := services.GetLineUserId(e.Source)
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

func clockIn(e webhook.MessageEvent, tokens []string) {
	clock(e, tokens, models.ClockIn, "出勤")
}

func clockOut(e webhook.MessageEvent, tokens []string) {
	clock(e, tokens, models.ClockOut, "退勤")
}

func breakBegin(e webhook.MessageEvent, tokens []string) {
	clock(e, tokens, models.BreakBegin, "休憩開始")
}
func breakEnd(e webhook.MessageEvent, tokens []string) {
	clock(e, tokens, models.BreakEnd, "休憩終了")
}

func clock(e webhook.MessageEvent, _ []string, clockType models.TimeClockType, clockName string) {
	lineUserId, ok := services.GetLineUserId(e.Source)
	if !ok {
		services.Reply(e.ReplyToken, "ユーザー情報の取得に失敗しました。")
		return
	}

	// 1つのLINEアカウントで複数ユーザーを登録している可能性があるため Find で複数レコードを取得
	var employees []models.Employee
	err := db.DB.Where("line_user_id = ?", lineUserId).Find(&employees).Error
	if err != nil || len(employees) == 0 {
		services.Reply(e.ReplyToken, "従業員が登録されていません。")
		return
	}

	now := time.Now()
	var okList []string
	var ngList []string

	for _, emp := range employees {
		_, err := services.RecordTimeClock(emp.ID, clockType, now)
		info := fmt.Sprintf("%sさん", emp.Name)
		if err != nil {
			log.Println(err)
			ngList = append(ngList, info)
		} else {
			okList = append(okList, info)
		}
	}

	var sb strings.Builder
	// 登録が1人だけの場合はシンプルなメッセージにする
	if len(employees) == 1 {
		if len(okList) == 1 {
			sb.WriteString(fmt.Sprintf("%sの%s打刻を行いました！", okList[0], clockName))
		} else {
			sb.WriteString(fmt.Sprintf("%sの%s打刻に失敗しました...", ngList[0], clockName))
		}
	} else {
		// 複数人いる場合のみリスト表示する
		if len(okList) > 0 {
			sb.WriteString(fmt.Sprintf("以下の従業員の%s打刻を行いました！\n", clockName))
			for _, info := range okList {
				sb.WriteString("・" + info + "\n")
			}
		}
		if len(ngList) > 0 {
			sb.WriteString(fmt.Sprintf("以下の従業員の%s打刻に失敗しました。\n", clockName))
			for _, info := range ngList {
				sb.WriteString("・" + info + "\n")
			}
		}
	}

	services.Reply(e.ReplyToken, sb.String())
}
