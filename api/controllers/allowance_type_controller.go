package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/db"
	"github.com/t2469/labor-management-system.git/helpers"
	"github.com/t2469/labor-management-system.git/models"
	"net/http"
)

func CreateAllowanceType(c *gin.Context) {
	var at models.AllowanceType
	if err := c.ShouldBindJSON(&at); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	at.CompanyID = companyID

	if err := db.DB.Create(&at).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, at)
}

func GetAllowanceTypes(c *gin.Context) {
	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var ats []models.AllowanceType
	if err := db.DB.Where("company_id = ?", companyID).Find(&ats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ats)
}
