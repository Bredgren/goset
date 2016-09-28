// TODO:
//   - saving/restoring game
//   - main menu
//   - "tutorial"
//   - leaderboard
package main

import (
	"math/rand"
	"time"

	"github.com/Bredgren/gogame"
)

func main() {
	gogame.Ready(func() {
		setup()
		gogame.SetMainLoop(mainLoop)
	})
}

const (
	displayW, displayH = 680, 500
)

var globalState = struct {
	lastTime         time.Duration
	gameStateMgr     StateMgr
	mainMenuState    *mainMenuState
	playState        *playState
	leaderboardState *leaderboardState
}{
	mainMenuState:    &mainMenuState{},
	playState:        &playState{},
	leaderboardState: &leaderboardState{},
}

func setup() {
	display := gogame.MainDisplay()
	display.SetMode(displayW, displayH)
	display.Fill(gogame.FillBlack)
	display.Flip()
	rand.Seed(time.Now().UnixNano())

	globalState.gameStateMgr.Goto(globalState.mainMenuState)
}

func mainLoop(t time.Duration) {
	dt := t - globalState.lastTime
	globalState.lastTime = t
	globalState.gameStateMgr.Current().Update(t, dt)
}
