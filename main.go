package main

import (
	"github.com/ZhenlyChen/ImageProcessing/makeJPEG"
)

func main() {
	// 八位图像压缩
	// EightBit.Process8Bit()
	// 图像切换效果
	// IrisImage.ProcessIris()
	// makeGif.Make( "./img/photo.jpg", "./dist/photo.gif")
	// makeGif.Make( "./img/cartoon.jpg", "./dist/cartoon.gif")
	// JPEG编码
	makeJPEG.Make("./img/photo.jpg", "./dist/photo.jpg")
}
