package main

import (
	"log"

	"github.com/magpie-engineering/CodingChallenges/mandelbrot/colourscale"
	"github.com/magpie-engineering/CodingChallenges/mandelbrot/mandelbrot"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lucasb-eyer/go-colorful"
)

type Game struct {
	pixels      []byte
	img         *ebiten.Image
	width       int
	height      int
	centreX     float64
	centreI     float64
	zoom        float64
	colourScale []colorful.Color
}

func NewGame(w, h, maxIter int) *Game {

	return &Game{
		pixels:      make([]byte, w*h*4), // RGBA, 4 bytes per pixel
		img:         ebiten.NewImage(w, h),
		width:       w,
		height:      h,
		zoom:        1,
		colourScale: colourscale.MakeColourScale(maxIter),
	}
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ebiten.SetFullscreen(false)
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.centreX -= 0.05 / g.zoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.centreX += 0.05 / g.zoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.centreI -= 0.05 / g.zoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.centreI += 0.05 / g.zoom
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.zoom *= 1.05
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.zoom *= 0.95
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
	mandelbrot.CalcMandelbrot(g.pixels, g.height, g.width, minX, maxX, minI, maxI, g.colourScale)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.img.WritePixels(g.pixels)
	screen.DrawImage(g.img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.width = outsideWidth
	g.height = outsideHeight
	g.pixels = make([]byte, g.width*g.height*4) // RGBA, 4 bytes per pixel
	g.img = ebiten.NewImage(g.width, g.height)
	return g.width, g.height
}

func main() {
	g := NewGame(640*2, 480*2, 255)
	ebiten.SetWindowSize(640*2, 480*2)
	ebiten.SetWindowTitle("Mandelbrot")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
