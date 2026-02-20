package annotate

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

func ApplyOps(dst draw.Image, ops []core.AnnotationOp) error {
	for _, op := range ops {
		switch op.Kind {
		case "rect":
			var p RectPayload
			if err := json.Unmarshal(op.Payload, &p); err != nil {
				return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
			}
			renderRect(dst, p)
		case "line":
			var p LinePayload
			if err := json.Unmarshal(op.Payload, &p); err != nil {
				return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
			}
			renderLine(dst, p)
		case "arrow":
			var p ArrowPayload
			if err := json.Unmarshal(op.Payload, &p); err != nil {
				return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
			}
			renderArrow(dst, p)
		case "text":
			var p TextPayload
			if err := json.Unmarshal(op.Payload, &p); err != nil {
				return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
			}
			renderText(dst, p)
		case "blur":
			var p BlurPayload
			if err := json.Unmarshal(op.Payload, &p); err != nil {
				return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
			}
			applyBlur(dst, p)
		case "pixelate":
			var p PixelatePayload
			if err := json.Unmarshal(op.Payload, &p); err != nil {
				return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
			}
			applyPixelate(dst, p)
		}
	}
	return nil
}

func parseColor(hex string) color.RGBA {
	if len(hex) != 7 || hex[0] != '#' {
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}
	r, err := strconv.ParseUint(hex[1:3], 16, 8)
	if err != nil {
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}
	g, err := strconv.ParseUint(hex[3:5], 16, 8)
	if err != nil {
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}
	b, err := strconv.ParseUint(hex[5:7], 16, 8)
	if err != nil {
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

func renderRect(dst draw.Image, p RectPayload) {
	c := parseColor(p.Color)
	r := image.Rect(p.X, p.Y, p.X+p.W, p.Y+p.H)
	if p.Fill {
		draw.Draw(dst, r, image.NewUniform(c), image.Point{}, draw.Over)
		return
	}
	s := p.StrokeWidth
	if s <= 0 {
		s = 2
	}
	for i := 0; i < s; i++ {
		drawLine(dst, p.X+i, p.Y+i, p.X+p.W-i, p.Y+i, c)
		drawLine(dst, p.X+i, p.Y+p.H-i, p.X+p.W-i, p.Y+p.H-i, c)
		drawLine(dst, p.X+i, p.Y+i, p.X+i, p.Y+p.H-i, c)
		drawLine(dst, p.X+p.W-i, p.Y+i, p.X+p.W-i, p.Y+p.H-i, c)
	}
}

func renderLine(dst draw.Image, p LinePayload) {
	c := parseColor(p.Color)
	s := p.StrokeWidth
	if s <= 0 {
		s = 2
	}
	for i := -s / 2; i <= s/2; i++ {
		drawLine(dst, p.X1+i, p.Y1, p.X2+i, p.Y2, c)
		drawLine(dst, p.X1, p.Y1+i, p.X2, p.Y2+i, c)
	}
}

func renderArrow(dst draw.Image, p ArrowPayload) {
	renderLine(dst, p.LinePayload)
	c := parseColor(p.Color)
	head := p.HeadSize
	if head <= 0 {
		head = 14
	}
	angle := math.Atan2(float64(p.Y2-p.Y1), float64(p.X2-p.X1))
	a1 := angle + math.Pi*0.82
	a2 := angle - math.Pi*0.82
	x3 := p.X2 + int(float64(head)*math.Cos(a1))
	y3 := p.Y2 + int(float64(head)*math.Sin(a1))
	x4 := p.X2 + int(float64(head)*math.Cos(a2))
	y4 := p.Y2 + int(float64(head)*math.Sin(a2))
	drawLine(dst, p.X2, p.Y2, x3, y3, c)
	drawLine(dst, p.X2, p.Y2, x4, y4, c)
}

func renderText(dst draw.Image, p TextPayload) {
	c := parseColor(p.Color)
	size := p.Size
	if size <= 0 {
		size = 2
	}
	x := p.X
	for _, r := range p.Text {
		renderGlyph(dst, x, p.Y, r, c, size)
		x += 6 * size
	}
}

func applyBlur(dst draw.Image, p BlurPayload) {
	if p.Radius <= 0 {
		p.Radius = 2
	}
	bounds := dst.Bounds()
	r := image.Rect(p.X, p.Y, p.X+p.W, p.Y+p.H).Intersect(bounds)
	src := image.NewRGBA(bounds)
	draw.Draw(src, bounds, dst, bounds.Min, draw.Src)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			var rs, gs, bs, count int
			for yy := y - p.Radius; yy <= y+p.Radius; yy++ {
				for xx := x - p.Radius; xx <= x+p.Radius; xx++ {
					if !image.Pt(xx, yy).In(bounds) {
						continue
					}
					cr, cg, cb, _ := src.At(xx, yy).RGBA()
					rs += int(cr >> 8)
					gs += int(cg >> 8)
					bs += int(cb >> 8)
					count++
				}
			}
			if count > 0 {
				dst.Set(x, y, color.RGBA{R: uint8(rs / count), G: uint8(gs / count), B: uint8(bs / count), A: 255})
			}
		}
	}
}

func applyPixelate(dst draw.Image, p PixelatePayload) {
	if p.Size <= 1 {
		p.Size = 8
	}
	bounds := dst.Bounds()
	r := image.Rect(p.X, p.Y, p.X+p.W, p.Y+p.H).Intersect(bounds)
	for y := r.Min.Y; y < r.Max.Y; y += p.Size {
		for x := r.Min.X; x < r.Max.X; x += p.Size {
			x2 := min(x+p.Size, r.Max.X)
			y2 := min(y+p.Size, r.Max.Y)
			cr, cg, cb, _ := dst.At(x, y).RGBA()
			block := color.RGBA{R: uint8(cr >> 8), G: uint8(cg >> 8), B: uint8(cb >> 8), A: 255}
			draw.Draw(dst, image.Rect(x, y, x2, y2), image.NewUniform(block), image.Point{}, draw.Src)
		}
	}
}

func drawLine(img draw.Image, x0, y0, x1, y1 int, c color.RGBA) {
	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		if image.Pt(x0, y0).In(img.Bounds()) {
			img.Set(x0, y0, c)
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func renderGlyph(dst draw.Image, x, y int, r rune, c color.RGBA, scale int) {
	seed := int(r)
	for gy := 0; gy < 7; gy++ {
		for gx := 0; gx < 5; gx++ {
			if ((seed >> ((gx + gy) % 8)) & 1) == 0 {
				continue
			}
			for sy := 0; sy < scale; sy++ {
				for sx := 0; sx < scale; sx++ {
					px := x + gx*scale + sx
					py := y + gy*scale + sy
					if image.Pt(px, py).In(dst.Bounds()) {
						dst.Set(px, py, c)
					}
				}
			}
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
