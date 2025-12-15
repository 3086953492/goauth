package utils

import (
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/3086953492/gokit/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"goauth/dto"
)

func GetFormFile(ctx *gin.Context, fieldName string, maxSize int64, allowedTypes []string) (dto.FileMeta, error) {
	if file, err := ctx.FormFile(fieldName); err == nil {
		fileData, err := file.Open()
		if err != nil {
			return dto.FileMeta{}, errors.InvalidInput().Msg("文件读取失败").Err(err).Build()
		}
		defer fileData.Close()

		if file.Size > maxSize {
			return dto.FileMeta{}, errors.InvalidInput().Msg("文件大小不能超过 " + strconv.FormatInt(maxSize, 10) + "MB").Build()
		}

		if !slices.Contains(allowedTypes, file.Header.Get("Content-Type")) {
			return dto.FileMeta{}, errors.InvalidInput().Msg("文件格式错误").Build()
		}

		return dto.FileMeta{
			Data:        fileData,
			Filename:    generateUniqueFilename(file.Filename),
			ContentType: file.Header.Get("Content-Type"),
			Size:        file.Size,
		}, nil
	}
	return dto.FileMeta{}, nil
}

// generateUniqueFilename 生成唯一文件名，格式: 时间戳_uuid.扩展名
func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Format("20060102150405")
	uniqueID := uuid.New().String()[:8]
	return timestamp + "_" + uniqueID + ext
}
