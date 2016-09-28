package main

import (
	"time"

	"github.com/Bredgren/gogame"
	"github.com/Bredgren/gogame/composite"
	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/geo"
	"github.com/Bredgren/gogame/ui"
)

type mainMenuState struct {
	nextState    GameState
	btnFont      *gogame.Font
	btnTextStyle *gogame.TextStyle
	btnPadding   float64
	btnFill      *gogame.FillStyle
	btnOutline   *gogame.StrokeStyle
	buttons      []*ui.BasicButton
	resumeBtn    *ui.BasicButton
}

const (
	btnIdle ui.ButtonState = iota
	btnHover
	btnSelect
)

func (s *mainMenuState) Enter() {
	s.nextState = nil
	s.btnFont = &gogame.Font{
		Size: 20,
	}
	s.btnTextStyle = &gogame.TextStyle{
		Colorer:  gogame.White,
		Align:    gogame.TextAlignCenter,
		Baseline: gogame.TextBaselineMiddle,
	}
	s.btnPadding = 10
	s.btnFill = &gogame.FillStyle{
		Colorer: gogame.Color{A: 0.2},
	}
	s.btnOutline = &gogame.StrokeStyle{
		Colorer: gogame.White,
		Width:   5,
	}
	if len(s.buttons) == 0 {
		s.makeBtns()
	}
	// TODO: actually check for saved game
	savedGame := true
	ri := -1
	for i, b := range s.buttons {
		if b == s.resumeBtn {
			ri = i
			break
		}
	}
	if savedGame && ri < 0 {
		s.buttons = append(s.buttons, s.resumeBtn)
	} else if !savedGame && ri >= 0 {
		s.buttons = append(s.buttons[:ri], s.buttons[ri+1:]...)
	}
}

func (s *mainMenuState) Exit() {
}

func (s *mainMenuState) Update(t, dt time.Duration) {
	if s.nextState != nil {
		// TODO: leaving animation
		globalState.gameStateMgr.Goto(s.nextState)
	}
	s.handleEvents()
	s.draw()
}

func (s *mainMenuState) handleEvents() {
	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
		handleCommonEvents(evt)
		switch evt.Type {
		case event.MouseButtonDown:
			data := evt.Data.(event.MouseData)
			if data.Button == 0 {
				for _, btn := range s.buttons {
					if btn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
						btn.State = btnSelect
					}
				}
			}
		case event.MouseButtonUp:
			data := evt.Data.(event.MouseData)
			if data.Button == 0 {
				for _, btn := range s.buttons {
					if btn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
						btn.Select()
						btn.State = btnIdle
					}
				}
			}
		case event.MouseMotion:
			data := evt.Data.(event.MouseMotionData)
			for _, btn := range s.buttons {
				if btn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
					if data.Buttons[0] {
						btn.State = btnSelect
					} else {
						btn.State = btnHover
					}
				} else {
					btn.State = btnIdle
				}
			}
		}
	}
}

func (s *mainMenuState) draw() {
	display := gogame.MainDisplay()
	display.Fill(gogame.FillBlack)
	for _, btn := range s.buttons {
		btn.DrawTo(display)
	}
	display.Flip()
}

func (s *mainMenuState) makeBtns() {
	display := gogame.MainDisplay()
	dr := display.Rect()

	s.resumeBtn = s.makeBtn("Resume", func() {
		gogame.Log("TODO: handle resume button")
	})

	b := s.makeBtn("Play", func() {
		s.nextState = globalState.playState
	})
	b.Rect.SetCenter(dr.Center())
	s.buttons = append(s.buttons, b)

	s.resumeBtn.Rect.SetCenterX(b.Rect.CenterX())
	s.resumeBtn.Rect.SetBottom(b.Rect.Top() - s.btnPadding)

	// s.makeHelpBtn()
	// s.makeLeaderboardBtn()
}

func (s *mainMenuState) btnRect(text string) geo.Rect {
	return geo.Rect{
		W: float64(s.btnFont.Width(text, s.btnTextStyle)) + s.btnPadding*2,
		H: float64(s.btnFont.Size) + s.btnPadding*2,
	}
}

func (s *mainMenuState) makeBtn(text string, callback func()) *ui.BasicButton {
	r := s.btnRect(text)
	idle := gogame.NewSurface(int(r.W), int(r.H))
	idle.Fill(s.btnFill)
	idle.DrawText(text, r.W/2, r.H/2, s.btnFont, s.btnTextStyle)

	hover := idle.Copy()
	hover.DrawRect(r, s.btnOutline)

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
