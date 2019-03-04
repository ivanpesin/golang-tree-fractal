package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"math"
	"os"

	"github.com/fogleman/gg"
)

var config struct {
	N             int
	height, width int
	delay         int
	bangle, tilt  int
	gif           bool
	outfile       string
}

var myPalette = []color.Color{
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0xff, 0x00, 0xff},
	color.RGBA{0xff, 0xff, 0xff, 0xff},
}

var pic *gif.GIF
var canvas *gg.Context

func init() {
	flag.IntVar(&config.N, "depth", 5, "How deep the branches go")
	flag.IntVar(&config.height, "height", 200, "Image height")
	flag.IntVar(&config.width, "width", 300, "Image width")
	flag.IntVar(&config.bangle, "angle", 45, "Angle between branches")
	flag.IntVar(&config.tilt, "tilt", 0, "Tilt is an additional angle applied to all branches")
	flag.StringVar(&config.outfile, "out", "out.gif", "Output filename")
	flag.IntVar(&config.delay, "delay", 10, "Animation delay")
	flag.BoolVar(&config.gif, "gif", false, "Generate animated gif")

	flag.Parse()
}

func main() {

	canvas = gg.NewContext(config.width, config.height)
	canvas.SetRGB(1, 1, 1)
	canvas.Clear()

	if config.gif {
		pic = &gif.GIF{LoopCount: config.N}
	}

	fmt.Print("Generating fractal tree: ")
	drawBranch(config.N,
		float64(config.width)/2, float64(config.height),
		float64(config.height)*0.4, 90)

	canvas.SavePNG("out.png")
	fmt.Print(" Done\n")

	if config.gif {
		f, err := os.Create(config.outfile)
		if err != nil {
			log.Fatalf("Failed to create file: %s", err.Error())
		}
		defer f.Close()
		gif.EncodeAll(f, pic)
	}
}

func drawBranch(n int, x, y, l, a float64) {

	if n == 0 {
		return
	}

	//	fmt.Printf("%2d: Drawing at [%2.0f,%2.0f], l = %2.0f, a = %f\n", n, x, y, l, a)
	// fmt.Print(".")

	rad := a * math.Pi / 180
	cx := x - (l * math.Cos(rad))
	cy := y - (l * math.Sin(rad))

	//fmt.Printf("  : Drawing to [%2.0f,%2.0f]\n", cx, cy)

	if n != 1 {
		canvas.SetRGB(0, 0, 0)
	} else {
		canvas.SetRGB(0, 1, 0)
	}
	canvas.DrawLine(x, y, cx, cy)
	canvas.Stroke()

	if config.gif {
		pimg := image.NewPaletted(canvas.Image().Bounds(), myPalette)
		draw.Draw(pimg, pimg.Rect, canvas.Image(), pimg.Bounds().Min, draw.Over)

		pic.Delay = append(pic.Delay, config.delay)
		pic.Image = append(pic.Image, pimg)
	}

	drawBranch(n-1, cx, cy, l*0.6, float64(config.tilt)+a-float64(config.bangle))
	drawBranch(n-1, cx, cy, l*0.6, float64(config.tilt)+a+float64(config.bangle))
}
