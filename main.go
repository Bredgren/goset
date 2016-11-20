// TODO:
//   - saving/restoring game
//   - main menu
//   - "tutorial"
//   - leaderboard
package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/Bredgren/gogame/fsm"
	"github.com/Bredgren/gogame/ggweb"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ggweb.Init(func() {
		g := newGame()
		ggweb.SetMainLoop(g.mainLoop)
	})
}

const (
	displayW, displayH = 680, 500
)

const (
	mainMenuState    fsm.State = "mainmenu"
	playState        fsm.State = "play"
	leaderboardState fsm.State = "leaderboard"
)

type gameState interface {
	Update(g *game, t, dt time.Duration) fsm.State
}

type game struct {
	*fsm.FSM
	lastTime time.Duration
	display  *ggweb.Surface
	state    gameState
	// sound     map[string]sound.Interface
	// cardBg   *flyingCardBg
}

func newGame() *game {
	var g game
	g = game{
		display: ggweb.NewSurfaceFromID("main"),
		FSM: &fsm.FSM{
			Transitions: []*fsm.Transition{
				{fsm.InitialState, mainMenuState, func() {
					g.state = newMainMenuState()
				}},

				{mainMenuState, playState, func() {}},
				{mainMenuState, leaderboardState, func() {}},

				{playState, mainMenuState, func() {}},

				{leaderboardState, mainMenuState, func() {}},
			},
		},
	}
	g.display.SetSize(displayW, displayH)
	g.display.StyleColor(ggweb.Fill, color.Black)
	g.display.DrawRect(ggweb.Fill, g.display.Rect())
	ggweb.RegisterEvents(g.display)
	ggweb.DisableContextMenu = true

	e := g.Goto(mainMenuState)
	if e != nil {
		panic(e)
	}

	// globalState.cardBg = &flyingCardBg{
	// 	surf: ggweb.NewSurface(display.Width(), display.Height()),
	// }

	// globalState.sound = make(map[string]sound.Interface)
	// globalState.sound["btnHover"] = sound.New("hover.wav")
	// globalState.sound["btnHover"].SetVolume(0.8)

	return &g
}

func (g *game) mainLoop(t time.Duration) {
	dt := t - g.lastTime
	g.lastTime = t
	if e := g.Goto(g.state.Update(g, t, dt)); e != nil {
		ggweb.Error(e.Error())
	}
}
