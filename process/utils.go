package process

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
)

// loads an image from a filepath with error propogation
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Loaded %s image\n", format)
	return img, nil
}

// converts an image.Image object to a matrix of color objects
func ImageToMatrix(img image.Image) [][]color.RGBA {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	matrix := make([][]color.RGBA, height)

	switch img := img.(type) {
	case *image.RGBA:
		// Direct access for RGBA images
		for y := 0; y < height; y++ {
			matrix[y] = make([]color.RGBA, width)
			for x := 0; x < width; x++ {
				idx := img.PixOffset(bounds.Min.X+x, bounds.Min.Y+y)
				matrix[y][x] = color.RGBA{
					R: img.Pix[idx],
					G: img.Pix[idx+1],
					B: img.Pix[idx+2],
					A: img.Pix[idx+3],
				}
			}
		}
	case *image.NRGBA:
		// Handle non-premultiplied RGBA
		for y := 0; y < height; y++ {
			matrix[y] = make([]color.RGBA, width)
			for x := 0; x < width; x++ {
				c := img.NRGBAAt(bounds.Min.X+x, bounds.Min.Y+y)
				matrix[y][x] = color.RGBA(c)
			}
		}
	default:
		// Fallback for other formats
		for y := 0; y < height; y++ {
			matrix[y] = make([]color.RGBA, width)
			for x := 0; x < width; x++ {
				c := color.RGBAModel.Convert(img.At(bounds.Min.X+x, bounds.Min.Y+y))
				matrix[y][x] = c.(color.RGBA)
			}
		}
	}
	return matrix
}

// converts a matrix of color.RGBA objects back to an image.Image
func MatrixToImage(matrix [][]color.RGBA) *image.RGBA {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return nil
	}

	height := len(matrix)
	width := len(matrix[0])

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, matrix[y][x])
		}
	}

	return img
}

// saveImageToFile saves an image.Image to a file with the specified filename
// The image format is determined by the file extension
func SaveImageToFile(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Always save as PNG
	return png.Encode(file, img)
}
