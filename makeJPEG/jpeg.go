package makeJPEG

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// AC系数
type factor struct {
	Length int // 0的长度
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
	imgColor := convertToYIQ(src)
	fmt.Println(len(imgColor), len(imgColor[0]))

	// AC[第i块数据][0:Y, 1:I, 2:Q][第i对AC系数]
	var AC [][][]factor
	// DC[第i块数据][0:Y, 1:I, 2:Q]
	var DC [][3]int

	for x := 0; x < len(imgColor); x += 8 {
		for y := 0; y < len(imgColor[x]); y += 8 {
			dtc := convertDTC(imgColor, x, y)
			dtcAfterQ := quantitative(dtc, x, y)
			ac, dc := traverseZigZag(dtcAfterQ, x, y)
			AC = append(AC, ac)
			DC = append(DC, dc)
		}
	}
	// allDC[0:Y, 1:I, 2:Q][第i块数据]
	allDC := dcDPCM(DC)
	DCBinary := huffmanDC(allDC)
	ACBinary := huffmanAC(AC)
	var finalData []byte
	finalData = append(finalData, DCBinary...)
	finalData = append(finalData, ACBinary...)
	return ioutil.WriteFile(dst, finalData, 0644)
}

func huffmanDC(dc [3][]int) []byte {
	table := make(map[int]int)
	var data []factor
	for colorChannel := range dc {
		for _, d := range dc[colorChannel] {
			size := getLength(d)
			table[size]++
			data = append(data, factor{size, d})
		}
	}
	dcTable := huffmanEncode(table)
	bitData := bitArray{}
	for i := range data {
		dc := data[i]
		sizeStr := dcTable[dc.Length]
		sizeByte := []byte(sizeStr)
		for _, b := range sizeByte {
			bitData.addBit(b - '0')
		}
		bitData.addData(dc.Data)
	}
	return bitData.Data
}

func huffmanAC(ac [][][]factor) []byte {
	table := make(map[int]int)
	var data []factor
	for i := range ac {
		for colorChannel := range ac[i] {
			for _, d := range ac[i][colorChannel] {
				size := getLength(d.Data)
				symbol1 := d.Length * 16 + size
				table[symbol1]++
				data = append(data, factor{symbol1, d.Data})
			}
		}
	}
	dcTable := huffmanEncode(table)
	bitData := bitArray{}
	for i := range data {
		dc := data[i]
		sizeStr := dcTable[dc.Length]
		sizeByte := []byte(sizeStr)
		for _, b := range sizeByte {
			bitData.addBit(b - '0')
		}
		bitData.addData(dc.Data)
	}
	return bitData.Data
}



func getLength(num int) int {
	num %= 2048
	if num < 0 {
		num = -num
	}
	l := 0
	for num > 0 {
		num >>= 1
		l++
	}
	return l
}


// DPCM编码
func dcDPCM(mat [][3]int) [3][]int {
	var dst [3][]int
	for k := 0; k < 3; k++ {
		var row []int
		last := mat[0][k]
		for i := range mat {
			if i == 0 {
				row = append(row, last)
			} else {
				last = mat[i][k]
				row = append(row, mat[i][k]-last)
			}
		}
		dst[k] = row
	}
	return dst
}

// 量化
func quantitative(src [][][3]int, x, y int) [][][3]int {
	// Y
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			src[x+i][y+j][0] = int(math.Round(float64(src[x+i][y+j][0]) / float64(matrixY[i][j])))
		}
	}
	// IQ
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			src[x+i][y+j][1] = int(math.Round(float64(src[x+i][y+j][1]) / float64(matrixIQ[i][j])))
			src[x+i][y+j][2] = int(math.Round(float64(src[x+i][y+j][2]) / float64(matrixIQ[i][j])))
		}
	}
	return src
}

// DTC离散余弦变换
func convertDTC(src [][][3]int, x, y int) [][][3]int {
	cSrc := make([][][3]int, 8)
	for k := 0; k < 8; k++ {
		cSrc[k] = make([][3]int, 8)
		copy(cSrc[k], src[x+k][y:y+8])
	}
	for k := 0; k < 3; k++ {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				cU, cV := 1.0, 1.0
				if i == 0 {
					cU = math.Sqrt(2) / 2.0
				}
				if j == 0 {
					cV = math.Sqrt(2) / 2.0
				}
				res := (cU * cV) / 4.0
				var value float64
				for a := 0; a < 8; a++ {
					for b := 0; b < 8; b++ {
						value += math.Cos(float64(2*a+1)*float64(i)*float64(math.Pi)/16.0) *
							math.Cos(float64(2*b+1)*float64(j)*float64(math.Pi)/16.0) *
							float64(cSrc[a][b][k])
					}
				}
				res = res * value
				src[x+i][y+j][k] = int(math.Round(res))
			}
		}
	}
	return src
}

// RGB to YIQ and reSampling
func convertToYIQ(src image.Image) [][][3]int {
	bounds := src.Bounds().Max
	mat := make([][][3]int, bounds.X+(8-bounds.X%8))
	for x := 0; x < bounds.X; x++ {
		row := make([][3]int, bounds.Y+(8-bounds.Y%8))
		for y := 0; y < bounds.Y; y++ {
			r, g, b, _ := src.At(x, y).RGBA()
			R, G, B := float32(r), float32(g), float32(b)
			row[y][0] = int(0.299*R + 0.587*G + 0.114*B)
			// 4:2:0 二次采样
			if y%2 == 0 {
				row[y][1] = int(0.596*R - 0.275*G - 0.321*B)
				row[y][2] = int(0.212*R - 0.523*G + 0.311*B)
				if x%2 == 1 {
					row[y][1] = (int(mat[x-1][y][1]) + int(row[y][1])) / 2
					row[y][2] = (int(mat[x-1][y][2]) + int(row[y][2])) / 2
					mat[x-1][y][1] = 0
					mat[x-1][y][2] = 0
				}
			}
		}
		mat[x] = row
	}
	return mat
}

// traverseZigZag 游长编码 ac dc
func traverseZigZag(src [][][3]int, x, y int) ([][]factor, [3]int) {
	res := make([][]factor, 3)
	var dc [3]int
	for k := 0; k < 3; k ++ {
		var zig []int
		i, j, up := 0, 0, 1
		turned := false
		d := [2][2]int{{1, -1}, {-1, 1}}
		corner := [2][4]int{{1, 0, 0, 1}, {0, 1, 1, 0}}
		for i < 8 && j < 8 {
			zig = append(zig, src[x+i][y+j][k])
			if i == 0 || j == 0 || i == 7 || j == 7 {
				if !turned {
					k := 2 * (up*(j/7) | (1-up)*(i/7))
					i += corner[up][k]
					j += corner[up][k+1]
					turned = true
					up = 1 - up
					continue
				} else {
					turned = false
				}
			}
			i += d[up][0]
			j += d[up][1]
		}

		var dst []factor
		count := 0
		for i := 1; i < 64; i++ {
			if zig[i] != 0 {
				for count > 15 {
					dst = append(dst, factor{15, zig[i]})
					count -= 15
				}
				dst = append(dst, factor{count, zig[i]})
				count = 0
			} else {
				count++
			}
		}
		dst = append(dst, factor{0, 0})
		res[k] = dst
		dc[k] = zig[0]
	}
	return res, dc
}
