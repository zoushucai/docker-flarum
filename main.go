package main

import (
	"embed"
	"fmt"

	"mygo/dbpk"
	"mygo/logger"
	"mygo/routers"
	"mygo/settings"
)

/*
	该路径是相对于go:embed 声明所属的文件进行寻址，不是根据go mod ,也不是main.go. 而且还不能使用 .. 这种路径
*/
//go:embed static/favicon.ico
var fav embed.FS

// @title			API 文档
// @version		1.0
// @description	API 详细的描述信息 ......
// @BasePath		/

// @contact.name	zsc
// @contact.url	https://github.com/zoushucai/random-images
// @contact.email	zscmoyujian@gmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	// 1.加载配置 (默认只要引入 settings 包即可,自动初始化)

	// 2.初始化日志
	if err := logger.Init(); err != nil {
		panic(err)
	}
	// 3.1 初始化数据 -- 是否重新写入新的数据  (长用于测试开发)
	ip, err := dbpk.InitData()
	if err != nil {
		panic(err)
	}
	// ip.SaveToJSON("data.json")
	// ip.SaveToCSV("data.csv")

	// 5. 注册路由
	r := routers.Setup(fav) // fav 网站图标,嵌入二进制
	routers.RegisterRoutes(r, ip)
	// 启动服务
	r.Run(fmt.Sprintf("0.0.0.0:%v", settings.App.Port))
	// Swagger 路由
	// docs.SwaggerInfo.Host = fmt.Sprintf("0.0.0.0:%v", settings.App.Port)
}
