package main

import (
	"github.com/Bredgren/gogame"
	"github.com/Bredgren/gogame/composite"
	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/geo"
	"github.com/Bredgren/gogame/ui"
)

func handleCommonEvents(evt event.Event) {
	switch evt.Type {
	case event.Quit:
		gogame.UnsetMainLoop()
	case event.VideoResize:
		// data := evt.Data.(event.ResizeData)
	}
}

const (
	btnIdle ui.ButtonState = iota
	btnHover
	btnSelect
)

func updateButtons(evt event.Event, buttons []*ui.BasicButton) {
	hovering := false
	switch evt.Type {
	case event.MouseButtonDown:
		data := evt.Data.(event.MouseData)
		if data.Button == 0 {
			for _, btn := range buttons {
				if btn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
					btn.State = btnSelect
				}
			}
		}
	case event.MouseButtonUp:
		data := evt.Data.(event.MouseData)
		if data.Button == 0 {
			for _, btn := range buttons {
				if btn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
					btn.Select()
					btn.State = btnIdle
				}
			}
		}
	case event.MouseMotion:
		data := evt.Data.(event.MouseMotionData)
		for _, btn := range buttons {
			if btn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
				prev := btn.State
				if data.Buttons[0] {
					btn.State = btnSelect
				} else {
					btn.State = btnHover
					hovering = true
				}
				if prev != btn.State {
					// globalState.sound["btnHover"].PlayFromStart()
				}
			} else {
				btn.State = btnIdle
			}
		}
	}
	if hovering {
		gogame.MainDisplay().SetCursor(gogame.CursorPointer)
	} else {
		gogame.MainDisplay().SetCursor(gogame.CursorDefault)
	}
}

var (
	btnFont = &gogame.Font{
		Size: 20,
	}
	btnTextStyle = &gogame.TextStyle{
		Colorer:  gogame.White,
		Align:    gogame.TextAlignCenter,
		Baseline: gogame.TextBaselineMiddle,
	}
	btnFill = &gogame.FillStyle{
		Colorer: gogame.Color{A: 0.2},
	}
	btnOutline = &gogame.StrokeStyle{
		Colorer: gogame.White,
		Width:   5,
	}
	btnPadding = 10.0
)

func btnRect(text string) geo.Rect {
	return geo.Rect{
		W: float64(btnFont.Width(text, btnTextStyle)) + btnPadding*2,
		H: float64(btnFont.Size) + btnPadding*2,
	}
}

func makeBtn(text string, callback func()) *ui.BasicButton {
	r := btnRect(text)
	idle := gogame.NewSurface(int(r.W), int(r.H))
	idle.Fill(btnFill)
	idle.DrawText(text, r.W/2, r.H/2, btnFont, btnTextStyle)

	hover := idle.Copy()
	hover.DrawRect(r, btnOutline)

	sel := idle.Copy()
	invert := gogame.NewSurface(int(r.W), int(r.H))
	invert.Fill(gogame.FillWhite)
	sel.BlitComp(invert, 0, 0, composite.Difference)

	return &ui.BasicButton{
		Rect:        r,
		DefaultSurf: idle,
		StateSurfs: map[ui.ButtonState]gogame.Surface{
			btnIdle:   idle,
			btnHover:  hover,
			btnSelect: sel,
		},
		Select: callback,
		State:  btnIdle,
	}
}

func makeBtnShape(shape gogame.Surface, callback func()) *ui.BasicButton {
	r := shape.Rect()
	idle := gogame.NewSurface(int(r.W), int(r.H))
	idle.Fill(btnFill)
	idle.Blit(shape, 0, 0)

	hover := idle.Copy()
	hover.DrawRect(r, btnOutline)

	sel := idle.Copy()
	invert := gogame.NewSurface(int(r.W), int(r.H))
	invert.Fill(gogame.FillWhite)
	sel.BlitComp(invert, 0, 0, composite.Difference)

	return &ui.BasicButton{
		Rect:        r,
		DefaultSurf: idle,
		StateSurfs: map[ui.ButtonState]gogame.Surface{
			btnIdle:   idle,
			btnHover:  hover,
			btnSelect: sel,
		},
		Select: callback,
		State:  btnIdle,
	}
}

func backArrowBtn(callback func()) *ui.BasicButton {
	w, h := 40.0, 40.0
	surf := gogame.NewSurface(int(w), int(h))
	surf.Fill(btnFill)

	shapeFill := gogame.FillStyle{
		Colorer: btnTextStyle.Colorer,
	}

	padding := 0.1
	head := 0.5
	surf.DrawLines([][2]float64{
		{padding * w, 0.5 * h},
		{head * w, padding * h},
		{head * w, h - padding*h},
	}, &shapeFill)

	width := 0.4
	surf.DrawRect(geo.Rect{
		X: head * w,
		Y: h/2 - (width*h)/2,
		W: (w - padding*w) - head*w,
		H: width * h,
	}, &shapeFill)

	return makeBtnShape(surf, callback)
}
