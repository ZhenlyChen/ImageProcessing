package EightBit

import (
	"fmt"
	"gocv.io/x/gocv"
	"sort"
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

type SortByR struct{ ColorSlice }
type SortByG struct{ ColorSlice }
type SortByB struct{ ColorSlice }

func (c SortByR) Less(i, j int) bool {
	return c.ColorSlice[i].R < c.ColorSlice[j].R
}
func (c SortByG) Less(i, j int) bool {
	return c.ColorSlice[i].G < c.ColorSlice[j].G
}
func (c SortByB) Less(i, j int) bool {
	return c.ColorSlice[i].B < c.ColorSlice[j].B
}

var colorTable ColorSlice

func DivRGB(data ColorSlice, deep int) {
	colorType := deep % 3
	half := len(data) / 2
	if colorType == 0 { // R
		sort.Sort(SortByR{data})
	} else if colorType == 1 { // G
		sort.Sort(SortByG{data})
	} else { // B
		sort.Sort(SortByB{data})
	}
	if deep >= 7 {
		var sumR, sumG, sumB int
		for _, c := range data[:half] {
			sumR += int(c.R)
			sumG += int(c.G)
			sumB += int(c.B)
		}
		colorTable = append(colorTable, RGBColor{
			R: uint8(sumR / half),
			G: uint8(sumG / half),
			B: uint8(sumB / half),
		})
		sumR, sumG, sumB = 0, 0, 0
		for _, c := range data[half:] {
			sumR += int(c.R)
			sumG += int(c.G)
			sumB += int(c.B)
		}
		colorTable = append(colorTable, RGBColor{
			R: uint8(sumR / half),
			G: uint8(sumG / half),
			B: uint8(sumB / half),
		})
	} else {
		DivRGB(data[:half], deep+1)
		DivRGB(data[half:], deep+1)
	}
}

func ToRGBColor(src gocv.Mat) (res ColorSlice) {
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

func To8Bit(src gocv.Mat) (res gocv.Mat) {
	res = src.Clone()
	DivRGB(ToRGBColor(src), 0)
	fmt.Println(colorTable)
	size := src.Size()
	for i := 0; i < size[0]; i++ {
		for j := 0; j < size[1]; j++ {
			oldColor := RGBColor{
				R: src.GetUCharAt(i, j*3),
				G: src.GetUCharAt(i, j*3+1),
				B: src.GetUCharAt(i, j*3+2),
			}
			newColor := getColor(oldColor)
			// fmt.Println(oldColor, newColor)
			res.SetUCharAt(i, j*3, newColor.R)
			res.SetUCharAt(i, j*3+1, newColor.G)
			res.SetUCharAt(i, j*3+2, newColor.B)
		}
		//panic("")
	}
	return
}

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
