package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	CompanyID uint   `json:"company_id"`
	IsAdmin   bool   `json:"is_admin"`
}

func CurrentAccount(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報がありません"})
	}

	var account models.Account
	if err := db.DB.Preload("Company").First(&account, accountID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "アカウント情報の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       account.ID,
		"email":    account.Email,
		"is_admin": account.IsAdmin,
		"company":  account.Company,
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	account := models.Account{
		Email:     req.Email,
		Password:  string(password),
		CompanyID: req.CompanyID,
		IsAdmin:   req.IsAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}
