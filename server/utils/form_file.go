package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FormFileResult 保存校验后的文件元信息（不包含打开的句柄）
type FormFileResult struct {
	FileHeader  *multipart.FileHeader
	Filename    string // 生成的唯一文件名
	ContentType string
}

// ValidateFormFile 校验表单文件（类型、大小），返回 FileHeader 和元信息，不打开文件
// 如果字段为空（用户未上传），返回 nil, nil
func ValidateFormFile(ctx *gin.Context, fieldName string, maxSize int64, allowedTypes []string) (*FormFileResult, error) {
	fh, err := ctx.FormFile(fieldName)
	if err != nil {
		// 字段为空或不存在，属于可选场景
		return nil, nil
	}

	// 校验文件大小
	if fh.Size > maxSize {
		return nil, fmt.Errorf("文件大小不能超过 %dMB", maxSize/(1024*1024))
	}

	// 校验 Content-Type
	contentType := fh.Header.Get("Content-Type")
	allowed := slices.Contains(allowedTypes, contentType)
	if !allowed {
		return nil, errors.New("文件格式错误")
	}

	return &FormFileResult{
		FileHeader:  fh,
		Filename:    generateUniqueFilename(fh.Filename),
		ContentType: contentType,
	}, nil
}

// generateUniqueFilename 生成唯一文件名，格式: 时间戳_uuid.扩展名
func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Format("20060102150405")
	uniqueID := uuid.New().String()[:8]
	return timestamp + "_" + uniqueID + ext
}
