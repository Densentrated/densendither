package process

import (
	"image/color"
	"math"
)

func lanczos(x float64, a float64) float64 {
	if x == 0 {
		return 1.0
	}
	if math.Abs(x) >= a {
		return 0.0
	}

	pix := math.Pi * x
	pixDivA := pix / a

	return (math.Sin(pix) / pix) * (math.Sin(pixDivA) / pixDivA)
}

func clamp(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v + 0.5) // Round to nearest
}

// resizes an image to the specified newWidth, newHeight
func ResizeLanczos3(src [][]color.RGBA, newWidth, newHeight int) [][]color.RGBA {
	srcHeight := len(src)
	if srcHeight == 0 {
		return nil
	}
	srcWidth := len(src[0])
	if srcWidth == 0 {
		return nil
	}

	const a = 3.0 // Lanczos3

	// Calculate scaling factors
	xRatio := float64(srcWidth) / float64(newWidth)
	yRatio := float64(srcHeight) / float64(newHeight)

	// Create destination matrix
	dst := make([][]color.RGBA, newHeight)
	for i := range dst {
		dst[i] = make([]color.RGBA, newWidth)
	}

	// Create temporary buffer for horizontal pass
	temp := make([][]color.RGBA, srcHeight)
	for i := range temp {
		temp[i] = make([]color.RGBA, newWidth)
	}

	// Horizontal pass (resize width)
	for y := 0; y < srcHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculate source position
			srcX := (float64(x)+0.5)*xRatio - 0.5

			// Calculate contributing pixel range
			xMin := int(math.Floor(srcX - a + 1))
			xMax := int(math.Floor(srcX + a))

			var r, g, b, alpha, weightSum float64

			// Accumulate contributions
			for sx := xMin; sx <= xMax; sx++ {
				// Handle boundaries
				clampedX := sx
				if clampedX < 0 {
					clampedX = 0
				} else if clampedX >= srcWidth {
					clampedX = srcWidth - 1
				}

				// Calculate weight
				weight := lanczos(srcX-float64(sx), a)

				// Accumulate weighted pixel values
				pixel := src[y][clampedX]
				r += float64(pixel.R) * weight
				g += float64(pixel.G) * weight
				b += float64(pixel.B) * weight
				alpha += float64(pixel.A) * weight
				weightSum += weight
			}

			// Normalize and store
			if weightSum != 0 {
				temp[y][x] = color.RGBA{
					R: clamp(r / weightSum),
					G: clamp(g / weightSum),
					B: clamp(b / weightSum),
					A: clamp(alpha / weightSum),
				}
			}
		}
	}

	// Vertical pass (resize height)
	for x := 0; x < newWidth; x++ {
		for y := 0; y < newHeight; y++ {
			// Calculate source position
			srcY := (float64(y)+0.5)*yRatio - 0.5

			// Calculate contributing pixel range
			yMin := int(math.Floor(srcY - a + 1))
			yMax := int(math.Floor(srcY + a))

			var r, g, b, alpha, weightSum float64

			// Accumulate contributions
			for sy := yMin; sy <= yMax; sy++ {
				// Handle boundaries
				clampedY := sy
				if clampedY < 0 {
					clampedY = 0
				} else if clampedY >= srcHeight {
					clampedY = srcHeight - 1
				}

				// Calculate weight
				weight := lanczos(srcY-float64(sy), a)

				// Accumulate weighted pixel values
				pixel := temp[clampedY][x]
				r += float64(pixel.R) * weight
				g += float64(pixel.G) * weight
				b += float64(pixel.B) * weight
				alpha += float64(pixel.A) * weight
				weightSum += weight
			}

			// Normalize and store
			if weightSum != 0 {
				dst[y][x] = color.RGBA{
					R: clamp(r / weightSum),
					G: clamp(g / weightSum),
					B: clamp(b / weightSum),
					A: clamp(alpha / weightSum),
				}
			}
		}
	}
	return dst
}
