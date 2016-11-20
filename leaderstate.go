package main

// import (
// 	"sort"
// 	"time"

// 	"github.com/Bredgren/gogame"
// 	"github.com/Bredgren/gogame/event"
// 	"github.com/Bredgren/gogame/ui"
// )

// type leaderboardState struct {
// 	leaving bool
// 	buttons []*ui.BasicButton
// 	data    []leaderboardEntry
// }

// func (s *leaderboardState) Enter() {
// 	s.leaving = false
// 	s.data = getLeaderboardData()
// 	sort.Sort(entryList(s.data))

// 	if len(s.buttons) == 0 {
// 		s.makeBtns()
// 	} else {
// 		for _, b := range s.buttons {
// 			b.State = btnIdle
// 		}
// 	}
// }

// func (s *leaderboardState) Exit() {
// 	gogame.MainDisplay().SetCursor(gogame.CursorDefault)
// }

// func (s *leaderboardState) Update(t, dt time.Duration) {
// 	if s.leaving {
// 		// TODO: leaving animation
// 		globalState.gameStateMgr.Goto(globalState.mainMenuState)
// 		return
// 	}
// 	s.handleEvents()
// 	s.draw()
// }

// func (s *leaderboardState) handleEvents() {
// 	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
// 		handleCommonEvents(evt)
// 		updateButtons(evt, s.buttons)
// 	}
// }

// func (s *leaderboardState) draw() {
// 	display := gogame.MainDisplay()
// 	display.Fill(gogame.FillBlack)
// 	for _, btn := range s.buttons {
// 		btn.DrawTo(display)
// 	}
// 	display.Flip()
// }

// type entryList []leaderboardEntry

// func (l entryList) Len() int {
// 	return len(l)
// }

// func (l entryList) Less(i, j int) bool {
// 	if l[i].numSets != l[j].numSets {
// 		return l[i].numSets > l[j].numSets
// 	}
// 	if l[i].numErrors != l[j].numErrors {
// 		return l[i].numErrors < l[j].numErrors
// 	}
// 	return l[i].time < l[j].time
// }

// func (l entryList) Swap(i, j int) {
// 	l[i], l[j] = l[j], l[i]
// }

// func (s *leaderboardState) makeBtns() {
// 	backBtn := backArrowBtn(func() {
// 		s.leaving = true
// 	})
// 	backBtn.Rect.SetBottomLeft(10, gogame.MainDisplay().Rect().H-10)
// 	s.buttons = append(s.buttons, backBtn)
// }
