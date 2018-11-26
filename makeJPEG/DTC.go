package makeJPEG

import (
	"math"
)

// DTC离散余弦变换
func convertDCT(src [][][3]int, x, y int) [][][3]int {
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


// IDTC逆离散余弦变换
func convertIDCT(src [][][3]int, x, y int) [][][3]int {
	cSrc := make([][][3]int, 8)
	for k := 0; k < 8; k++ {
		cSrc[k] = make([][3]int, 8)
		copy(cSrc[k], src[x+k][y:y+8])
	}
	for k := 0; k < 3; k++ {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				var res float64
				for u := 0; u < 8; u++ {
					for v := 0; v < 8; v++ {
						cU, cV := 1.0, 1.0
						if u == 0 {
							cU = math.Sqrt(2) / 2.0
						}
						if v == 0 {
							cV = math.Sqrt(2) / 2.0
						}
						value := (cU * cV) / 4.0
						value *= math.Cos(float64(2*i+1)*float64(u)*float64(math.Pi)/16.0) *
							math.Cos(float64(2*j+1)*float64(v)*float64(math.Pi)/16.0) *
							float64(cSrc[u][v][k])
						// fmt.Println(value)
						res += value
					}
				}
				// fmt.Println("res:", res)
				src[x+i][y+j][k] = int(math.Round(res))
			}
		}
	}
	return src
}