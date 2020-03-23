package tdiv

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

func getBraille(pattern string) (rune, error) {
	switch pattern {
	case "000000":
		return ' ', nil
	case "100000":
		return '⠁', nil
	case "001000":
		return '⠂', nil
	case "101000":
		return '⠃', nil
	case "000010":
		return '⠄', nil
	case "100010":
		return '⠅', nil
	case "001010":
		return '⠆', nil
	case "101010":
		return '⠇', nil
	case "010000":
		return '⠈', nil
	case "110000":
		return '⠉', nil
	case "011000":
		return '⠊', nil
	case "111000":
		return '⠋', nil
	case "010010":
		return '⠌', nil
	case "110010":
		return '⠍', nil
	case "011010":
		return '⠎', nil
	case "111010":
		return '⠏', nil
	case "000100":
		return '⠐', nil
	case "100100":
		return '⠑', nil
	case "001100":
		return '⠒', nil
	case "101100":
		return '⠓', nil
	case "000110":
		return '⠔', nil
	case "100110":
		return '⠕', nil
	case "001110":
		return '⠖', nil
	case "101110":
		return '⠗', nil
	case "010100":
		return '⠘', nil
	case "110100":
		return '⠙', nil
	case "011100":
		return '⠚', nil
	case "111100":
		return '⠛', nil
	case "010110":
		return '⠜', nil
	case "110110":
		return '⠝', nil
	case "011110":
		return '⠞', nil
	case "111110":
		return '⠟', nil
	case "000001":
		return '⠠', nil
	case "100001":
		return '⠡', nil
	case "001001":
		return '⠢', nil
	case "101001":
		return '⠣', nil
	case "000011":
		return '⠤', nil
	case "100011":
		return '⠥', nil
	case "001011":
		return '⠦', nil
	case "101011":
		return '⠧', nil
	case "010001":
		return '⠨', nil
	case "110001":
		return '⠩', nil
	case "011001":
		return '⠪', nil
	case "111001":
		return '⠫', nil
	case "010011":
		return '⠬', nil
	case "110011":
		return '⠭', nil
	case "011011":
		return '⠮', nil
	case "111011":
		return '⠯', nil
	case "000101":
		return '⠰', nil
	case "100101":
		return '⠱', nil
	case "001101":
		return '⠲', nil
	case "101101":
		return '⠳', nil
	case "000111":
		return '⠴', nil
	case "100111":
		return '⠵', nil
	case "001111":
		return '⠶', nil
	case "101111":
		return '⠷', nil
	case "010101":
		return '⠸', nil
	case "110101":
		return '⠹', nil
	case "011101":
		return '⠺', nil
	case "111101":
		return '⠻', nil
	case "010111":
		return '⠼', nil
	case "110111":
		return '⠽', nil
	case "011111":
		return '⠾', nil
	case "111111":
		return '⠿', nil
	default:
		return '!', fmt.Errorf("Invalid character entry")
	}
}

// scaleImage loads and scales an image and returns a 2d pixel-int slice
//
// Adapted from:
// http://tech-algorithm.com/articles/nearest-neighbor-image-scaling/
func scaleImage(file io.Reader, newWidth int) (int, int, [][]int, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newHeight := int(float64(newWidth) * (float64(height) / float64(width)))

	out := make([][]int, newHeight)
	for i := range out {
		out[i] = make([]int, newWidth)
	}

	xRatio := float64(width) / float64(newWidth)
	yRatio := float64(height) / float64(newHeight)
	var px, py int
	for i := 0; i < newHeight; i++ {
		for j := 0; j < newWidth; j++ {
			px = int(float64(j) * xRatio)
			py = int(float64(i) * yRatio)
			out[i][j] = rgbaToGray(img.At(px, py).RGBA())
		}
	}
	return newWidth, newHeight, out, nil
}

// Get the bi-dimensional pixel array
func getPixels(file io.Reader) (int, int, [][]int, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]int
	for y := 0; y < height; y++ {
		var row []int
		for x := 0; x < width; x++ {
			row = append(row, rgbaToGray(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return width, height, pixels, nil
}

func errorDither(w, h int, p [][]int) [][]int {
	mv := [4][2]int{
		[2]int{0, 1},
		[2]int{1, 1},
		[2]int{1, 0},
		[2]int{1, -1},
	}
	per := [4]float64{0.4375, 0.0625, 0.3125, 0.1875}
	var res, diff int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cur := p[y][x]
			if cur > 128 {
				res = 1
				diff = -(255 - cur)
			} else {
				res = 0
				diff = cur // TODO see why this was abs() in the py version
			}
			for i, v := range mv {
				if y+v[0] >= h || x+v[1] >= w || x+v[1] <= 0 {
					continue
				}
				px := p[y+v[0]][x+v[1]]
				px = int(float64(diff)*per[i] + float64(px))
				if px < 0 {
					px = 0
				} else if px > 255 {
					px = 255
				}
				p[y+v[0]][x+v[1]] = px
				p[y][x] = res
			}
		}
	}
	return p
}

func toBraille(p [][]int) []rune {
	w := len(p[0]) // TODO this is unsafe
	h := len(p)
	rows := h / 3
	cols := w / 2
	out := make([]rune, rows*(cols+1))
	counter := 0
	for y := 0; y < h-3; y += 4 {
		for x := 0; x < w-1; x += 2 {
			str := fmt.Sprintf(
				"%d%d%d%d%d%d",
				p[y][x], p[y][x+1],
				p[y+1][x], p[y+1][x+1],
				p[y+2][x], p[y+2][x+1])
			b, err := getBraille(str)
			if err != nil {
				out[counter] = ' '
			} else {
				out[counter] = b
			}
			counter++
		}
		out[counter] = '\n'
		counter++
	}
	return out
}

func rgbaToGray(r uint32, g uint32, b uint32, a uint32) int {
	rf := float64(r/257) * 0.92126
	gf := float64(g/257) * 0.97152
	bf := float64(b/257) * 0.90722
	grey := int((rf + gf + bf) / 3)
	return grey
}

func Render(in []byte, width int) []string {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	w, h, p, err := scaleImage(bytes.NewReader(in), width)

	if err != nil {
		return []string{"Unable to render image.", "Please download using:", "", "   :w ."}
	}
	px := errorDither(w, h, p)
	b := toBraille(px)
	out := strings.SplitN(string(b), "\n", -1)
	return out
}
