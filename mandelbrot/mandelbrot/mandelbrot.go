package mandelbrot

import (
	"runtime"
	"sync"

	"github.com/lucasb-eyer/go-colorful"
)

func CalcMandelbrot(screen []byte, screenHeight int, screenWidth int, minX float64, maxX float64, minI float64, maxI float64, colourScale []colorful.Color) {
	n_goroutines := runtime.NumCPU()
	rows_per_routine := screenHeight / n_goroutines

	var wg sync.WaitGroup

	for routine := range n_goroutines {
		wg.Go(func() {
			processRows(rows_per_routine*routine, rows_per_routine*(routine+1), screen, screenHeight, screenWidth, minX, maxX, minI, maxI, colourScale)
		})
	}
	wg.Wait()
}

func processRows(startRow int, endRow int, screen []byte, screenHeight int, screenWidth int, minX float64, maxX float64, minI float64, maxI float64, colourScale []colorful.Color) {
	iRange := maxI - minI
	xRange := maxX - minX
	maxIter := len(colourScale) - 1

	for screenI := startRow; screenI < endRow && screenI < screenHeight; screenI++ {
		i := (float64(screenI)/float64(screenHeight))*iRange + minI
		for screenX := 0; screenX < screenWidth; screenX++ {
			x := (float64(screenX)/float64(screenWidth))*xRange + minX
			iterCount := mandelbrotCount(x, i, maxIter)

			idx := (screenI*screenWidth + screenX) * 4
			pixelColour := colourScale[iterCount]

			screen[idx+0] = byte(pixelColour.R * 255) // R
			screen[idx+1] = byte(pixelColour.G * 255) // G
			screen[idx+2] = byte(pixelColour.B * 255) // B
			screen[idx+3] = 255                       // A

		}
	}
}

func mandelbrotCount(x float64, i float64, maxIter int) int {
	var zx, zi, tmpzx float64
	var iterCount int
	i2 := i * i
	q := (x-0.25)*(x-0.25) + i2
	if q*(q+(x-0.25)) < 0.25*i2 {
		return maxIter // In main cardioid
	}
	if (x+1.0)*(x+1.0)+i2 < 0.0625 {
		return maxIter // In period-2 bulb
	}

	tmpzx = zx
	zx = zx*zx - zi*zi + x
	zi = 2*tmpzx*zi + i
	for ; (zx*zx+zi*zi) <= 4 && iterCount < maxIter; iterCount++ {
		tmpzx = zx
		zx = zx*zx - zi*zi + x
		zi = 2*tmpzx*zi + i
	}
	return iterCount

}
