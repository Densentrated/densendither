package main

import (
	"densendither/palette"
	"densendither/process"
	"fmt"
)

func main() {

	var pico8Colors = [16]string{
		"#000000", // black
		"#1a0d26", // very dark purple
		"#2d1b3d", // dark purple
		"#4a2c5a", // medium dark purple
		"#663d73", // medium purple
		"#824e8c", // medium purple
		"#9f5fa5", // medium light purple
		"#bc70be", // light purple
		"#d981d7", // lighter purple
		"#f692f0", // very light purple
		"#e6ccff", // pale purple
		"#f0e6ff", // very pale purple
		"#f7f3ff", // almost white purple
		"#fcfaff", // near white
		"#fefeff", // off white
		"#ffffff", // white
	}

	purple_pallete := palette.Palette{
		Name:   "12526",
		Colors: pico8Colors[:],
	}

	city_pic, _ := process.LoadImage("city6.png")
	city_pic_matrix := process.ImageToMatrix(city_pic)
	city_pic_dithered := process.OrderedDither(city_pic_matrix, purple_pallete)
	process.SaveImageToFile(process.MatrixToImage(city_pic_dithered), "city6floydithered.png")
	fmt.Println("Image has been Dithered!")

}
