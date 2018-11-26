package makeJPEG

import "math"

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

// 反量化
func Iquantitative(src [][][3]int, x, y int) [][][3]int {

	// Y
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			src[x+i][y+j][0] = int(math.Round(float64(src[x+i][y+j][0]) * float64(matrixY[i][j])))
		}
	}
	// IQ
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			src[x+i][y+j][1] = int(math.Round(float64(src[x+i][y+j][1]) * float64(matrixIQ[i][j])))
			src[x+i][y+j][2] = int(math.Round(float64(src[x+i][y+j][2]) * float64(matrixIQ[i][j])))
		}
	}

	return src
}