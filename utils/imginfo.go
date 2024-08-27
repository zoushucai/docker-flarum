package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg" // 支持 JPEG 格式, 不然解析图片失败
	_ "image/png"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"math/rand/v2"
)

type ImageInfo struct {
	Sub    string `json:"sub" csv:"sub"`       // 文件所在的目录 (不含文件名),  Sub + File 则为文件路径
	File   string `json:"file" csv:"file"`     // 文件名(含后缀)
	Suf    string `json:"suf" csv:"suf"`       // 后缀名
	Width  int    `json:"width" csv:"width"`   // 图片文件的宽
	Height int    `json:"height" csv:"height"` // 图片文件的高
}

// 添加一个属性	SaveDb  , 用来保存 []*ImageInfo 的地址
type ImageProcessor struct {
	Dir       string        // 要处理的目录
	SaveDb    *[]*ImageInfo // 保存图像信息的指针. 默认为 nil
	BatchSize int           // 批量插入的大小
}

// 定义常见图片扩展名
var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".tiff": true,
	".webp": true,
}

// GetImagesInfo 获取目录下所有图片的详细信息
func (ip *ImageProcessor) GetImagesInfo() error {
	paths, err := ip.GetImagesPaths() // 获取图片路径
	if err != nil {
		return fmt.Errorf("failed to get image paths: %v", err)
	}
	// 提前分配 infos 切片的容量以提升性能
	infos := make([]*ImageInfo, 0, len(paths))
	for _, path := range paths {
		if path == "" {
			continue // 忽略可能的空路径
		}
		info, err := ip.GetSingleImageInfo(path)
		if err != nil {
			slog.Warn(" get info failed", slog.String("path", path), slog.Any("err", err))
			continue
		} else {
			infos = append(infos, info)
		}

	}
	ip.SaveDb = &infos
	return nil

}

// GetSingleImageInfo 获取单个图片的详细信息
func (ip *ImageProcessor) GetSingleImageInfo(path string) (*ImageInfo, error) {
	// 获取文件名、后缀名和所在目录
	fileName := filepath.Base(path)
	suf := filepath.Ext(path)
	sub := filepath.Dir(path)

	// 获取图片的宽高
	width, height, err := ip.GetImageSize(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get image size for %s: %v", path, err)
	}
	return &ImageInfo{
		Sub:    sub,
		File:   fileName,
		Suf:    suf,
		Width:  width,
		Height: height,
	}, nil
}

// GetImagesPaths 遍历指定目录，返回所有图像文件的路径
//
// 如果遍历目录时遇到错误，则返回该错误。函数会检查文件扩展名是否在 imageExtensions 中，
// 只有符合条件的文件路径才会被添加到结果切片中。
//
// 参数:
//   - 无
//
// 返回值:
//   - []string: 包含所有图像文件路径的切片
//   - error: 遇到的错误，如果没有错误则返回 nil
func (ip *ImageProcessor) GetImagesPaths() (paths []string, err error) {

	// 使用 filepath.WalkDir 遍历目录
	err = filepath.WalkDir(ip.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // 返回遍历中的错误
		}

		// 如果是文件且具有图像扩展名，则将其路径添加到切片中
		if !d.IsDir() && imageExtensions[strings.ToLower(filepath.Ext(path))] {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %v", ip.Dir, err)
	}

	return paths, nil
}

// GetImageSize 获取图片的宽高
func (ip *ImageProcessor) GetImageSize(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	// 解码图片文件
	imgConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	// 获取图片的宽和高
	return imgConfig.Width, imgConfig.Height, nil
}

// CalculateMd5 计算文件的 MD5 哈希值
func (ip *ImageProcessor) CalculateMd5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

/*
	由于没有用到数据库, 因此需要自己写方法, 对数据进行增删改查的方法
	这里只用了查,因此只写查的方法
*/

// // SaveToCSV 将图片信息保存为 CSV 文件
// func (ip *ImageProcessor) SaveToCSV(filename string) error {
// 	DB := ip.SaveDb
// 	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
// 	if err != nil {
// 		return fmt.Errorf("failed to open file %s: %v", filename, err)
// 	}
// 	defer file.Close()

// 	// 使用 gocsv 库写入 CSV 文件
// 	if err := gocsv.MarshalFile(DB, file); err != nil {
// 		return fmt.Errorf("failed to write to file %s: %v", filename, err)
// 	}
// 	return nil
// }

// // SaveToJSON 将图片信息保存为 JSON 文件
// func (ip *ImageProcessor) SaveToJSON(filename string) error {
// 	DB := ip.SaveDb // 获取保存数据的指针
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file %s: %v", filename, err)
// 	}
// 	defer file.Close()
// 	encoder := json.NewEncoder(file)
// 	encoder.SetIndent("", "  ") // 美化输出

// 	if err := encoder.Encode(DB); err != nil {
// 		return fmt.Errorf("failed to encode data to JSON: %v", err)
// 	}
// 	return nil
// }

// FilterBySub 根据指定的 Sub 过滤图像信息 (可以链式调用)
// (这样有一个硬伤,会把SaveDb清空,或者查询查询这数据没了,因此需要创建一个副本)
func (ip *ImageProcessor) FilterBySub(sub string) *ImageProcessor {
	if ip.SaveDb == nil {
		ip.SaveDb = &[]*ImageInfo{} // 如果 SaveDb 为 nil，则初始化为空切片
		return ip
	}

	if sub == "" {
		return ip // 如果 sub 为空，返回当前对象
	}

	var filteredInfos []*ImageInfo
	for _, info := range *ip.SaveDb {
		if strings.HasSuffix(info.Sub, sub) {
			filteredInfos = append(filteredInfos, info)
		}
	}
	// 为了防止修改 SaveDb，创建一个新的 ImageProcessor 对象
	newip := &ImageProcessor{
		BatchSize: ip.BatchSize,
		Dir:       ip.Dir,
		SaveDb:    &filteredInfos,
	}
	return newip
}

// RandomImageInfo 随机返回 SaveDb 中的一条 ImageInfo
func (ip *ImageProcessor) RandomImageInfo() *ImageInfo {
	if ip.SaveDb == nil || len(*ip.SaveDb) == 0 {
		return nil
	}
	// 使用当前时间作为随机数生成器的种子
	// 生成一个随机索引
	fmt.Println("len(*ip.SaveDb): ", len(*ip.SaveDb))
	randomIndex := rand.IntN(len(*ip.SaveDb))
	// 返回随机选择的 ImageInfo
	return (*ip.SaveDb)[randomIndex]
}

// GetByIndex 根据单个下标或下标切片返回数据，并更新结构体以支持链式调用
// (这样有一个硬伤,会把SaveDb清空,或者查询查询这数据没了)
func (ip *ImageProcessor) GetByIndex(indices interface{}) *ImageProcessor {
	if ip.SaveDb == nil || len(*ip.SaveDb) == 0 {
		return ip
	}
	var filteredInfos []*ImageInfo
	switch v := indices.(type) {
	case int:
		// 单个下标
		if v >= 0 && v < len(*ip.SaveDb) {
			filteredInfos = append(filteredInfos, (*ip.SaveDb)[v])
		}
	case []int:
		// 下标切片
		for _, index := range v {
			if index >= 0 && index < len(*ip.SaveDb) {
				filteredInfos = append(filteredInfos, (*ip.SaveDb)[index])
			}
		}
	default:
		return ip // 如果 indices 类型无效，返回原始结构体
	}
	// 为了防止修改 SaveDb，创建一个新的 ImageProcessor 对象
	newip := &ImageProcessor{
		BatchSize: ip.BatchSize,
		Dir:       ip.Dir,
		SaveDb:    &filteredInfos,
	}
	return newip
}

// FirstInfo 返回 SaveDb 的第一个元素，如果 SaveDb 为 nil 或空，则返回错误信息
func (ip *ImageProcessor) FirstInfo() *ImageInfo {
	if ip.SaveDb == nil || len(*ip.SaveDb) == 0 {
		panic("no image info available")
	}
	return (*ip.SaveDb)[0]
}
