package mandelbrot

import (
	"runtime"
	"sync"

	"github.com/lucasb-eyer/go-colorful"
)

func CalcMandelbrot(screen []byte, screenHeight int, screenWidth int, minX float64, maxX float64, minI float64, maxI float64, colourScale []colorful.Color) {
	n_goroutines := runtime.NumCPU()
	rows_per_routine := screenHeight / n_goroutines
	iRange := maxI - minI
	xRange := maxX - minX
	maxIter := len(colourScale) - 1

	hInv := iRange / float64(screenHeight)
	wInv := xRange / float64(screenWidth)

	var wg sync.WaitGroup

	for routine := range n_goroutines {
		wg.Go(func() {
			for screenI := rows_per_routine * routine; screenI < rows_per_routine*(routine+1) && screenI < screenHeight; screenI++ {
				i := float64(screenI)*hInv + minI
				for screenX := 0; screenX < screenWidth; screenX++ {
					x := float64(screenX)*wInv + minX
					iterCount := mandelbrotCount(x, i, maxIter)

					idx := (screenI*screenWidth + screenX) * 4
					pixelColour := colourScale[iterCount]

					screen[idx+0] = byte(pixelColour.R * 255) // R
					screen[idx+1] = byte(pixelColour.G * 255) // G
					screen[idx+2] = byte(pixelColour.B * 255) // B
					screen[idx+3] = 255                       // A

				}
			}
		})
	}
	wg.Wait()
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
