package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/helpers"
	"github.com/t2469/attendance-system.git/models"
	"github.com/t2469/attendance-system.git/services"
	"net/http"
	"strconv"
	"time"
)

type ClockReqInput struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Time       string `json:"time" binding:"required"`
	Reason     string `json:"reason"`
}

type ClockReqOutput struct {
	models.ClockRequest
	EmployeeName string `json:"employee_name"`
}

func CreateClockRequest(c *gin.Context) {
	// URLパラメータから対象の打刻IDを取得
	clkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clock id"})
		return
	}

	var input ClockReqInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 打刻種別を TimeClockType に変換
	tcType := models.TimeClockType(input.Type)
	switch tcType {
	case models.ClockIn, models.ClockOut, models.BreakBegin, models.BreakEnd:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time clock type: " + input.Type})
		return
	}

	// トークンから会社IDを取得
	compID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// リクエストされた社員IDが、その会社に属しているかチェック
	if err := helpers.CheckEmployeeAccess(input.EmployeeID, compID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// 対象のTimeClockが存在し、かつその社員のものであるかチェック
	var clk models.TimeClock
	if err := db.DB.First(&clk, clkID).Error; err != nil || clk.EmployeeID != input.EmployeeID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to request change for this clock"})
		return
	}

	// リクエストされた時刻をパースする
	reqTime, err := helpers.ParseTimestamp(input.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
		return
	}

	// ClockRequestを作成（StatusはデフォルトでPending）
	req := models.ClockRequest{
		EmployeeID: input.EmployeeID,
		ClockID:    uint(clkID),
		Type:       tcType,
		Time:       reqTime,
		Status:     models.Pending,
		Reason:     input.Reason,
	}

	if err := db.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func GetClockRequests(c *gin.Context) {
	// JWTから会社IDを取得
	compID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	employeeID := c.Query("employee_id")
	status := c.Query("status")

	var req []ClockReqOutput

	query := db.DB.
		Table("clock_requests").
		Select("clock_requests.*, employees.name as employee_name").
		Joins("JOIN employees ON employees.id = clock_requests.employee_id").
		Where("employees.company_id = ?", compID)

	// 従業員IDで絞り込み
	if employeeID != "" {
		query = query.Where("clock_requests.employee_id = ?", employeeID)
	}

	// ステータスで絞り込み
	if status != "" {
		query = query.Where("clock_requests.status = ?", status)
	}

	if err := query.Order("clock_requests.created_at DESC").Scan(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// ApproveClockRequest は打刻修正申請を承認し、TimeClockとWorkRecordを更新する（管理者専用）
func ApproveClockRequest(c *gin.Context) {
	// パスパラメータからリクエストIDを取得
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	// JWTトークンから認証情報を取得
	accountID, err := helpers.GetAccountID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	isAdmin, err := helpers.GetIsAdmin(c)
	if err != nil || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 対象の修正申請を取得
	var req models.ClockRequest
	if err := db.DB.First(&req, reqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	// 既にレビュー済みならエラー
	if req.Status != models.Pending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request already reviewed"})
		return
	}

	// 対象の従業員が同一会社に属しているか確認
	var emp models.Employee
	if err := db.DB.First(&emp, req.EmployeeID).Error; err != nil || emp.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to review request for this employee"})
		return
	}

	// 修正対象のTimeClockを取得
	var clock models.TimeClock
	if err := db.DB.First(&clock, req.ClockID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "original clock not found"})
		return
	}

	// 打刻の内容を申請内容で更新
	clock.Type = req.Type
	clock.Timestamp = req.Time
	if err := db.DB.Save(&clock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update clock"})
		return
	}

	// 該当日の WorkRecord を再集計
	day := clock.Timestamp.In(time.Local).Truncate(24 * time.Hour)
	if err := services.UpsertWorkRecord(req.EmployeeID, day); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update work record"})
		return
	}

	// 申請ステータスを更新
	now := time.Now()
	req.Status = models.Approved
	req.ReviewedBy = &accountID
	req.ReviewedAt = &now
	if err := db.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "request approved",
		"request": req,
	})
}

// RejectClockRequest は打刻修正申請を却下する（管理者のみ）
func RejectClockRequest(c *gin.Context) {
	// リクエストIDをURLパラメータから取得
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	// JWTトークンから認証情報を取得
	accountID, err := helpers.GetAccountID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	isAdmin, err := helpers.GetIsAdmin(c)
	if err != nil || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}
	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 対象のClockRequestをDBから取得
	var req models.ClockRequest
	if err := db.DB.First(&req, reqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	// すでに処理済みの場合はエラー
	if req.Status != models.Pending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request already reviewed"})
		return
	}

	// 対象の従業員が同一会社に属しているか確認
	var emp models.Employee
	if err := db.DB.First(&emp, req.EmployeeID).Error; err != nil || emp.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to review request for this employee"})
		return
	}

	// 却下処理（打刻の更新やWorkRecordの再計算は行わず、申請状態のみ更新）
	now := time.Now()
	req.Status = models.Rejected
	req.ReviewedBy = &accountID
	req.ReviewedAt = &now

	if err := db.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "request rejected",
		"request": req,
	})
}
