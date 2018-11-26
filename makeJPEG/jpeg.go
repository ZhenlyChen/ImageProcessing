package makeJPEG

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// 系数
type factor struct {
	Length int // 长度
	Data   int // 数据
}

func Make(src, dst string) {
	file1, err := os.OpenFile(src, os.O_RDONLY, 0)
	check(err)
	defer file1.Close()
	img, err := jpeg.Decode(file1)
	check(err)
	err = encode(img, dst)
	check(err)
}

var matrixY = [8][8]int{
	{16, 11, 10, 16, 24, 40, 51, 61},
	{12, 12, 14, 19, 26, 58, 60, 55},
	{14, 13, 16, 24, 40, 57, 69, 56},
	{14, 17, 22, 29, 51, 87, 80, 62},
	{18, 22, 37, 56, 68, 109, 103, 77},
	{24, 35, 55, 64, 81, 104, 113, 92},
	{49, 64, 78, 87, 103, 121, 120, 101},
	{72, 92, 95, 98, 112, 100, 103, 99},
}

var matrixIQ = [8][8]int{
	{17, 18, 24, 47, 99, 99, 99, 99},
	{18, 21, 26, 66, 99, 99, 99, 99},
	{24, 26, 56, 99, 99, 99, 99, 99},
	{47, 66, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
}

func encode(src image.Image, dst string) error {
	fmt.Println("Encode JPEG for ", dst)
	fmt.Println("RGB to YUV")
	imgColor := convertToYIQ(src)
	fmt.Println("Size: ", len(imgColor), len(imgColor[0]))
	// AC[第i块数据][0:Y, 1:I, 2:Q][第i对AC系数]
	var AC [][][]factor
	// DC[第i块数据][0:Y, 1:I, 2:Q]
	var DC [][3]int
	fmt.Println("DTC...")
	fmt.Println("Quantitative...")
	fmt.Println("ZigZag...")
	fmt.Println("DPCM...")
	for x := 0; x < len(imgColor); x += 8 {
		for y := 0; y < len(imgColor[x]); y += 8 {
			convertDTC(imgColor, x, y)
			quantitative(imgColor, x, y)
			ac, dc := traverseZigZag(imgColor, x, y)
			AC = append(AC, ac)
			DC = append(DC, dc)
		}
	}
	// allDC[0:Y, 1:I, 2:Q][第i块数据]
	allDC := dcDPCM(DC)
	fmt.Println("Huffman for DC...")
	DCBinary := huffmanDC(allDC)
	fmt.Println("Huffman for AC...")
	ACBinary := huffmanAC(AC)
	fmt.Println("Output DC...")
	err := ioutil.WriteFile(dst+".dc", []byte(fmt.Sprint(allDC)), 0644)
	check(err)
	fmt.Println("Output AC...")
	err = ioutil.WriteFile(dst+".ac", []byte(fmt.Sprint(AC)), 0644)
	check(err)
	fmt.Println("Output Binary...")
	var finalData []byte
	finalData = append(finalData, DCBinary...)
	finalData = append(finalData, ACBinary...)
	err = ioutil.WriteFile(dst+".binary", finalData, 0644)
	check(err)

	// 解码
	fmt.Println("Decode...")
	fmt.Println("Inverse quantitative...")
	fmt.Println("Inverse DTC...")
	for x := 0; x < len(imgColor); x += 8 {
		for y := 0; y < len(imgColor[x]); y += 8 {
			Iquantitative(imgColor, x, y)
			convertIDTC(imgColor, x, y)
		}
	}
	fmt.Println("YUV to RGB...")
	convertRGB(imgColor)


	fmt.Println("Write Image...")
	bounds := src.Bounds().Max
	dstImage := image.NewRGBA(src.Bounds())
	for x := 0; x < bounds.Y; x++ {
		for y := 0; y < bounds.X; y++ {
			dstImage.Set(y, x, color.RGBA{
				R: uint8(imgColor[x][y][0]),
				G: uint8(imgColor[x][y][1]),
				B: uint8(imgColor[x][y][2]), A: 255 })
		}
	}
	distFile, err := os.Create(dst)
	check(err)
	fmt.Println("Write File...")
	err = jpeg.Encode(distFile, dstImage, &jpeg.Options{Quality: 75})
	check(err)

	fmt.Println("Finish!")
	return nil
}

