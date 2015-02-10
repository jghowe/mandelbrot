package mandelbrot

import (
	"image"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

func Render(dimensions *image.Point, centerX float64, centerY float64, zoom float64, colors []string, colorSpacing int) *image.NRGBA {
	var (
		px int
		minCx, minCy, widthCx, heightCy float64 = calculateExtents(centerX, centerY, zoom)
		maxIteration float64 = 4096;
		palette	[]colorful.Color
	)

	m := image.NewNRGBA(image.Rect(0,0,dimensions.X,dimensions.Y))

	insideColor, _ := colorful.Hex(colors[len(colors)-1])

	// extend palette with color spacing via blending
	palette = make([]colorful.Color, (len(colors)-1) * colorSpacing)

	for c := 0; c<len(colors)-1; c++ {
		c2 := int(math.Mod(float64(c+1), float64(len(colors)-1)))

		color, _ := colorful.Hex(colors[c])
		color2, _ := colorful.Hex(colors[c2])

		for i := 0; i<colorSpacing; i++ {
			palette[c*colorSpacing+i] = color.BlendRgb(color2, float64(i)/float64(colorSpacing))
		}
	}

	ch := make(chan int)

	// compute the mandlebot set
	for px = 0; px < dimensions.X; px++ {
		// use concurrency to take advantage of multiple cores
		go func(px int) {
			x0 := (float64(px) / float64(dimensions.X)) * widthCx + minCx

			for py := 0; py < dimensions.Y; py++ {

				y0 := (float64(py) / float64(dimensions.Y)) * heightCy + minCy

				x := 0.0
				y := 0.0

				xSquare := 0.0
				ySquare := 0.0

				iteration := 0.0

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

				index := py*m.Stride + px*4

				m.Pix[index] = uint8(color.R * 255)
				m.Pix[index+1] = uint8(color.G * 255)
				m.Pix[index+2] = uint8(color.B * 255)
				m.Pix[index+3] = 255
			}

			ch <- 1
		}(px)
	}

	for px = 0; px < dimensions.X; px++ {
		_ = <- ch
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
