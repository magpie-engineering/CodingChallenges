package colourscale

import (
	"github.com/lucasb-eyer/go-colorful"
)

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}

func MakeColourScale(scaleLen int) []colorful.Color {
	// The "keypoints" of the gradient.
	keypoints := GradientTable{
		{colorful.Color{R: 0.0, G: 7.0 / 255.0, B: 100.0 / 255.0}, 0.0},
		{colorful.Color{R: 32.0 / 255.0, G: 107.0 / 255.0, B: 203.0 / 255.0}, 0.16},
		{colorful.Color{R: 234.0 / 255.0, G: 1.1, B: 1.0}, 0.42},
		{colorful.Color{R: 1.0, G: 170.0 / 255.0, B: 0.0}, 0.64},
		{colorful.Color{R: 196.0 / 255.0, G: 30.0 / 255.0, B: 58.0 / 255.0}, 0.85},
		{colorful.Color{R: 0.0, G: 7.0 / 255.0, B: 100.0 / 255.0}, 1.0},
	}

	colourScale := make([]colorful.Color, scaleLen)
	for idx := range colourScale {
		colourScale[idx] = keypoints.GetInterpolatedColorFor(float64(idx) / float64(scaleLen))
	}

	return colourScale

}
