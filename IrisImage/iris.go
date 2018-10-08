package IrisImage

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func GetIrisImage(img1, img2 image.Image, per float64) (image.Image, error) {
	// 图片色彩模型为rCbCr，先统一转换为RGBA
	b1 := img1.Bounds()
	m1 := image.NewRGBA(image.Rect(0, 0, b1.Dx(), b1.Dy()))
	draw.Draw(m1, m1.Bounds(), img1, b1.Min, draw.Src)

	b2 := img2.Bounds()
	img2.At(1, 2)
	m2 := image.NewRGBA(image.Rect(0, 0, b2.Dx(), b2.Dy()))
	draw.Draw(m2, m2.Bounds(), img2, b2.Min, draw.Src)

	res := image.NewGray(img1.Bounds())
	// 计算圆心
	midX := img1.Bounds().Dx() / 2
	midY := img1.Bounds().Dy() / 2
	// 计算半径
	radius := per * getDis(0, 0, midX, midY)
	for x := 0; x < img1.Bounds().Dx(); x++ {
		for y := 0; y < img1.Bounds().Dy(); y++ {
			var thisColor color.Color
			// 判断距离
			if getDis(x, y, midX, midY) < radius {
				thisColor = m2.At(x, y)
			} else {
				thisColor = m1.At(x, y)
			}
			// 填充颜色
			red, _, _, _ := thisColor.RGBA()
			res.Set(x, y, color.Gray{Y: uint8(red)})
		}
	}
	return res, nil
}

func getDis(x1, y1, x2, y2 int) float64 {
	lenX := math.Abs(float64(x1 - x2))
	lenY := math.Abs(float64(y1 - y2))
	return math.Sqrt(float64(lenX*lenX + lenY*lenY))
}
