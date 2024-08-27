package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mygo/settings"

	"github.com/natefinch/lumberjack"
)

var logWriter *lumberjack.Logger

// EnsureFileExists 检查文件是否存在。如果不存在，则创建该文件及其父目录。
func EnsureFileExists(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		// 文件已存在
		return nil
	}
	// 文件不存在，确保父目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	file.Close() // 关闭文件

	return nil
}

func Init() error {
	EnsureFileExists(settings.APPCONFIG.Log.Filename)
	// 从配置文件中读取日志配置
	logWriter = &lumberjack.Logger{
		Filename:   settings.APPCONFIG.Log.Filename,   // 日志文件路径
		MaxSize:    settings.APPCONFIG.Log.MaxSize,    // 文件最大尺寸（以MB为单位）
		MaxBackups: settings.APPCONFIG.Log.MaxBackups, // 保留的最大旧文件数量
		MaxAge:     settings.APPCONFIG.Log.MaxAge,     // 保留旧文件的最大天数
		Compress:   settings.APPCONFIG.Log.Compress,   // 是否压缩/归档旧文件
		LocalTime:  true,                              // 使用本地时间创建时间戳
	}

	opt := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format(time.DateTime))
			}
			if a.Key == slog.SourceKey {
				// 修改调用的文件, 采用相对路径
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
				a.Value = slog.AnyValue(source) // 使用相对的路径
			}
			return a
		},
	}
	var handler slog.Handler
	if strings.ToLower(settings.APPCONFIG.Log.Format) == "text" {
		handler = slog.NewTextHandler(logWriter, nil)
	} else {
		handler = slog.NewJSONHandler(logWriter, opt)
	}
	slog.SetDefault(slog.New(handler))
	return nil
}
