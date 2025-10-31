package security

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

func NewCORSMiddleware(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 设置CORS头
		// Access-Control-Allow-Origin 只能是单个源或 "*"，不能是逗号分隔的列表
		if len(config.AllowOrigins) > 0 {
			if config.AllowOrigins[0] == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				// 检查请求的 Origin 是否在允许列表中
				if slices.Contains(config.AllowOrigins, origin) {
					c.Header("Access-Control-Allow-Origin", origin)
				}
			}
		}

		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 缓存预检请求结果24小时

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
