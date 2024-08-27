package dbpk

import (
	"fmt"
	"log/slog"
	"mygo/settings"
	"mygo/utils"
)

/*
这个文件,主要用于处理数据,把数据直接读取到内存中,而不是写入到数据库,然后从数据库在读取
毕竟尽可能简化
*/

func InitData() (*utils.ImageProcessor, error) {
	ip := &utils.ImageProcessor{
		Dir:       settings.App.ImgPath,
		BatchSize: settings.App.BatchSize,
	}
	err := ip.GetImagesInfo() // 获取图片信息

	if err != nil {
		return nil, fmt.Errorf("无法读取图片信息: %v", err)
	}
	slog.Info("成功读取图片的信息",
		slog.String("图片路径:", ip.Dir),
		slog.Int("图片数量:", len(*ip.SaveDb)),
	)

	return ip, nil
}
