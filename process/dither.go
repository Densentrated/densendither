package process

import (
	"densendither/palette"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

var bayer8x8 = [8][8]int{
	{0, 32, 8, 40, 2, 34, 10, 42},
	{48, 16, 56, 24, 50, 18, 58, 26},
	{12, 44, 4, 36, 14, 46, 6, 38},
	{60, 28, 52, 20, 62, 30, 54, 22},
	{3, 35, 11, 43, 1, 33, 9, 41},
	{51, 19, 59, 27, 49, 17, 57, 25},
	{15, 47, 7, 39, 13, 45, 5, 37},
	{63, 31, 55, 23, 61, 29, 53, 21},
}

func HexToRBGA(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 && len(hex) != 8 {
		return color.RGBA{}, fmt.Errorf("invalid hex color length: %d", len(hex))
	}

	values, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return color.RGBA{}, err
	}

	if len(hex) == 6 {
		return color.RGBA{
			R: uint8(values >> 16),
			G: uint8(values >> 8),
			B: uint8(values),
			A: 244,
		}, nil
	}

	return color.RGBA{
		R: uint8(values >> 24),
		G: uint8(values >> 16),
		B: uint8(values >> 8),
		A: uint8(values),
	}, nil
}

func findClosestColor(hex string, colorPalette palette.Palette) color.RGBA {
	minDist := math.MaxFloat64
	var closest color.RGBA
	currentColor, _ := HexToRBGA(hex)

	for _, p := range colorPalette.Colors {
		currentPaletteColor, _ := HexToRBGA(p)
		dr := float64(currentColor.R) - float64(currentPaletteColor.R)
		dg := float64(currentColor.G) - float64(currentPaletteColor.G)
		db := float64(currentColor.B) - float64(currentPaletteColor.B)
		dist := dr*dr + dg*dg + db*db

		if dist < minDist {
			minDist = dist
			closest = currentPaletteColor
		}
	}
	return closest
}

func getBayerThreshold(x int, y int) float64 {
	value := float64(bayer8x8[y%8][x%8]+1) / 65.0
	return value
}

func OrderedDither(image [][]color.RGBA, colorPalette palette.Palette) [][]color.RGBA {
	height := len(image)
	width := len(image[0])

	result := make([][]color.RGBA, height)
	for i := range result {
		result[i] = make([]color.RGBA, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := image[y][x]
			threshold := getBayerThreshold(x, y)

			r := clamp(float64(oldPixel.R) + (threshold-.5)*64)
			g := clamp(float64(oldPixel.G) + (threshold-.5)*64)
			b := clamp(float64(oldPixel.B) + (threshold-.5)*64)

			hexCode := fmt.Sprintf("#%02x%02x%02x", int(r), int(g), int(b))

			result[y][x] = findClosestColor(hexCode, colorPalette)
		}
	}

	return result
}

func FloydSteinbergDither(image [][]color.RGBA, colorPalette palette.Palette) [][]color.RGBA {
	height := len(image)
	if height == 0 {
		return image
	}
	width := len(image[0])
	if width == 0 {
		return image
	}

	work := make([][]struct {
		R, G, B float64
		A       uint8
	}, height)

	for y := 0; y < height; y++ {
		work[y] = make([]struct {
			R, G, B float64
			A       uint8
		}, width)
		for x := 0; x < width; x++ {
			work[y][x].R = float64(image[y][x].R)
			work[y][x].G = float64(image[y][x].G)
			work[y][x].B = float64(image[y][x].B)
			work[y][x].A = image[y][x].A
		}
	}
	result := make([][]color.RGBA, height)
	for i := range result {
		result[i] = make([]color.RGBA, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldR := clamp(work[y][x].R)
			oldG := clamp(work[y][x].G)
			oldB := clamp(work[y][x].B)
			oldA := work[y][x].A

			hexCode := fmt.Sprintf("#%02x%02x%02x", oldR, oldG, oldB)

			newPixel := findClosestColor(hexCode, colorPalette)
			// Preserve original alpha channel
			newPixel.A = oldA
			result[y][x] = newPixel

			errR := work[y][x].R - float64(newPixel.R)
			errG := work[y][x].G - float64(newPixel.G)
			errB := work[y][x].B - float64(newPixel.B)

			// Distribute error to neighboring pixels using Floyd-Steinberg matrix
			// Right pixel (x+1, y) gets 7/16 of error
			if x+1 < width {
				work[y][x+1].R += errR * 7.0 / 16.0
				work[y][x+1].G += errG * 7.0 / 16.0
				work[y][x+1].B += errB * 7.0 / 16.0
			}

			// Bottom-left pixel (x-1, y+1) gets 3/16 of error
			if y+1 < height && x-1 >= 0 {
				work[y+1][x-1].R += errR * 3.0 / 16.0
				work[y+1][x-1].G += errG * 3.0 / 16.0
				work[y+1][x-1].B += errB * 3.0 / 16.0
			}

			// Bottom pixel (x, y+1) gets 5/16 of error
			if y+1 < height {
				work[y+1][x].R += errR * 5.0 / 16.0
				work[y+1][x].G += errG * 5.0 / 16.0
				work[y+1][x].B += errB * 5.0 / 16.0
			}

			// Bottom-right pixel (x+1, y+1) gets 1/16 of error
			if y+1 < height && x+1 < width {
				work[y+1][x+1].R += errR * 1.0 / 16.0
				work[y+1][x+1].G += errG * 1.0 / 16.0
				work[y+1][x+1].B += errB * 1.0 / 16.0
			}
		}
	}
	return result
}
