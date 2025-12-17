package security

import (
	"slices"
	"strings"

	"github.com/3086953492/gokit/config/types"
	"github.com/gin-gonic/gin"
)

func NewCORSMiddleware(corsConfig types.CorsMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 设置CORS头
		// Access-Control-Allow-Origin 只能是单个源或 "*"，不能是逗号分隔的列表
		if len(corsConfig.AllowOrigins) > 0 {
			if corsConfig.AllowOrigins[0] == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				// 检查请求的 Origin 是否在允许列表中
				if slices.Contains(corsConfig.AllowOrigins, origin) {
					c.Header("Access-Control-Allow-Origin", origin)
				}
			}
		}

		c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowMethods, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowHeaders, ","))
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
