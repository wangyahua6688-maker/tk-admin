package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type RoleController struct{ db *gorm.DB }

func NewRoleController(db *gorm.DB) *RoleController { return &RoleController{db: db} }

func (rc *RoleController) List(c *gin.Context) {
	var roles []models.Role
	rc.db.Preload("Permissions").Find(&roles)
	c.JSON(http.StatusOK, gin.H{"data": roles})
}
func (rc *RoleController) Create(c *gin.Context) {
	var r models.Role
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rc.db.Create(&r)
	c.JSON(http.StatusOK, gin.H{"data": r})
}
