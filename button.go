package main

import (
	"image/color"

	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/geo"
	"github.com/Bredgren/gogame/ggweb"
)

type ButtonState int

const (
	buttonIdle = ButtonState(iota)
	buttonHover
	buttonSelect
)

type Button struct {
	Rect     geo.Rect
	Surfs    map[ButtonState]*ggweb.Surface
	Callback func()
	State    ButtonState
}

func newTextButton(text string, size, x, y float64, callback func()) *Button {
	textHeight := int(size)
	const padding = 5
	font := ggweb.Font{
		Size: size,
	}

	s := ggweb.NewSurface(0, 0)
	s.SetFont(&font)
	textWidth := int(s.TextWidth(text))

	idleSurf := ggweb.NewSurface(textWidth+padding*2, textHeight+padding*2)
	idleSurf.SetFont(&font)

	idleSurf.StyleColor(ggweb.Fill, color.RGBA{A: 100})
	idleSurf.DrawRect(ggweb.Fill, idleSurf.Rect())

	idleSurf.StyleColor(ggweb.Stroke, color.White)
	idleSurf.SetLineWidth(size / 5)
	idleSurf.DrawRect(ggweb.Stroke, idleSurf.Rect())

	idleSurf.StyleColor(ggweb.Fill, color.White)
	idleSurf.SetTextAlign(ggweb.TextAlignCenter)
	idleSurf.SetTextBaseline(ggweb.TextBaselineMiddle)
	idleSurf.DrawText(ggweb.Fill, text, idleSurf.Rect().CenterX(), idleSurf.Rect().CenterY())

	hoverSurf := ggweb.NewSurface(textWidth+padding*2, textHeight+padding*2)
	hoverSurf.SetFont(&font)

	hoverSurf.StyleColor(ggweb.Fill, color.White)
	hoverSurf.DrawRect(ggweb.Fill, hoverSurf.Rect())

	hoverSurf.StyleColor(ggweb.Fill, color.Black)
	hoverSurf.SetTextAlign(ggweb.TextAlignCenter)
	hoverSurf.SetTextBaseline(ggweb.TextBaselineMiddle)
	hoverSurf.DrawText(ggweb.Fill, text, hoverSurf.Rect().CenterX(), hoverSurf.Rect().CenterY())

	r := idleSurf.Rect()
	r.SetTopLeft(x, y)
	return &Button{
		Rect: r,
		Surfs: map[ButtonState]*ggweb.Surface{
			buttonIdle:   idleSurf,
			buttonHover:  hoverSurf,
			buttonSelect: hoverSurf,
		},
		Callback: callback,
	}
}

func (b *Button) handleEvent(evt event.Event) {
	switch evt.Type {
	case event.MouseButtonUp:
		data := evt.Data.(event.MouseData)
		if data.Button == 0 && b.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
			b.State = buttonSelect
			b.Callback()
		}
	case event.MouseMotion:
		data := evt.Data.(event.MouseMotionData)
		if b.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
			if b.State != buttonHover {
				// globalState.sound["btnHover"].PlayFromStart()
			}
			b.State = buttonHover
		} else {
			b.State = buttonIdle
		}
	}
}

func (b *Button) drawTo(surf *ggweb.Surface) {
	surf.Blit(b.Surfs[b.State], b.Rect.X, b.Rect.Y)
}
