package biz

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin/config"
	"go-admin/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 声明当前变量。
var (
	// maxUploadImageSize 单次图片上传体积上限（5MB）。
	maxUploadImageSize int64 = 5 << 20
	// allowedImageContentTypes 声明允许的图片 MIME 类型集合。
	allowedImageContentTypes = map[string]struct{}{
		// 处理当前语句逻辑。
		"image/jpeg": {},
		// 处理当前语句逻辑。
		"image/png": {},
		// 处理当前语句逻辑。
		"image/gif": {},
		// 处理当前语句逻辑。
		"image/webp": {},
	}
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
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "no file uploaded")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if file.Size <= 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "empty file is not allowed")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if file.Size > maxUploadImageSize {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "file size exceeds 5MB limit")
		// 返回当前处理结果。
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	// 判断条件并进入对应分支逻辑。
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid file type, only images are allowed")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	contentType, err := detectImageContentType(file)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "failed to read upload file")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if _, ok := allowedImageContentTypes[contentType]; !ok {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid file content type")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	cfg := config.GetConfig()
	// 定义并初始化当前变量。
	savePath := strings.TrimSpace(cfg.Upload.SavePath)
	// 定义并初始化当前变量。
	baseURL := strings.TrimSpace(cfg.Upload.BaseURL)
	// 判断条件并进入对应分支逻辑。
	if savePath == "" || baseURL == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, "upload config is invalid")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if !strings.HasPrefix(baseURL, "/") {
		// 更新当前变量或字段值。
		baseURL = "/" + strings.TrimPrefix(baseURL, "/")
	}
	// 判断条件并进入对应分支逻辑。
	if strings.Contains(baseURL, "..") {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, "upload base_url is invalid")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	absSavePath, err := filepath.Abs(savePath)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, "failed to resolve save directory")
		// 返回当前处理结果。
		return
	}

	// 确保保存目录存在
	if err := os.MkdirAll(absSavePath, 0755); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, "failed to create save directory")
		// 返回当前处理结果。
		return
	}

	// 生成唯一文件名
	newFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	// 定义并初始化当前变量。
	dst := filepath.Join(absSavePath, newFilename)
	// 定义并初始化当前变量。
	absDst, err := filepath.Abs(dst)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, "failed to resolve upload file path")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if !strings.HasPrefix(absDst, absSavePath+string(os.PathSeparator)) {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, "invalid upload file path")
		// 返回当前处理结果。
		return
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, absDst); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, "failed to save file")
		// 返回当前处理结果。
		return
	}

	// 返回相对路径 (baseURL + filename)
	// 如果 baseURL 是 /uploads, 结果是 /uploads/xxx.jpg
	relativePath := filepath.Join(baseURL, newFilename)
	// 在 Windows 上 filepath.Join 会用 \, 统一转换为 /
	relativePath = filepath.ToSlash(relativePath)

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{
		// 处理当前语句逻辑。
		"url": relativePath,
	})
}

// detectImageContentType 检测上传文件 MIME 类型。
func detectImageContentType(fileHeader *multipart.FileHeader) (string, error) {
	// 定义并初始化当前变量。
	src, err := fileHeader.Open()
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", err
	}
	// 注册延迟执行逻辑。
	defer src.Close()
	// 定义并初始化当前变量。
	buf := make([]byte, 512)
	// 定义并初始化当前变量。
	n, err := src.Read(buf)
	// 判断条件并进入对应分支逻辑。
	if err != nil && err != io.EOF {
		// 返回当前处理结果。
		return "", err
	}
	// 判断条件并进入对应分支逻辑。
	if n <= 0 {
		// 返回当前处理结果。
		return "", fmt.Errorf("empty file")
	}
	// 定义并初始化当前变量。
	contentType := strings.ToLower(strings.TrimSpace(http.DetectContentType(buf[:n])))
	// 定义并初始化当前变量。
	parts := strings.Split(contentType, ";")
	// 返回当前处理结果。
	return strings.TrimSpace(parts[0]), nil
}
