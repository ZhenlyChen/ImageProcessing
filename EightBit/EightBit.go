package EightBit

import (
	"sort"

	"gocv.io/x/gocv"
)

type RGBColor struct {
	R uint8
	G uint8
	B uint8
}

type ColorSlice []RGBColor

func (c ColorSlice) Len() int {
	return len(c)
}

func (c ColorSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type sortByR struct{ ColorSlice }
type sortByG struct{ ColorSlice }
type sortByB struct{ ColorSlice }

func (c sortByR) Less(i, j int) bool {
	return c.ColorSlice[i].R < c.ColorSlice[j].R
}
func (c sortByG) Less(i, j int) bool {
	return c.ColorSlice[i].G < c.ColorSlice[j].G
}
func (c sortByB) Less(i, j int) bool {
	return c.ColorSlice[i].B < c.ColorSlice[j].B
}

// colorTable LUT
var colorTable ColorSlice

// divRGB 划分颜色，生成LUT
func divRGB(data ColorSlice, deep int) {
	length := len(data)
	if deep >= 8 { // 已经划分了7次，再划分一次就可以生成256个区域
		var sumR, sumG, sumB int
		for _, c := range data[:] {
			sumR += int(c.R)
			sumG += int(c.G)
			sumB += int(c.B)
		}
		colorTable = append(colorTable, RGBColor{
			R: uint8(sumR / length),
			G: uint8(sumG / length),
			B: uint8(sumB / length),
		})
		return
	}
	switch deep % 3 {
	case 0:
		sort.Sort(sortByR{data})
	case 1:
		sort.Sort(sortByG{data})
	case 2:
		sort.Sort(sortByB{data})
	}
	// 继续划分
	divRGB(data[:length/2], deep+1)
	divRGB(data[length/2:], deep+1)
}

func toRGBColor(src gocv.Mat) (res ColorSlice) {
	size := src.Size()
	for i := 0; i < size[0]; i++ {
		for j := 0; j < size[1]; j++ {
			res = append(res, RGBColor{
				R: src.GetUCharAt(i, j*3),
				G: src.GetUCharAt(i, j*3+1),
				B: src.GetUCharAt(i, j*3+2),
			})
		}
	}
	return
}

// To8Bit 转成 8 位照片
func To8Bit(src gocv.Mat) (res gocv.Mat) {
	res = src.Clone()
	divRGB(toRGBColor(src), 0)
	size := src.Size()
	for i := 0; i < size[0]; i++ {
		for j := 0; j < size[1]; j++ {
			oldColor := RGBColor{
				R: src.GetUCharAt(i, j*3),
				G: src.GetUCharAt(i, j*3+1),
				B: src.GetUCharAt(i, j*3+2),
			}
			newColor := getColor(oldColor)
			res.SetUCharAt(i, j*3, newColor.R)
			res.SetUCharAt(i, j*3+1, newColor.G)
			res.SetUCharAt(i, j*3+2, newColor.B)
		}
	}
	return
}

// 寻找欧氏距离最短的颜色
func getColor(src RGBColor) RGBColor {
	index := 0
	dis := getDis(src, colorTable[0])
	for i, c := range colorTable {
		nDis := getDis(src, c)
		if nDis < dis {
			index = i
			dis = nDis
		}
	}
	return colorTable[index]
}

func getDis(a, c RGBColor) int {
	var r, g, b int
	r = int(a.R) - int(c.R)
	g = int(a.G) - int(c.G)
	b = int(a.B) - int(c.B)
	return r*r + g*g + b*b
}
