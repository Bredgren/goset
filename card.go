package main

import (
	"image/color"
	"math"

	"github.com/Bredgren/gogame/geo"
	"github.com/Bredgren/gogame/ggweb"
)

type count int

const (
	one count = iota
	two
	three
)

type shape int

const (
	oval shape = iota
	diamond
	tilde
)

var (
	red    = color.RGBA{255, 0, 0, 255}
	green  = color.RGBA{0, 230, 0, 255}
	purple = color.RGBA{255, 0, 255, 255}
)

type fill int

const (
	empty fill = iota
	solid
	line
)

type card struct {
	count count
	fill  fill
	color color.Color
	shape shape
}

func (c *card) surface(w, h float64) *ggweb.Surface {
	surf := ggweb.NewSurface(int(w), int(h))
	surf.StyleColor(ggweb.Fill, color.White)
	surf.DrawRect(ggweb.Fill, surf.Rect())
	var shapeSurf *ggweb.Surface
	shapeW := w * 0.75
	shapeH := shapeW / 2
	switch c.shape {
	case oval:
		shapeSurf = makeOval(shapeW, shapeH, c.color, c.fill)
	case diamond:
		shapeSurf = makeDiamond(shapeW, shapeH, c.color, c.fill)
	case tilde:
		shapeSurf = makeTilde(shapeW, shapeH, c.color, c.fill)
	}
	centerRect := geo.Rect{W: shapeW, H: shapeH}
	centerRect.SetCenterX(w / 2)
	centerRect.SetCenterY(h / 2)
	switch c.count {
	case one:
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y)
	case two:
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y-shapeH*0.75)
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y+shapeH*0.75)
	case three:
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y-shapeH*1.25)
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y)
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y+shapeH*1.25)
	}
	return surf
}

func makeOval(w, h float64, c color.Color, f fill) *ggweb.Surface {
	s := ggweb.NewSurface(int(w), int(h))
	if f == solid {
		s.StyleColor(ggweb.Fill, c)
		s.DrawRect(ggweb.Fill, geo.Rect{X: h / 2, Y: 0, W: w - h, H: h})
		s.DrawArc(ggweb.Fill, geo.Rect{X: 0, Y: 0, W: h, H: h}, math.Pi/2, 3*math.Pi/2, true)
		s.DrawArc(ggweb.Fill, geo.Rect{X: w - h, Y: 0, W: h, H: h}, -math.Pi/2, math.Pi/2, true)
	} else {
		s.StyleColor(ggweb.Stroke, c)
		lw := 0.03 * w
		s.SetLineWidth(lw)
		if f == line {
			s.Blit(makeLines(w, h, c, lw), 0, 0)
			mask := ggweb.NewSurface(int(w), int(h))
			mask.StyleColor(ggweb.Fill, color.White)
			mask.DrawRect(ggweb.Fill, geo.Rect{X: h / 2, Y: 0, W: w - h, H: h})
			mask.DrawArc(ggweb.Fill, geo.Rect{X: 0, Y: 0, W: h, H: h}, math.Pi/2, 3*math.Pi/2, true)
			mask.DrawArc(ggweb.Fill, geo.Rect{X: w - h, Y: 0, W: h, H: h}, -math.Pi/2, math.Pi/2, true)
			s.SetCompositeOp(ggweb.DestinationIn)
			s.Blit(mask, 0, 0)
			s.SetCompositeOp(ggweb.SourceOver)
		}
		s.DrawLine(h/2, 1, w-h/2, 1)
		s.DrawLine(h/2, h-1, w-h/2, h-1)
		s.DrawArc(ggweb.Stroke, geo.Rect{X: 1, Y: 1, W: h, H: h - 2}, math.Pi/2, 3*math.Pi/2, true)
		s.DrawArc(ggweb.Stroke, geo.Rect{X: w - h - 1, Y: 1, W: h, H: h - 2}, -math.Pi/2, math.Pi/2, true)
	}
	return s
}

func makeDiamond(w, h float64, c color.Color, f fill) *ggweb.Surface {
	s := ggweb.NewSurface(int(w), int(h))
	path := ggweb.NewPath()
	path.MoveTo(w/2, 0)
	path.LineTo(w, h/2)
	path.LineTo(w/2, h)
	path.LineTo(0, h/2)
	path.Close()
	if f == solid {
		s.StyleColor(ggweb.Fill, c)
		s.DrawPath(ggweb.Fill, path)
	} else {
		s.StyleColor(ggweb.Stroke, c)
		lw := 0.03 * w
		s.SetLineWidth(lw)
		if f == line {
			s.Blit(makeLines(w, h, c, lw), 0, 0)
			mask := ggweb.NewSurface(int(w), int(h))
			mask.StyleColor(ggweb.Fill, color.White)
			mask.DrawPath(ggweb.Fill, path)
			s.SetCompositeOp(ggweb.DestinationIn)
			s.Blit(mask, 0, 0)
			s.SetCompositeOp(ggweb.SourceOver)
		}
		// points[0][1] += style.Width / 2
		// points[1][0] -= style.Width / 2
		// points[2][1] -= style.Width / 2
		// points[3][0] += style.Width / 2
		// points[4][1] += style.Width / 2
		s.DrawPath(ggweb.Stroke, path)
	}
	return s
}

func makeTilde(w, h float64, c color.Color, f fill) *ggweb.Surface {
	s := ggweb.NewSurface(int(w), int(h))
	p1x, p1y := w*0.04, h*0.5
	p2x, p2y := w*0.77, h*0.125
	cp1x, cp1y := w*0.2, h*0.5-h
	cp1x2, cp1y2 := p2x-w*0.2, h*0.5
	cp2x, cp2y := p2x+w*0.1, 0.0
	cp2x2, cp2y2 := w+w*0.01, 0.0

	path := ggweb.NewPath()
	path.MoveTo(p1x, p1y)
	path.BezierCurveTo(cp1x, cp1y, cp1x2, cp1y2, p2x, p2y)
	path.BezierCurveTo(cp2x, cp2y, cp2x2, cp2y2, w-p1x, h-p1y)
	path.BezierCurveTo(w-cp1x, h-cp1y, w-cp1x2, h-cp1y2, w-p2x, h-p2y)
	path.BezierCurveTo(w-cp2x, h-cp2y, w-cp2x2, h-cp2y2, p1x, p1y)

	if f == solid {
		s.StyleColor(ggweb.Fill, c)
		s.DrawPath(ggweb.Fill, path)
	} else {
		s.StyleColor(ggweb.Stroke, c)
		lw := 0.03 * w
		s.SetLineWidth(lw)
		if f == line {
			s.Blit(makeLines(w, h, c, lw), 0, 0)
			mask := ggweb.NewSurface(int(w), int(h))
			mask.StyleColor(ggweb.Fill, color.White)
			mask.DrawPath(ggweb.Fill, path)
			s.SetCompositeOp(ggweb.DestinationIn)
			s.Blit(mask, 0, 0)
			s.SetCompositeOp(ggweb.SourceOver)
		}
		s.DrawPath(ggweb.Stroke, path)
	}
	return s
}

func makeLines(w, h float64, c color.Color, lw float64) *ggweb.Surface {
	s := ggweb.NewSurface(int(w), int(h))
	s.StyleColor(ggweb.Stroke, c)
	s.SetLineWidth(lw)
	for x := 0.1 * w; x < w; x += 0.1 * w {
		xx := x
		if int(math.Ceil(lw))%2 != 0 {
			xx = math.Floor(x) + 0.5
		}
		s.DrawLine(xx, 0, xx, h)
	}
	return s
}

func isSet(c1, c2, c3 card) bool {
	return ((c1.count == c2.count && c1.count == c3.count) || (c1.count != c2.count && c1.count != c3.count && c2.count != c3.count)) &&
		((c1.fill == c2.fill && c1.fill == c3.fill) || (c1.fill != c2.fill && c1.fill != c3.fill && c2.fill != c3.fill)) &&
		((c1.color == c2.color && c1.color == c3.color) || (c1.color != c2.color && c1.color != c3.color && c2.color != c3.color)) &&
		((c1.shape == c2.shape && c1.shape == c3.shape) || (c1.shape != c2.shape && c1.shape != c3.shape && c2.shape != c3.shape))
}
