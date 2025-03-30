package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

// GetCompanyID コンテキストからJWTのcompany_idをuintとして取得
func GetCompanyID(c *gin.Context) (uint, error) {
	val, exists := c.Get("company_id")
	if !exists {
		return 0, errors.New("company_id not found in context")
	}

	if f, ok := val.(float64); ok {
		return uint(f), nil
	}

	return 0, errors.New("company_id is not a valid number")
}

// GetAccountID JWTからaccount_idを取得
func GetAccountID(c *gin.Context) (uint, error) {
	val, exists := c.Get("account_id")
	if !exists {
		return 0, errors.New("account_id not found in JWT")
	}

	if f, ok := val.(float64); ok {
		return uint(f), nil
	}

	return 0, errors.New("account_id is invalid")
}
