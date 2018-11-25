package makeGif

import (
	"fmt"
	"github.com/andybons/gogif"
	"image"
	"image/gif"
	"image/jpeg"
	"os"
)

func Make(src, dst string) {
	file1, err := os.OpenFile(src, os.O_RDONLY, 0)
	check(err)
	defer file1.Close()
	img, err := jpeg.Decode(file1)
	check(err)
	distGif, err := os.Create(dst)
	check(err)
	outGif := &gif.GIF{}
	bounds := img.Bounds()
	palettedImage := image.NewPaletted(bounds, nil)
	quantizer := gogif.MedianCutQuantizer{NumColor: 256}
	quantizer.Quantize(palettedImage, bounds, img, image.ZP)
	// Add new frame to animated GIF
	outGif.Image = append(outGif.Image, palettedImage)
	outGif.Delay = append(outGif.Delay, 0)
	err = gif.EncodeAll(distGif, outGif)
	if err != nil {
		fmt.Println(err)
	}
}


func check(err error) {
	if err != nil {
		panic(err)
	}
}