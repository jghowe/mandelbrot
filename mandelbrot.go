package mandelbrot

import (
	"image"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

func Render(dimensions *image.Point, centerX float64, centerY float64, zoom float64, colors []string, repeatStep int) *image.NRGBA {
	var (
		px, py int
		x, y, x0, y0, xSquare, ySquare float64
		minCx, minCy, widthCx, heightCy float64 = calculateExtents(centerX, centerY, zoom)
		iteration, maxIteration float64 = 0, 4096
		palette	[]colorful.Color
	)

	m := image.NewNRGBA(image.Rect(0,0,dimensions.X,dimensions.Y))

	// extend palette to repeat interval via blending
	palette = make([]colorful.Color, len(colors) * repeatStep)

	for c := 0; c<len(colors); c++ {
		c2 := int(math.Mod(float64(c+1), float64(len(colors))))

		color, _ := colorful.Hex(colors[c])
		color2, _ := colorful.Hex(colors[c2])

		for i := 0; i<repeatStep; i++ {
			palette[c*repeatStep+i] = color.BlendRgb(color2, float64(i)/float64(repeatStep))
		}
	}

	// compute the mandlebot set
	for px = 0; px < dimensions.X; px++ {
		x0 = (float64(px) / float64(dimensions.X)) * widthCx + minCx

		for py = 0; py < dimensions.Y; py++ {

			y0 = (float64(py) / float64(dimensions.Y)) * heightCy + minCy

			x = 0
			y = 0

			xSquare = 0
			ySquare = 0

			for iteration = 0; (xSquare + ySquare) < 4 && iteration < maxIteration; iteration++ {
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

			index := py*m.Stride + px*4

			m.Pix[index] = uint8(color.R * 255)
			m.Pix[index+1] = uint8(color.G * 255)
			m.Pix[index+2] = uint8(color.B * 255)
			m.Pix[index+3] = 255
		}
	}

	return m
}

func calculateExtents(centerX float64, centerY float64, zoom float64) (float64, float64, float64, float64) {

	widthDefaultX := 3.0
	heightDefaultY := 2.0

	widthCx := widthDefaultX / zoom
	heightCy := heightDefaultY / zoom

	minCx := centerX - widthCx / 2
	minCy := centerY - heightCy /2

	return minCx, minCy, widthCx, heightCy
}
