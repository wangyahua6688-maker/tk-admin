package biz

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-admin-full/config"
	"go-admin-full/internal/utils"
)

type UploadController struct {
}

func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadImage 处理图片上传
func (uc *UploadController) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "no file uploaded")
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		utils.JSONError(c, http.StatusBadRequest, "invalid file type, only images are allowed")
		return
	}

	cfg := config.GetConfig()
	savePath := cfg.Upload.SavePath
	baseURL := cfg.Upload.BaseURL

	// 确保保存目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to create save directory")
		return
	}

	// 生成唯一文件名
	newFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	dst := filepath.Join(savePath, newFilename)

	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	// 返回相对路径 (baseURL + filename)
	// 如果 baseURL 是 /uploads, 结果是 /uploads/xxx.jpg
	relativePath := filepath.Join(baseURL, newFilename)
	// 在 Windows 上 filepath.Join 会用 \, 统一转换为 /
	relativePath = filepath.ToSlash(relativePath)

	utils.JSONOK(c, gin.H{
		"url": relativePath,
	})
}
