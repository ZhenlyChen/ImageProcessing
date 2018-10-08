package main

import (
	"github.com/ZhenlyChen/ImageProcessing/EightBit"
	"github.com/ZhenlyChen/ImageProcessing/IrisImage"
	"github.com/andybons/gogif"
	"gocv.io/x/gocv"
	"image"
	"image/gif"
	"image/jpeg"
	"os"
	"strconv"
)

func main() {
	process8Bit()
	// processIris()
}

func process8Bit() {
	img := gocv.IMRead("./img/redapple.jpg", gocv.IMReadColor)
	resImg := EightBit.To8Bit(img)
	distFile, err := os.Create("./dist/goodapple.jpg")
	check(err)
	dist, err := resImg.ToImage()
	check(err)
	// png.Encode(distFile, dist)
	err = jpeg.Encode(distFile, dist, &jpeg.Options{Quality: 100})
	check(err)
}

func processIris() {
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
		resImg, err := IrisImage.GetIrisImage(img1, img2, i)
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
