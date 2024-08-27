package pics

/*
简单的图像处理
*/
import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/anthonynsimon/bild/transform"
)

// ProcessImage 读取图像，进行缩放，并返回图像的字节切片
func ProcessImage(path string, width, height int) ([]byte, error) {
	// 打开图像文件
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	// 解码图像
	img, imgType, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	// 缩放图像
	var newImage image.Image
	if width == 0 && height == 0 {
		newImage = img
	} else {
		if width == 0 {
			// 如果只指定高度，按比例计算宽度
			ratio := float64(height) / float64(img.Bounds().Dy())
			width = int(ratio * float64(img.Bounds().Dx()))
		} else if height == 0 {
			// 如果只指定宽度，按比例计算高度
			ratio := float64(width) / float64(img.Bounds().Dx())
			height = int(ratio * float64(img.Bounds().Dy()))
		}
		// NearestNeighbor
		newImage = transform.Resize(img, width, height, transform.Lanczos)
	}

	// 编码图像为字节切片
	var buffer bytes.Buffer
	switch imgType {
	case "jpeg":
		err = jpeg.Encode(&buffer, newImage, nil)
	case "png":
		err = png.Encode(&buffer, newImage)
	case "gif":
		err = gif.Encode(&buffer, newImage, nil)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", imgType)
	}

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
