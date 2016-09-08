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
	oval := makeOval(45, 20, red, solid)
	display.Blit(oval, 10, 10)
	oval2 := makeOval(45, 20, green, empty)
	display.Blit(oval2, 10, 40)
	oval3 := makeOval(45, 20, purple, line)
	display.Blit(oval3, 10, 70)

	oval4 := makeOval(100, 50, red, solid)
	display.Blit(oval4, 100, 10)
	oval5 := makeOval(100, 50, green, empty)
	display.Blit(oval5, 100, 70)
	oval6 := makeOval(100, 50, purple, line)
	display.Blit(oval6, 100, 130)

	display.Flip()
}

func makeOval(w, h int, color gogame.Color, f fill) gogame.Surface {
	s := gogame.NewSurface(w, h)
	if f == solid {
		style := &gogame.FillStyle{
			Colorer: color,
		}
		s.DrawRect(geo.Rect{X: float64(h) / 2, Y: 0, W: float64(w) - float64(h), H: float64(h)}, style)
		s.DrawArc(geo.Rect{X: 0, Y: 0, W: float64(h), H: float64(h)}, math.Pi/2, 3*math.Pi/2, style)
		s.DrawArc(geo.Rect{X: float64(w - h), Y: 0, W: float64(h), H: float64(h)}, -math.Pi/2, math.Pi/2, style)
	} else {
		style := &gogame.StrokeStyle{
			Colorer: color,
			Width:   0.03 * float64(w),
		}
		if f == line {
			s.Blit(makeLines(w, h, style), 0, 0)
			mask := gogame.NewSurface(w, h)
			mask.DrawRect(geo.Rect{X: float64(h) / 2, Y: 0, W: float64(w) - float64(h), H: float64(h)}, gogame.FillWhite)
			mask.DrawArc(geo.Rect{X: 0, Y: 0, W: float64(h), H: float64(h)}, math.Pi/2, 3*math.Pi/2, gogame.FillWhite)
			mask.DrawArc(geo.Rect{X: float64(w - h), Y: 0, W: float64(h), H: float64(h)}, -math.Pi/2, math.Pi/2, gogame.FillWhite)
			s.BlitComp(mask, 0, 0, composite.DestinationIn)
		}
		s.DrawLine(float64(h)/2, 1, float64(w)-float64(h)/2, 1, style)
		s.DrawLine(float64(h)/2, float64(h)-1, float64(w)-float64(h)/2, float64(h)-1, style)
		s.DrawArc(geo.Rect{X: 1, Y: 1, W: float64(h), H: float64(h) - 2}, math.Pi/2, 3*math.Pi/2, style)
		s.DrawArc(geo.Rect{X: float64(w-h) - 1, Y: 1, W: float64(h), H: float64(h) - 2}, -math.Pi/2, math.Pi/2, style)
	}
	return s
}

func makeLines(w, h int, style *gogame.StrokeStyle) gogame.Surface {
	s := gogame.NewSurface(w, h)
	wf := float64(w)
	for x := 0.1 * wf; x < wf; x += 0.1 * wf {
		xx := x
		if int(math.Ceil(style.Width))%2 != 0 {
			xx = math.Floor(x) + 0.5
		}
		s.DrawLine(xx, 0, xx, float64(h), style)
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
