package makeGif

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"math"
	"os"

	"github.com/andybons/gogif"
)

func Make(src, dst string) {
	// 打开原图
	fmt.Println("Encode for GIF: ", src)
	file1, err := os.OpenFile(src, os.O_RDONLY, 0)
	check(err)
	// JPEG 解码
	img, err := jpeg.Decode(file1)
	check(err)
	rgbSrc := getRGB(img)
	distGif, err := os.Create(dst)
	check(err)
	// 生成 GIF
	outGif := &gif.GIF{}
	bounds := img.Bounds()
	// 生成 调色板
	palettedImage := image.NewPaletted(bounds, nil)
	// 中值取值算法计算256色
	quantizer := gogif.MedianCutQuantizer{NumColor: 256}
	// 量化
	quantizer.Quantize(palettedImage, bounds, img, image.ZP)
	// 加入GIF数据
	outGif.Image = append(outGif.Image, palettedImage)
	outGif.Delay = append(outGif.Delay, 0)
	// GIF 编码
	err = gif.EncodeAll(distGif, outGif)
	if err != nil {
		fmt.Println(err)
	}
	file1.Close()
	fmt.Println("Encode for GIF finish: ", dst)


	fmt.Println("Calculate MSE for ", dst)
	file2, err := os.OpenFile(dst, os.O_RDONLY, 0)
	check(err)
	imgGif, err := gif.Decode(file2)
	check(err)
	rgbDst := getRGB(imgGif)

	// 均方差计算
	sum := 0.0
	count := 0
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			for k := 0; k < 3; k ++ {
				sum += math.Pow(float64(rgbSrc[x][y][k] - rgbDst[x][y][k]), 2.0)
				count++
			}
		}
	}
	fmt.Println("MSE: ", float64(sum) / float64(count))
}

func getRGB(src image.Image) [][][3]int {
	bounds := src.Bounds().Max
	mat := make([][][3]int, bounds.X)
	for x := 0; x < bounds.X; x++ {
		row := make([][3]int, bounds.Y)
		for y := 0; y < bounds.Y; y++ {
			r,g,b,_ := src.At(x, y).RGBA()
			R, G, B := float32(r/257), float32(g/257), float32(b/257)
			row[y][0], row[y][1], row[y][2] =  int(R), int(G), int(B)
		}
		mat[x] = row
	}
	return mat
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
