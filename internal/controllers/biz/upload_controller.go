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

// UploadController 定义UploadController相关结构。
type UploadController struct {
}

// NewUploadController 创建UploadController实例。
func NewUploadController() *UploadController {
	// 返回当前处理结果。
	return &UploadController{}
}

// UploadImage 处理图片上传
func (uc *UploadController) UploadImage(c *gin.Context) {
	// 定义并初始化当前变量。
	file, err := c.FormFile("file")
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "no file uploaded")
		// 返回当前处理结果。
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	// 判断条件并进入对应分支逻辑。
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "invalid file type, only images are allowed")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	cfg := config.GetConfig()
	// 定义并初始化当前变量。
	savePath := cfg.Upload.SavePath
	// 定义并初始化当前变量。
	baseURL := cfg.Upload.BaseURL

	// 确保保存目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusInternalServerError, "failed to create save directory")
		// 返回当前处理结果。
		return
	}

	// 生成唯一文件名
	newFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	// 定义并初始化当前变量。
	dst := filepath.Join(savePath, newFilename)

	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusInternalServerError, "failed to save file")
		// 返回当前处理结果。
		return
	}

	// 返回相对路径 (baseURL + filename)
	// 如果 baseURL 是 /uploads, 结果是 /uploads/xxx.jpg
	relativePath := filepath.Join(baseURL, newFilename)
	// 在 Windows 上 filepath.Join 会用 \, 统一转换为 /
	relativePath = filepath.ToSlash(relativePath)

	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, gin.H{
		// 处理当前语句逻辑。
		"url": relativePath,
	})
}
