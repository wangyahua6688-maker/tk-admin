package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "go-admin-full/internal/models"
    "go-admin-full/internal/services"
    "gorm.io/gorm"
)

type PermissionController struct {
    svc *services.PermissionService
}

func NewPermissionController(db *gorm.DB) *PermissionController {
    return &PermissionController{svc: services.NewPermissionService(db)}
}

func (pc *PermissionController) List(c *gin.Context) {
    items, err := pc.svc.List(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": items})
}

func (pc *PermissionController) Create(c *gin.Context) {
    var p models.Permission
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := pc.svc.Create(c.Request.Context(), &p); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": p})
}

func (pc *PermissionController) Get(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    p, err := pc.svc.Get(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": p})
}

func (pc *PermissionController) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := pc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
