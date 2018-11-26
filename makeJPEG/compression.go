package makeJPEG


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
				row = append(row, mat[i][k]-last)
				last = mat[i][k]
			}
		}
		dst[k] = row
	}
	return dst
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
				dst = append(dst, factor{count, zig[i]})
				count = 0
			} else {
				count++
				for count > 15 {
					dst = append(dst, factor{15, 0})
					count -= 15
				}
			}
		}
		dst = append(dst, factor{0, 0})
		res[k] = dst
		dc[k] = zig[0]
	}
	return res, dc
}
