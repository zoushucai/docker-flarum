package routers

import (
	"embed"
	"mygo/middlewares"
	"mygo/utils"
	"net/http"

	_ "mygo/docs"

	"github.com/fufuok/favicon"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup 初始化 Gin 引擎
func Setup(fav embed.FS) *gin.Engine {
	// r := gin.Default()

	r := gin.New()              // 初始化 Gin 引擎
	r.Use(middlewares.Logger()) // 使用 slog 日志中间件
	r.Use(gin.Recovery())       // 使用 Gin 内置的恢复中间件

	// 使用 CORS 中间件 --- 自定义的 CORS 中间件
	// r.Use(middlewares.Cors()) // 使用 CORS 中间件--解决跨域

	// 使用 cors 中间件--- 第三方包
	r.Use(cors.Default())

	// 设置可信代理 IP
	r.SetTrustedProxies([]string{
		"127.0.0.1",      // 本地 IP4
		"::1",            // 本地 IPv6
		"10.0.0.0/8",     // 内网
		"172.16.0.0/12",  // 内网
		"192.168.0.0/16", // 内网
	})
	// 设置最大上传文件大小
	r.MaxMultipartMemory = 10 << 20 // 10MB
	r.Use(Favicon(fav))
	return r
}

func Favicon(fav embed.FS) gin.HandlerFunc {
	return favicon.New(favicon.Config{
		File:       "static/favicon.ico",
		FileSystem: http.FS(fav),
	})
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine, ip *utils.ImageProcessor) {
	// 第一个是 url 的路径, 第二个是计算机中的路径, 二者建立一个一一映射
	r.Static(ip.Dir, ip.Dir)

	// 设置 favicon -- 这只是提供了文件访问的功能，不提嵌入的功能
	// ginServer.StaticFile("/favicon.ico", "./static/favicon.ico")
	// ginServer.StaticFile("/favicon.png", "./static/favicon.png")

	// 重定向 "/" --> v1/random
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/v1/random")
	})

	// API 路由
	v1 := r.Group("/v1")
	{
		// 传统的传参
		v1.GET("/random", RandomImage(ip)) // 随机获取图片
		v1.POST("/upload", FileUpload(ip)) // 上传图片
	}
	// 添加 404 路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "404 Not Found"})
	})

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
