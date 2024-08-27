package middlewares

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// 排除不需要记录的请求，如静态资源请求
	whitelist = []string{
		"/v1/random",
		"/v1/upload",
	}
	blacklist = []string{
		"/favicon.ico",
		"/favicon.png",
		"/static/",
		"/js/",
		"/css/",
		"/swagger/",
	}
)

// 直接用 slog 记录日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		// 跳过不需要记录的路径
		if !shouldLog(c.Request.URL.Path, whitelist, blacklist) {
			return
		}

		statusCode := c.Writer.Status()
		latency := time.Since(startTime)

		// 准备日志字段
		slog.With(
			slog.Int64("ts", startTime.UnixMilli()),               // 时间戳, 毫秒
			slog.Int("status", statusCode),                        // 状态码
			slog.Duration("latency", latency),                     // 请求耗时
			slog.String("method", c.Request.Method),               // 请求方法
			slog.String("path", c.Request.URL.Path),               // 请求路径
			slog.String("query", c.Request.URL.RawQuery),          // 请求参数
			slog.String("ip", c.ClientIP()),                       // 客户端 IP
			slog.String("user_agent", c.Request.UserAgent()),      // 请求 User-Agent
			slog.String("accept", c.Request.Header.Get("Accept")), // Accept
		).Info("Request received")

	}
}

// shouldLog 判断请求路径是否需要被记录
func shouldLog(path string, whitelist, blacklist []string) bool {
	// 检查路径是否在白名单中 -- 需要记录
	for _, whitePath := range whitelist {
		if strings.HasPrefix(path, whitePath) {
			return true
		}
	}

	// 检查路径是否在黑名单中 -- 不需要记录
	for _, blackPath := range blacklist {
		if strings.HasPrefix(path, blackPath) {
			return false
		}
	}

	// 默认情况下都需要记录
	return true
}
