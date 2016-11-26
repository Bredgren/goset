package main

import (
	"image/color"
	"time"

	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/fsm"
	"github.com/Bredgren/gogame/geo"
	"github.com/Bredgren/gogame/ggweb"
)

type MainMenu struct {
	nextState  fsm.State
	buttons    []*Button
	cardGroups [4]*FlyingCardGroup
}

func newMainMenuState() GameState {
	m := MainMenu{
		nextState: mainMenuState,
	}
	m.buttons = []*Button{
		newTextButton("Play", 10, 10, func() {
			ggweb.Log("Play")
			m.nextState = playState
		}),
		newTextButton("Leaderboard", 10, 50, func() {
			ggweb.Log("Leaderboard")
			m.nextState = leaderboardState
		}),
		newTextButton("Help", 10, 90, func() {
			ggweb.Log("Help")
			m.nextState = helpState
		}),
	}
	// TODO: Add resume button if there is a saved game
	for i := range m.cardGroups {
		m.cardGroups[i] = &FlyingCardGroup{}
	}
	return &m
}

func (m *MainMenu) Update(g *game, t, dt time.Duration) fsm.State {
	m.handleEvents()
	for _, cg := range m.cardGroups {
		if cg == nil || !cg.Active {
			r := g.display.Rect()
			mass := geo.RandNum(5, 15)
			targetDist := geo.RandNum(50, 60)
			k := geo.RandNum(35, 45)
			damp := geo.RandNum(14, 16)
			const thick = 75
			cg.Activate(geo.RandVecRects([]geo.Rect{
				{X: r.Left() - thick, Y: -thick, W: r.W + thick, H: thick},
				{X: r.Right(), Y: -thick, W: thick, H: r.H + thick},
				{X: r.Left(), Y: r.Bottom(), W: r.W + thick, H: thick},
				{X: r.Left() - thick, Y: r.Top(), W: thick, H: r.H + thick},
			}), mass, targetDist, k, damp)
		}
		cg.Update(t, dt)
	}

	m.draw(g.display, t, dt)
	return m.nextState
}

func (m *MainMenu) handleEvents() {
	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
		// handleCommonEvents(evt)
		// updateButtons(evt, s.buttons)
		for _, b := range m.buttons {
			b.handleEvent(evt)
		}
	}
}

func (m *MainMenu) draw(display *ggweb.Surface, t, dt time.Duration) {
	display.StyleColor(ggweb.Fill, color.Black)
	display.DrawRect(ggweb.Fill, display.Rect())

	for _, cg := range m.cardGroups {
		cg.Draw(display, t)
	}

	// 	// Draw tItle
	// 	titleFont := gogame.Font{
	// 		Size: 75,
	// 	}
	// 	titleStyle := gogame.TextStyle{
	// 		Colorer:  gogame.White,
	// 		Align:    gogame.TextAlignCenter,
	// 		Baseline: gogame.TextBaselineMiddle,
	// 	}
	// 	display.DrawText("SET", display.Rect().CenterX(), 10+float64(titleFont.Size), &titleFont, &titleStyle)

	display.SetCursor(ggweb.CursorDefault)
	for _, b := range m.buttons {
		b.drawTo(display)
		if b.State == buttonHover {
			display.SetCursor(ggweb.CursorPointer)
		}
	}
}

// func (s *mainMenuState) makeBtns() {
// 	btnSpacing := 10.0

// 	s.resumeBtn = makeBtn("Resume", func() {
// 		gogame.Log("TODO: handle resume button")
// 		// s.nextState = globalState.playState
// 	})

// 	playBtn := makeBtn("Play", func() {
// 		s.nextState = globalState.playState
// 	})
// 	playBtn.Rect.SetMidLeft(40, gogame.MainDisplay().Rect().CenterY())
// 	s.buttons = append(s.buttons, playBtn)

// 	s.resumeBtn.Rect.SetLeft(playBtn.Rect.Left())
// 	s.resumeBtn.Rect.SetBottom(playBtn.Rect.Top() - btnSpacing)

// 	leaderBtn := makeBtn("Leaderboard", func() {
// 		s.nextState = globalState.leaderboardState
// 	})
// 	leaderBtn.Rect.SetLeft(playBtn.Rect.Left())
// 	leaderBtn.Rect.SetTop(playBtn.Rect.Bottom() + btnSpacing)
// 	s.buttons = append(s.buttons, leaderBtn)

// 	helpBtn := makeBtn("Help", func() {
// 		gogame.Log("TODO: handle help button")
// 		// s.nextState = globalState.helpState
// 	})
// 	helpBtn.Rect.SetLeft(playBtn.Rect.Left())
// 	helpBtn.Rect.SetTop(leaderBtn.Rect.Bottom() + btnSpacing)
// 	s.buttons = append(s.buttons, helpBtn)
// }
