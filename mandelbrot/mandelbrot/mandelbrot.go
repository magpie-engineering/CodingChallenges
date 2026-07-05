package mandelbrot

func CalcMandelbrot(screen []byte, screenHeight int, screenWidth int, minX float64, maxX float64, minI float64, maxI float64, maxIter byte) {
	iRange := maxI - minI
	xRange := maxX - minX

	for screenI := 0; screenI < screenHeight; screenI++ {
		i := (float64(screenI)/float64(screenHeight))*iRange + minI
		for screenX := 0; screenX < screenWidth; screenX++ {
			x := (float64(screenX)/float64(screenWidth))*xRange + minX
			iterCount := mandelbrotCount(x, i, maxIter)

			idx := (screenI*screenWidth + screenX) * 4

			screen[idx+0] = 255 - iterCount // R
			screen[idx+1] = 255 - iterCount // G
			screen[idx+2] = 255 - iterCount // B
			screen[idx+3] = 255             // A

		}
	}

}

func mandelbrotCount(x float64, i float64, maxIter uint8) byte {
	var zx, zi, tmpzx float64
	var iterCount byte
	if x == 0 && i == 0 {
		// will never escape
		return maxIter
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
