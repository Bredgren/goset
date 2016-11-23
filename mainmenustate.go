package main

import (
	"image/color"
	"time"

	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/fsm"
	"github.com/Bredgren/gogame/ggweb"
)

// import (
// 	"time"

// 	"github.com/Bredgren/gogame"
// 	"github.com/Bredgren/gogame/event"
// 	"github.com/Bredgren/gogame/ui"
// )

type mainMenu struct {
	// nextState GameState
	buttons []*Button
	// resumeBtn *ui.BasicButton
}

func newMainMenuState() gameState {
	buttons := []*Button{
		newTextButton("Play", 10, 10, func() {
			ggweb.Log("Play")
		}),
		newTextButton("Leaderboard", 10, 50, func() {
			ggweb.Log("Leaderboard")
		}),
		newTextButton("Help", 10, 90, func() {
			ggweb.Log("Help")
		}),
	}
	return &mainMenu{
		buttons: buttons,
	}
}

// func (s *mainMenuState) Enter() {
// 	s.nextState = nil
// 	if len(s.buttons) == 0 {
// 		s.makeBtns()
// 	} else {
// 		for _, b := range s.buttons {
// 			b.State = btnIdle
// 		}
// 	}

// 	// TODO: actually check for saved game
// 	savedGame := true
// 	ri := -1
// 	for i, b := range s.buttons {
// 		if b == s.resumeBtn {
// 			ri = i
// 			break
// 		}
// 	}
// 	if savedGame && ri < 0 {
// 		s.buttons = append(s.buttons, s.resumeBtn)
// 	} else if !savedGame && ri >= 0 {
// 		s.buttons = append(s.buttons[:ri], s.buttons[ri+1:]...)
// 	}
// }

// func (s *mainMenuState) Exit() {
// 	gogame.MainDisplay().SetCursor(gogame.CursorDefault)
// }

func (s *mainMenu) Update(g *game, t, dt time.Duration) fsm.State {
	// 	if s.nextState != nil {
	// 		// TODO: leaving animation
	// 		globalState.gameStateMgr.Goto(s.nextState)
	// 	}
	s.handleEvents()
	// 	globalState.cardBg.update(t, dt)
	s.draw(g.display)
	return mainMenuState
}

func (s *mainMenu) handleEvents() {
	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
		// handleCommonEvents(evt)
		// updateButtons(evt, s.buttons)
		for _, b := range s.buttons {
			b.handleEvent(evt)
		}
	}
}

func (s *mainMenu) draw(display *ggweb.Surface) {
	// 	display := gogame.MainDisplay()
	// 	display.Blit(globalState.cardBg.surf, 0, 0)

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

	display.StyleColor(ggweb.Fill, color.Black)
	display.DrawRect(ggweb.Fill, display.Rect())
	display.SetCursor(ggweb.CursorDefault)
	for _, b := range s.buttons {
		b.drawTo(display)
		if b.State == buttonHover {
			display.SetCursor(ggweb.CursorPointer)
		}
	}
	// 	display.Flip()
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
