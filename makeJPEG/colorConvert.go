package makeJPEG

import "image"

// RGB to YIQ and reSampling
func convertToYIQ(src image.Image) [][][3]int {
	bounds := src.Bounds().Max
	mat := make([][][3]int, bounds.Y+(8-bounds.Y%8))
	for x := 0; x < bounds.Y; x++ {
		row := make([][3]int, bounds.X+(8-bounds.X%8))
		for y := 0; y < bounds.X; y++ {
			r, g, b, _ := src.At(y, x).RGBA()
			R, G, B := float32(r/257), float32(g/257), float32(b/257)
			// YUV
			row[y][0] = int(uint8(0.299*R + 0.587*G + 0.114*B))
			// 4:2:0 二次采样
			if y%2 == 0 {
				row[y][1] = int(uint8(-0.1687*R - 0.3313*G + 0.5*B + 128))
				row[y][2] = int(uint8(0.5*R - 0.4187*G - 0.0813*B + 128))
				if x%2 == 1 {
					row[y][1] = (int(mat[x-1][y][1]) + int(row[y][1])) / 2
					row[y][2] = (int(mat[x-1][y][2]) + int(row[y][2])) / 2
					row[y+1][1] = row[y][1]
					row[y+1][2] = row[y][2]
					mat[x-1][y][1] = row[y][1]
					mat[x-1][y][2] = row[y][2]
					mat[x-1][y+1][1] = row[y][1]
					mat[x-1][y+1][2] = row[y][2]
				}
			}
		}
		mat[x] = row
	}
	for i := bounds.Y; i < len(mat); i++ {
		row := make([][3]int, bounds.X+(8-bounds.X%8))
		mat[i] = row
	}
	return mat
}

func convertRGB(mat [][][3]int) [][][3]int {
	weight := len(mat)
	height := len(mat[0])
	for x := 0; x < weight; x++ {
		for y := 0; y < height; y++ {
			Y, U, V := mat[x][y][0], mat[x][y][1],mat[x][y][2]
			mat[x][y][0] = Y + int(1.402 * float32(V- 128))
			mat[x][y][1] = Y - int(0.34414 * float32(U- 128)) - int(0.71414 * float32(V- 128))
			mat[x][y][2] = Y + int(1.772 * float32(U- 128))
			for i := 0; i < 3; i++ {
				if mat[x][y][i] < 0 {
					mat[x][y][i] = 0
				} else if mat[x][y][i] > 255 {
					mat[x][y][i] = 255
				}
			}
		}
	}
	return mat
}