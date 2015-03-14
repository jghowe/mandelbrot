package mandelbrot

import (
	"image"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// This function renders the Mandelbrot set using the complex coordinates (centerX, centerY) and zoom factor.  The colors parameter is
// a slice of color hex string that will be used to color the image.  The last color in the slice is used to color the inside
// of the fractal.  The colorSpacing parameter determines how often the color scheme repeats.  A low value will cause the
// colors to repeat often while a high value will cause the causes the colors to be stretched out.
func Render(dimensions *image.Point, centerX float64, centerY float64, zoom float64, colors []string, colorSpacing int) *image.NRGBA {
	var (
		minCx, minCy, widthCx, heightCy float64 = calculateExtents(centerX, centerY, zoom)
	)

	m := image.NewNRGBA(image.Rect(0,0,dimensions.X,dimensions.Y))

	insideColor, _ := colorful.Hex(colors[len(colors)-1])

	palette := makePalette(colors, colorSpacing)

	ch := make(chan pixelRow)

	// compute the mandlebot set
	for py := 0; py < dimensions.Y; py++ {
		y0 := (float64(py) / float64(dimensions.Y)) * heightCy + minCy

		// use concurrency to take advantage of multiple cores
		go renderRow(py, y0, dimensions.X, minCx, minCy, widthCx, heightCy, palette, insideColor, ch)
	}

	// receive rendered rows and write to image
	for py := 0; py < dimensions.Y; py++ {
		// receive result
		pixelRow := <- ch

		// write color values to image
		for px := 0; px < dimensions.X; px++ {
			index := pixelRow.Row*m.Stride + px*4
			color := pixelRow.Pixels[px]

			m.Pix[index] = uint8(color.R * 255)
			m.Pix[index+1] = uint8(color.G * 255)
			m.Pix[index+2] = uint8(color.B * 255)
			m.Pix[index+3] = 255
		}
	}

	return m
}

// struct for receiving rendered rows
type pixelRow struct {
	Row int
	Pixels []colorful.Color
}

// render one pixel row of the mandelbrot set
func renderRow(py int, y0 float64, width int, minCx float64, minCy float64, widthCx float64, heightCy float64, palette []colorful.Color, insideColor colorful.Color, ch chan pixelRow) {
	pixels := make([]colorful.Color, width)

	for px := 0; px < width; px++ {
		x0 := (float64(px) / float64(width)) * widthCx + minCx

		pixels[px]= renderPixel(px, py, x0, y0, palette, insideColor)
	}

	// send result back using channel
	ch <- pixelRow{py, pixels}
}

// render one pixel of the mandelbrot set with the given location and palette
func renderPixel(px int, py int, x0 float64, y0 float64, palette []colorful.Color, insideColor colorful.Color) (colorful.Color) {
	var x, y, xSquare, ySquare, iteration float64
	var maxIteration float64 = 4096

	for ; (xSquare + ySquare) < 4 && iteration < maxIteration; iteration++ {
		y = (x + y) * (x + y) - xSquare - ySquare + y0
		x = xSquare - ySquare + x0
		xSquare = x * x
		ySquare = y * y
	}

	if (iteration < maxIteration) {
		zn := xSquare + ySquare
		nu := math.Log( math.Log(zn) / 2 / math.Log(2) ) / math.Log(2)
		iteration = iteration + 1 - nu
	}

	color1 := palette[ int(math.Floor(math.Mod( iteration, float64(len(palette))))) ]
	color2 := palette[ int(math.Floor(math.Mod( iteration + 1, float64(len(palette))))) ]

	color := color1.BlendRgb(color2, math.Mod(iteration, 1))

	if (iteration >= maxIteration) {
		color = insideColor
	}

	return color
}

// make palette by extending colors via blending
func makePalette(colors []string, colorSpacing int) ([]colorful.Color) {
	length := len(colors)-1
	palette := make([]colorful.Color, (length * colorSpacing))

	for c := 0; c<length; c++ {
		c2 := (c+1) % length

		color, _ := colorful.Hex(colors[c])
		color2, _ := colorful.Hex(colors[c2])

		for i := 0; i<colorSpacing; i++ {
			palette[c*colorSpacing+i] = color.BlendRgb(color2, float64(i)/float64(colorSpacing))
		}
	}

	return palette
}

// calculate extents in complex space based on center and zoom factor
func calculateExtents(centerX float64, centerY float64, zoom float64) (float64, float64, float64, float64) {
	widthDefaultX := 3.0
	heightDefaultY := 2.0

	widthCx := widthDefaultX / zoom
	heightCy := heightDefaultY / zoom

	minCx := centerX - widthCx / 2
	minCy := centerY - heightCy /2

	return minCx, minCy, widthCx, heightCy
}
