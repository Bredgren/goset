package main

import (
	"image/color"
	"math/rand"
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

func newMainMenuState(display *ggweb.Surface) GameState {
	m := MainMenu{
		nextState: mainMenuState,
	}
	btnX := 40.0
	btnY := display.Rect().CenterY()
	m.buttons = []*Button{
		newTextButton("Play", btnX, btnY, func() {
			ggweb.Log("Play")
			m.nextState = playState
		}),
		newTextButton("Leaderboard", btnX, btnY+50, func() {
			ggweb.Log("Leaderboard")
			m.nextState = leaderboardState
		}),
		newTextButton("Help", btnX, btnY+100, func() {
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
		if (cg == nil || !cg.Active) && rand.Float64() < 0.01 {
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

	titleFont := ggweb.Font{
		Size: 75,
	}
	display.SetFont(&titleFont)
	display.StyleColor(ggweb.Fill, color.White)
	display.SetTextAlign(ggweb.TextAlignCenter)
	display.SetTextBaseline(ggweb.TextBaselineMiddle)
	display.DrawText(ggweb.Fill, "SET", display.Rect().CenterX(), 10+float64(titleFont.Size))

	display.SetCursor(ggweb.CursorDefault)
	for _, b := range m.buttons {
		b.drawTo(display)
		if b.State == buttonHover {
			display.SetCursor(ggweb.CursorPointer)
		}
	}
}
