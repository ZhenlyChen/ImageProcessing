package main

import (
	"fmt"
	"github.com/ZhenlyChen/ImageProcessing/IrisImage"
	"gocv.io/x/gocv"
	"image"
	"image/jpeg"
	"os"
	"strconv"
)

func main() {
	processIris()
}

func processIris() {
	var err error
	var file1, file2 *os.File
	var img1, img2 image.Image
	file1, err = os.OpenFile("./img/Nobel.jpg",os.O_RDONLY,0)
	check(err)
	defer file1.Close()
	file2, err = os.OpenFile("./img/lena.jpg", os.O_RDONLY, 0)
	check(err)
	defer file2.Close()
	img1, err = jpeg.Decode(file1)
	check(err)
	img2, err = jpeg.Decode(file2)
	check(err)


	var names []string
	for i := 0.01; i <= 1; i += 0.01 {
		resImg ,err := IrisImage.GetIrisImage(img1, img2, i)
		check(err)
		name := "./dist/iris"+strconv.Itoa(int(i * 100))+".jpg"
		names = append(names, name)
		fmt.Println(name)
		distFile, err := os.Create(name)
		check(err)
		err = jpeg.Encode(distFile, resImg, &jpeg.Options{Quality: 100})
		check(err)
	}

	video, err := gocv.VideoWriterFile("./dist/iris.avi", "MJPG", 60, img1.Bounds().Dy(), img1.Bounds().Dx(), true)
	check(err)
	for _, n := range names {
		img := gocv.IMRead(n, gocv.IMReadColor)
		video.Write(img)
	}
	window := gocv.NewWindow("Video")
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
	v.Close()
	fmt.Println("success!")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}