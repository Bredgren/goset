package main

import (
	"time"

	"github.com/Bredgren/gogame"
	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/ui"
)

type mainMenuState struct {
	nextState GameState
	buttons   []*ui.BasicButton
	resumeBtn *ui.BasicButton
}

func (s *mainMenuState) Enter() {
	s.nextState = nil
	if len(s.buttons) == 0 {
		s.makeBtns()
	} else {
		for _, b := range s.buttons {
			b.State = btnIdle
		}
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
	globalState.cardBg.update(t, dt)
	s.draw()
}

func (s *mainMenuState) handleEvents() {
	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
		handleCommonEvents(evt)
		updateButtons(evt, s.buttons)
	}
}

func (s *mainMenuState) draw() {
	display := gogame.MainDisplay()
	display.Blit(globalState.cardBg.surf, 0, 0)

	// Draw tItle
	titleFont := gogame.Font{
		Size: 75,
	}
	titleStyle := gogame.TextStyle{
		Colorer:  gogame.White,
		Align:    gogame.TextAlignCenter,
		Baseline: gogame.TextBaselineMiddle,
	}
	display.DrawText("SET", display.Rect().CenterX(), 10+float64(titleFont.Size), &titleFont, &titleStyle)

	for _, btn := range s.buttons {
		btn.DrawTo(display)
	}
	display.Flip()
}

func (s *mainMenuState) makeBtns() {
	btnSpacing := 10.0

	s.resumeBtn = makeBtn("Resume", func() {
		gogame.Log("TODO: handle resume button")
		// s.nextState = globalState.playState
	})

	playBtn := makeBtn("Play", func() {
		s.nextState = globalState.playState
	})
	// playBtn.Rect.SetCenter(gogame.MainDisplay().Rect().Center())
	playBtn.Rect.SetMidLeft(40, gogame.MainDisplay().Rect().CenterY())
	s.buttons = append(s.buttons, playBtn)

	// s.resumeBtn.Rect.SetCenterX(playBtn.Rect.CenterX())
	s.resumeBtn.Rect.SetLeft(playBtn.Rect.Left())
	s.resumeBtn.Rect.SetBottom(playBtn.Rect.Top() - btnSpacing)

	leaderBtn := makeBtn("Leaderboard", func() {
		s.nextState = globalState.leaderboardState
	})
	// leaderBtn.Rect.SetCenterX(playBtn.Rect.CenterX())
	leaderBtn.Rect.SetLeft(playBtn.Rect.Left())
	leaderBtn.Rect.SetTop(playBtn.Rect.Bottom() + btnSpacing)
	s.buttons = append(s.buttons, leaderBtn)

	helpBtn := makeBtn("Help", func() {
		gogame.Log("TODO: handle help button")
		// s.nextState = globalState.helpState
	})
	// helpBtn.Rect.SetCenterX(playBtn.Rect.CenterX())
	helpBtn.Rect.SetLeft(playBtn.Rect.Left())
	helpBtn.Rect.SetTop(leaderBtn.Rect.Bottom() + btnSpacing)
	s.buttons = append(s.buttons, helpBtn)
}
