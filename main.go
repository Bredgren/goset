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
	width, height := 500, 600
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

	display.Blit(makeOval(45, 20, red, solid), 10, 10)
	display.Blit(makeOval(45, 20, green, empty), 10, 40)
	display.Blit(makeOval(45, 20, purple, line), 10, 70)
	display.Blit(makeOval(100, 50, red, solid), 100, 10)
	display.Blit(makeOval(100, 50, green, empty), 100, 70)
	display.Blit(makeOval(100, 50, purple, line), 100, 130)

	display.Blit(makeDimond(45, 20, red, solid), 220, 10)
	display.Blit(makeDimond(45, 20, green, empty), 220, 40)
	display.Blit(makeDimond(45, 20, purple, line), 220, 70)

	display.Blit(makeTilde(45, 20, red, solid), 300, 10)
	display.Blit(makeTilde(45, 20, green, empty), 300, 40)
	display.Blit(makeTilde(45, 20, purple, line), 300, 70)

	display.Blit(makeTilde(45, 20, red, solid), 10, 550)
	display.DrawRect(geo.Rect{X: 10, Y: 550, W: 45, H: 20}, &gogame.StrokeStyle{
		Colorer: gogame.White,
		Width:   1,
	})

	display.Flip()
}

func makeOval(w, h float64, color gogame.Color, f fill) gogame.Surface {
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

func makeDimond(w, h float64, color gogame.Color, f fill) gogame.Surface {
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

func makeTilde(w, h float64, color gogame.Color, f fill) gogame.Surface {
	s := gogame.NewSurface(int(w), int(h))
	points := [][2]float64{
		{w * 0.05, h / 2}, {w * 0.15, h/2 - h}, {w - w/4 - w/10, h/8 + h},
		{w - w/4, h / 8}, {w - w/4 + w/10, h/8 - h}, {w, h/2 - h},
		{w, h / 2}, {w, h/2 + h}, {w/4 + w/10, h - h/8 - h*0.3},
		{w * 0.2, h - h/8}, {w/4 - w/10, h}, {-w * 0.01, h},
		{w * 0.05, h / 2},
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

// type scf struct {
// 	shape shape
// 	color color
// 	fill  fill
// }
//
// var shapes map[scf]gogame.Surface

// func makeShapes(w, h int) {
// 	for s := oval; s < tilde; s++ {
// 	}
// }

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
	red    = gogame.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0}
	green  = gogame.Color{R: 0.0, G: 1.0, B: 0.0, A: 1.0}
	purple = gogame.Color{R: 1.0, G: 0.0, B: 1.0, A: 1.0}
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
