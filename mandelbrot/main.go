package main

import (
	"log"

	"github.com/magpie-engineering/CodingChallenges/mandelbrot/mandelbrot"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	pixels  []byte
	img     *ebiten.Image
	width   int
	height  int
	centreX float64
	centreI float64
	zoom    float64
}

func NewGame(w, h int) *Game {
	return &Game{
		pixels: make([]byte, w*h*4), // RGBA, 4 bytes per pixel
		img:    ebiten.NewImage(w, h),
		width:  w,
		height: h,
		zoom:   1,
	}
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.centreX -= 0.25 / g.zoom
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.centreX += 0.25 / g.zoom
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.centreI -= 0.25 / g.zoom
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.centreI += 0.25 / g.zoom
	}
	if _, wheelY := ebiten.Wheel(); wheelY != 0 {
		if wheelY > 0 {
			g.zoom *= 1.1
		} else {
			g.zoom *= 0.9
		}

	}

	minX := -2/g.zoom + g.centreX
	maxX := 2/g.zoom + g.centreX
	minI := -2/g.zoom + g.centreI
	maxI := 2/g.zoom + g.centreI
	mandelbrot.CalcMandelbrot(g.pixels, g.height, g.width, minX, maxX, minI, maxI, 255)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.img.WritePixels(g.pixels)
	screen.DrawImage(g.img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func main() {
	g := NewGame(640*2, 480*2)
	ebiten.SetWindowSize(640*2, 480*2)
	ebiten.SetWindowTitle("Mandelbrot")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
