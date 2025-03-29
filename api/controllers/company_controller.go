package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"net/http"
)

func CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBind(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, company)
}
