package IrisImage

import (
	"github.com/andybons/gogif"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"math"
	"os"
	"strconv"
)


func ProcessIris() {
	// 读取源图像
	var err error
	var file1, file2 *os.File
	var img1, img2 image.Image
	file1, err = os.OpenFile("./img/Nobel.jpg", os.O_RDONLY, 0)
	check(err)
	defer file1.Close()
	file2, err = os.OpenFile("./img/lena.jpg", os.O_RDONLY, 0)
	check(err)
	defer file2.Close()
	img1, err = jpeg.Decode(file1)
	check(err)
	img2, err = jpeg.Decode(file2)
	check(err)

	// 生成不同半径的图片
	var names []string
	var subimages []image.Image
	for i := 0.01; i <= 1; i += 0.01 {
		resImg, err := GetIrisImage(img1, img2, i)
		subimages = append(subimages, resImg)
		check(err)
		name := "./dist/iris" + strconv.Itoa(int(i*100)) + ".jpg"
		names = append(names, name)
		distFile, err := os.Create(name)
		check(err)
		err = jpeg.Encode(distFile, resImg, &jpeg.Options{Quality: 100})
		check(err)
	}

	// 生成gif
	distGif, err := os.Create("./dist/iris.gif")
	check(err)
	outGif := &gif.GIF{}
	for _, simage := range subimages {
		bounds := simage.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)

		quantizer := gogif.MedianCutQuantizer{NumColor: 64}
		quantizer.Quantize(palettedImage, bounds, simage, image.ZP)

		// Add new frame to animated GIF
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 0)
	}
	gif.EncodeAll(distGif, outGif)

	// 生成video
	video, err := gocv.VideoWriterFile("./dist/iris.avi", "MJPG", 60, img1.Bounds().Dy(), img1.Bounds().Dx(), true)
	check(err)
	for _, n := range names {
		img := gocv.IMRead(n, gocv.IMReadColor)
		video.Write(img)
	}
	video.Close()
	/*window := gocv.NewWindow("Video")
	v, err := gocv.VideoCaptureFile("./dist/iris.avi")
	check(err)
	frame := gocv.NewMat()
	for {
		v.Read(&frame)
		if frame.Empty() {
			break
		}
		window.IMShow(frame)
		if window.WaitKey(33) >= 0 {
			break
		}
	}
	v.Close()*/
	// fmt.Println("success!")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}


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
