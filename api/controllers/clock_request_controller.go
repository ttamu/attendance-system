package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/helpers"
	"github.com/t2469/attendance-system.git/models"
	"net/http"
	"strconv"
)

type ClockRequestInput struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Time       string `json:"time" binding:"required"`
	Reason     string `json:"reason"`
}

func CreateClockRequest(c *gin.Context) {
	// URLパラメータから対象の打刻IDを取得
	clkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clock id"})
		return
	}

	var input ClockRequestInput
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

	var requests []models.ClockRequest

	query := db.DB.
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

	if err := query.Order("created_at DESC").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}
