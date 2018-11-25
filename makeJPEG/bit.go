package makeJPEG

type bitArray struct {
	Data []byte
	Pos uint
}

func (b *bitArray) addData(num int) {
	data := toBinByte(num)
	for _, d := range data {
		b.addBit(d)
	}
}

func toBinByte(num int) (res []byte) {
	for num > 0 {
		i := num % 2
		res = append(res, byte(i))
		num /= 2
	}
	return res
}

func (b *bitArray) addBit(value byte) {
	if value != 0 && value != 1 {
		panic("Error bit value")
	}
	if  b.Pos == 0 {
		b.Data = append(b.Data, 0)
	}
	b.Data[len(b.Data) - 1] |= value << b.Pos
	b.Pos = b.Pos + 1
	if b.Pos == 8 {
		b.Pos = 0
	}
}
