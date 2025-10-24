package main

import (
	"densendither/palette"
	"densendither/process"
)

func main() {

	var pico8Colors = [16]string{
		"#000000",
		"#1D2B53",
		"#7E2553",
		"#008751",
		"#AB5236",
		"#5F574F",
		"#C2C3C7",
		"#FFF1E8",
		"#FF004D",
		"#FFA300",
		"#FFEC27",
		"#00E436",
		"#29ADFF",
		"#83769C",
		"#FF77A8",
		"#FFCCAA",
	}

	purple_pallete := palette.Palette{
		Name:   "12526",
		Colors: pico8Colors[:],
	}

	city_pic, _ := process.LoadImage("city3.png")
	city_pic_matrix := process.ImageToMatrix(city_pic)
	city_pic_dithered := process.OrderedDither(city_pic_matrix, purple_pallete)
	process.SaveImageToFile(process.MatrixToImage(city_pic_dithered), "city3dithered.png")

}
