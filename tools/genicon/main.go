// tools/genicon generates the 22x22 template icon for the tray app.
// Run: go run ./tools/genicon
// Output: internal/adapters/tray/icon.png
package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const size = 22

func main() {
	img := image.NewNRGBA(image.Rect(0, 0, size, size))

	blend := func(px, py int, a float64) {
		if px < 0 || px >= size || py < 0 || py >= size {
			return
		}
		existing := img.NRGBAAt(px, py)
		na := uint8(math.Min(255, float64(existing.A)+a*255))
		img.SetNRGBA(px, py, color.NRGBA{0, 0, 0, na})
	}

	stroke := func(x, y, radius float64) {
		ri := int(math.Ceil(radius + 1))
		for dx := -ri; dx <= ri; dx++ {
			for dy := -ri; dy <= ri; dy++ {
				dist := math.Sqrt(float64(dx*dx + dy*dy))
				a := math.Min(1, math.Max(0, radius-dist+0.5))
				if a > 0 {
					blend(int(math.Round(x))+dx, int(math.Round(y))+dy, a)
				}
			}
		}
	}

	drawBezier := func(r float64, x0, y0, x1, y1, x2, y2, x3, y3 float64) {
		for i := 0; i <= 600; i++ {
			t := float64(i) / 600
			u := 1 - t
			x := u*u*u*x0 + 3*u*u*t*x1 + 3*u*t*t*x2 + t*t*t*x3
			y := u*u*u*y0 + 3*u*u*t*y1 + 3*u*t*t*y2 + t*t*t*y3
			stroke(x, y, r)
		}
	}

	drawLine := func(r, x0, y0, x1, y1 float64) {
		n := int(math.Sqrt((x1-x0)*(x1-x0)+(y1-y0)*(y1-y0)) * 10)
		if n < 1 {
			n = 1
		}
		for i := 0; i <= n; i++ {
			t := float64(i) / float64(n)
			stroke(x0+t*(x1-x0), y0+t*(y1-y0), r)
		}
	}

	// === Bow (left half, x=1..10) ===
	// C-shape: from (10, 2) down to (10, 20), bowing left to x≈2 at midpoint.
	drawBezier(1.0, 10, 2, -1, 2, -1, 20, 10, 20)

	// String: thin line along the flat side of the bow.
	drawLine(0.5, 10, 2, 10, 20)

	// === Arrow (right half, x=14..18) ===
	// Shaft: x=15, y=7 to y=21.
	drawLine(1.0, 15, 7, 15, 21)

	// Arrowhead: tip (15,1), left edge to (12,8), right edge to (18,8), base (12,8)-(18,8).
	drawLine(1.0, 15, 1, 12, 8)
	drawLine(1.0, 15, 1, 18, 8)
	drawLine(1.0, 12, 8, 18, 8)

	f, err := os.Create("internal/adapters/tray/icon.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
