package main

import (
	"image"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/StephaneBunel/bresenham"
)

const iso = 1

type metaballsWidget struct {
	widget.BaseWidget

	model  *ensemble
	raster *canvas.Raster
}

func newMetaballsWidget(m *ensemble) *metaballsWidget {
	mw := &metaballsWidget{model: m}
	mw.raster = canvas.NewRaster(mw.draw)
	mw.ExtendBaseWidget(mw)
	return mw
}

func (mw *metaballsWidget) animate() {
	go func() {
		for range time.Tick(time.Millisecond * 50) {
			mw.model.move()
			mw.Refresh()
		}
	}()
}

func (mw *metaballsWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mw.raster)
}

func (mw *metaballsWidget) draw(w, h int) image.Image {
	fgcolor := theme.ForegroundColor()
	g := grid(w, h)
	gx := float32(g) / float32(w)
	gy := float32(g) / float32(h)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for row := 0; row < h; row += g {
		y := float32(row) / float32(h)
		for col := 0; col < w; col += g {
			x := float32(col) / float32(w)
			m := mw.model
			a := m.value(x, y)
			b := m.value(x+gx, y)
			c := m.value(x+gx, y+gy)
			d := m.value(x, y+gy)

			a1, a2 := lerp(col, col+g, (iso-a)/(b-a)), row
			b1, b2 := col+g, lerp(row, row+g, (iso-b)/(c-b))
			c1, c2 := lerp(col, col+g, (iso-d)/(c-d)), row+g
			d1, d2 := col, lerp(row, row+g, (iso-a)/(d-a))

			switch state(a, b, c, d) {
			case 1, 14:
				bresenham.DrawLine(img, c1, c2, d1, d2, fgcolor)
			case 2, 13:
				bresenham.DrawLine(img, b1, b2, c1, c2, fgcolor)
			case 3, 12:
				bresenham.DrawLine(img, b1, b2, d1, d2, fgcolor)
			case 4:
				bresenham.DrawLine(img, a1, a2, b1, b2, fgcolor)
			case 5:
				bresenham.DrawLine(img, a1, a2, d1, d2, fgcolor)
				bresenham.DrawLine(img, b1, b2, c1, c2, fgcolor)
			case 6:
				bresenham.DrawLine(img, a1, a2, c1, c2, fgcolor)
			case 7, 8:
				bresenham.DrawLine(img, a1, a2, d1, d2, fgcolor)
			case 9:
				bresenham.DrawLine(img, a1, a2, c1, c2, fgcolor)
			case 10:
				bresenham.DrawLine(img, a1, a2, b1, b2, fgcolor)
				bresenham.DrawLine(img, c1, c2, d1, d2, fgcolor)
			case 11:
				bresenham.DrawLine(img, a1, a2, b1, b2, fgcolor)
			}
		}
	}
	return img
}

func grid(w, h int) int {
	return int(math.Ceil(math.Sqrt(float64(w*h) / 16384)))
}

func state(tl, tr, br, bl float32) int {
	res := 0
	if tl >= iso {
		res |= 8
	}
	if tr >= iso {
		res |= 4
	}
	if br >= iso {
		res |= 2
	}
	if bl >= iso {
		res |= 1
	}
	return res
}

func lerp(a, b int, t float32) int {
	if t < 0 {
		return a
	}
	if t > 1 {
		return b
	}
	return int(math.Round(float64(a) + float64(b-a)*float64(t)))
}
