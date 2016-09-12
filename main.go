package main

import (
	"math"
	"time"

	"github.com/Bredgren/gogame"
	"github.com/Bredgren/gogame/composite"
	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/geo"
)

func main() {
	ready := gogame.Ready()
	go func() {
		<-ready
		setup()
		gogame.SetMainLoop(mainLoop)
	}()
}

func setup() {
	width, height := 500, 700
	display := gogame.MainDisplay()
	display.SetMode(width, height)
	display.Fill(gogame.FillBlack)
	display.Flip()
}

func mainLoop(t time.Duration) {
	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
		switch evt.Type {
		case event.Quit:
			gogame.UnsetMainLoop()
		case event.VideoResize:
			// data := evt.Data.(event.ResizeData)
		case event.KeyDown:
			// data := evt.Data.(event.KeyData)
		case event.KeyUp:
			// data := evt.Data.(event.KeyData)
		case event.MouseButtonDown:
			// data := evt.Data.(event.MouseData)
		case event.MouseButtonUp:
			// data := evt.Data.(event.MouseData)
		case event.MouseMotion:
			// data := evt.Data.(event.MouseMotionData)
		}
	}
	display := gogame.MainDisplay()
	display.Fill(gogame.FillBlack)

	display.Blit(makeCard(60, 100, one, empty, red, oval), 10, 10)
	display.Blit(makeCard(60, 100, two, empty, red, oval), 80, 10)
	display.Blit(makeCard(60, 100, three, empty, red, oval), 150, 10)

	display.Blit(makeCard(60, 100, one, solid, red, oval), 10, 120)
	display.Blit(makeCard(60, 100, two, solid, red, oval), 80, 120)
	display.Blit(makeCard(60, 100, three, solid, red, oval), 150, 120)

	display.Blit(makeCard(60, 100, one, line, red, oval), 10, 230)
	display.Blit(makeCard(60, 100, two, line, red, oval), 80, 230)
	display.Blit(makeCard(60, 100, three, line, red, oval), 150, 230)

	display.Blit(makeCard(60, 100, one, empty, green, dimond), 220, 10)
	display.Blit(makeCard(60, 100, two, empty, green, dimond), 290, 10)
	display.Blit(makeCard(60, 100, three, empty, green, dimond), 360, 10)

	display.Blit(makeCard(60, 100, one, solid, green, dimond), 220, 120)
	display.Blit(makeCard(60, 100, two, solid, green, dimond), 290, 120)
	display.Blit(makeCard(60, 100, three, solid, green, dimond), 360, 120)

	display.Blit(makeCard(60, 100, one, line, green, dimond), 220, 230)
	display.Blit(makeCard(60, 100, two, line, green, dimond), 290, 230)
	display.Blit(makeCard(60, 100, three, line, green, dimond), 360, 230)

	display.Blit(makeCard(60, 100, one, empty, purple, tilde), 10, 340)
	display.Blit(makeCard(60, 100, two, empty, purple, tilde), 80, 340)
	display.Blit(makeCard(60, 100, three, empty, purple, tilde), 150, 340)

	display.Blit(makeCard(60, 100, one, solid, purple, tilde), 10, 450)
	display.Blit(makeCard(60, 100, two, solid, purple, tilde), 80, 450)
	display.Blit(makeCard(60, 100, three, solid, purple, tilde), 150, 450)

	display.Blit(makeCard(60, 100, one, line, purple, tilde), 10, 560)
	display.Blit(makeCard(60, 100, two, line, purple, tilde), 80, 560)
	display.Blit(makeCard(60, 100, three, line, purple, tilde), 150, 560)

	display.Flip()
}

func makeOval(w, h float64, c color, f fill) gogame.Surface {
	color := gogame.Color(c)
	s := gogame.NewSurface(int(w), int(h))
	if f == solid {
		style := &gogame.FillStyle{
			Colorer: color,
		}
		s.DrawRect(geo.Rect{X: h / 2, Y: 0, W: w - h, H: h}, style)
		s.DrawArc(geo.Rect{X: 0, Y: 0, W: h, H: h}, math.Pi/2, 3*math.Pi/2, style)
		s.DrawArc(geo.Rect{X: w - h, Y: 0, W: h, H: h}, -math.Pi/2, math.Pi/2, style)
	} else {
		style := &gogame.StrokeStyle{
			Colorer: color,
			Width:   0.03 * float64(w),
		}
		if f == line {
			s.Blit(makeLines(w, h, style), 0, 0)
			mask := gogame.NewSurface(int(w), int(h))
			mask.DrawRect(geo.Rect{X: h / 2, Y: 0, W: w - h, H: h}, gogame.FillWhite)
			mask.DrawArc(geo.Rect{X: 0, Y: 0, W: h, H: h}, math.Pi/2, 3*math.Pi/2, gogame.FillWhite)
			mask.DrawArc(geo.Rect{X: w - h, Y: 0, W: h, H: h}, -math.Pi/2, math.Pi/2, gogame.FillWhite)
			s.BlitComp(mask, 0, 0, composite.DestinationIn)
		}
		s.DrawLine(h/2, 1, w-h/2, 1, style)
		s.DrawLine(h/2, h-1, w-h/2, h-1, style)
		s.DrawArc(geo.Rect{X: 1, Y: 1, W: h, H: h - 2}, math.Pi/2, 3*math.Pi/2, style)
		s.DrawArc(geo.Rect{X: w - h - 1, Y: 1, W: h, H: h - 2}, -math.Pi/2, math.Pi/2, style)
	}
	return s
}

func makeDimond(w, h float64, c color, f fill) gogame.Surface {
	color := gogame.Color(c)
	s := gogame.NewSurface(int(w), int(h))
	points := [][2]float64{
		{w / 2, 0},
		{w, h / 2},
		{w / 2, h},
		{0, h / 2},
		{w / 2, 0},
	}
	if f == solid {
		s.DrawLines(points, &gogame.FillStyle{Colorer: color})
	} else {
		style := &gogame.StrokeStyle{
			Colorer: color,
			Width:   0.03 * w,
		}
		if f == line {
			s.Blit(makeLines(w, h, style), 0, 0)
			mask := gogame.NewSurface(int(w), int(h))
			mask.DrawLines(points, gogame.FillWhite)
			s.BlitComp(mask, 0, 0, composite.DestinationIn)
		}
		s.DrawLines(points, style)
	}
	return s
}

func makeTilde(w, h float64, c color, f fill) gogame.Surface {
	color := gogame.Color(c)
	s := gogame.NewSurface(int(w), int(h))
	p1x, p1y := w*0.04, h*0.5
	p2x, p2y := w*0.77, h*0.125
	cp1x, cp1y := w*0.2, h*0.5-h
	cp1x2, cp1y2 := p2x-w*0.2, h*0.5
	cp2x, cp2y := p2x+w*0.1, 0.0
	cp2x2, cp2y2 := w+w*0.01, 0.0
	points := [][2]float64{
		{p1x, p1y}, {cp1x, cp1y}, {cp1x2, cp1y2},
		{p2x, p2y}, {cp2x, cp2y}, {cp2x2, cp2y2},
		{w - p1x, h - p1y}, {w - cp1x, h - cp1y}, {w - cp1x2, h - cp1y2},
		{w - p2x, h - p2y}, {w - cp2x, h - cp2y}, {w - cp2x2, h - cp2y2},
		{p1x, p1y},
	}
	if f == solid {
		s.DrawBezierCurves(points, &gogame.FillStyle{Colorer: color})
	} else {
		style := &gogame.StrokeStyle{
			Colorer: color,
			Width:   0.03 * float64(w),
		}
		if f == line {
			s.Blit(makeLines(w, h, style), 0, 0)
			mask := gogame.NewSurface(int(w), int(h))
			mask.DrawBezierCurves(points, gogame.FillWhite)
			s.BlitComp(mask, 0, 0, composite.DestinationIn)
		}
		s.DrawBezierCurves(points, style)
	}
	return s
}

func makeLines(w, h float64, style *gogame.StrokeStyle) gogame.Surface {
	s := gogame.NewSurface(int(w), int(h))
	wf := float64(w)
	for x := 0.1 * wf; x < wf; x += 0.1 * wf {
		xx := x
		if int(math.Ceil(style.Width))%2 != 0 {
			xx = math.Floor(x) + 0.5
		}
		s.DrawLine(xx, 0, xx, h, style)
	}
	return s
}

func makeCard(w, h float64, n count, f fill, c color, s shape) gogame.Surface {
	surf := gogame.NewSurface(int(w), int(h))
	surf.Fill(gogame.FillWhite)
	var shapeSurf gogame.Surface
	shapeW := w * 0.75
	shapeH := shapeW / 2
	switch s {
	case oval:
		shapeSurf = makeOval(shapeW, shapeH, c, f)
	case dimond:
		shapeSurf = makeDimond(shapeW, shapeH, c, f)
	case tilde:
		shapeSurf = makeTilde(shapeW, shapeH, c, f)
	}
	centerRect := geo.Rect{W: shapeW, H: shapeH}
	centerRect.SetCenterX(w / 2)
	centerRect.SetCenterY(h / 2)
	switch n {
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

type count int

const (
	one count = iota
	two
	three
)

type shape int

const (
	oval shape = iota
	dimond
	tilde
)

type color gogame.Color

var (
	red    = color(gogame.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0})
	green  = color(gogame.Color{R: 0.0, G: 0.9, B: 0.0, A: 1.0})
	purple = color(gogame.Color{R: 1.0, G: 0.0, B: 1.0, A: 1.0})
)

type fill int

const (
	empty fill = iota
	solid
	line
)

type card struct {
	count count
	shape shape
	color color
	fill  fill
}

func (c *card) surface(w, h int) gogame.Surface {
	s := gogame.NewSurface(w, h)
	return s
}
