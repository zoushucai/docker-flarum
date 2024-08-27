package settings

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
)

/*
为了避免引入太多的第三方库,以及配置信息, 直接把配置信息用结构体来表示
*/

// AppConfig 定义应用的配置信息
type AppConfig struct {
	Name      string // 应用名称
	Version   string // 应用版本
	Mode      string // 环境名称
	Port      int    // 应用端口
	ImgPath   string // 图片路径
	BatchSize int    // 批量操作数量
}

// LogConfig 定义日志的配置信息
type LogConfig struct {
	Level      string // 日志级别
	Filename   string // 日志文件名
	MaxSize    int    // 日志文件最大保存的大小 (MB)
	MaxAge     int    // 日志文件最多保存多少天
	MaxBackups int    // 日志文件最多保存多少个备份
	Compress   bool   // 是否压缩日志文件
	Format     string // 日志格式 (json or text)
}

// Config 包含所有配置的根结构体
type Config struct {
	App AppConfig // 应用配置
	Log LogConfig // 日志配置
}

// 全局配置变量
var APPCONFIG = Config{
	App: AppConfig{
		Name:      "webapp",
		Version:   "1.0.0",
		Mode:      "dev",
		Port:      9113,
		ImgPath:   "./assets", // 图片的路径
		BatchSize: 1000,       //(未使用)
	},
	Log: LogConfig{
		Level:      "debug",           // 日志级别(没有用到)
		Filename:   "logs/webapp.log", // 日志文件名
		MaxSize:    10,                // 日志文件最大保存的大小 (MB)
		MaxAge:     30,                // 日志文件最多保存多少天
		MaxBackups: 7,                 // 日志文件最多保存多少个备份
		Compress:   true,              // 是否压缩日志文件
		Format:     "json",            // 日志格式 (json or text)
	},
}
var App = &APPCONFIG.App
var Log = &APPCONFIG.Log

/*
ParseFlags 解析命令行参数, 返回参数指针
flag.Int 语法的定义命令行参数, 格式: 参数, 默认值, 注释
*/
func ParseFlags() {

	port := flag.Int("port", App.Port, "Port to run the server on")
	imgDir := flag.String("imgdir", App.ImgPath, "Directory for images") // 定义 imgdir 参数，默认值为 ./assets
	flag.Parse()                                                         // 解析命令行参数
	// return port, imgDir
	// 使用命令行参数覆盖默认配置
	App.Port = *port
	App.ImgPath = filepath.Clean(*imgDir)
}

// 打印路径信息
func PrintPathInfo() {
	workdir, err := os.Getwd()
	if err != nil {
		return
	}
	slog.Info("当前工作目录:", slog.String("路径", workdir))
	exePath, err := os.Executable()
	if err != nil {
		return
	}
	slog.Info("当前程序路径:", slog.String("路径", exePath))
}

func init() {
	ParseFlags()
	PrintPathInfo()
}
