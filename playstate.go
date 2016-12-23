package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/fsm"
	"github.com/Bredgren/gogame/ggweb"
)

// import (
// 	"fmt"
// 	"math"
// 	"math/rand"
// 	"sort"
// 	"time"

// 	"github.com/Bredgren/gogame"
// 	"github.com/Bredgren/gogame/composite"
// 	"github.com/Bredgren/gogame/event"
// 	"github.com/Bredgren/gogame/geo"
// 	"github.com/Bredgren/gogame/key"
// 	"github.com/Bredgren/gogame/ui"
// )

type PlayState struct {
	SD        SaveData
	Paused    bool
	Buttons   []*Button
	NextState fsm.State
	// 	deck struct {
	// 		cards []card
	// 		rect  geo.Rect
	// 		hover bool
	// 	}
	// 	gameOver        bool
	// 	activeCards     []card
	// 	cardPos         []geo.Rect
	// 	hoverIndex      int
	// 	selectedCards   [3]int
	// 	cardRect        geo.Rect // Location and size of top-left card
	// 	cardAreaWidth   int
	// 	cardGap         float64
	// 	numCards        int // Target number of cards on table
	// 	maxActiveCards  int
	// 	score           int
	// 	errors          int
	// 	lastScoreChange time.Duration
	// 	lastErrorChange time.Duration
	// 	scalingCards    map[scaleAnim]bool
	// 	errorCards      map[int]time.Duration
	// 	playBtns        []*ui.BasicButton
	// 	pauseBtns       []*ui.BasicButton
}

func newPlayState(display *ggweb.Surface, sd SaveData) *PlayState {
	p := PlayState{
		SD:        sd,
		NextState: playState,
	}

	r := display.Rect()
	p.Buttons = []*Button{
		newTextButton("Save & Quit", 20, r.Left()+10, r.Bottom()-40, func() {
			p.NextState = mainMenuState
		}),
	}

	return &p
}

func (p *PlayState) Update(g *game, t, dt time.Duration) fsm.State {
	if !p.Paused {
		p.SD.PlayTime += dt
	}

	p.handleEvents()
	p.draw(g.display, t, dt)

	p.saveState()

	return p.NextState
}

func (p *PlayState) handleEvents() {
	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
		for _, b := range p.Buttons {
			b.handleEvent(evt)
		}
	}
}

func (p *PlayState) draw(display *ggweb.Surface, t, dt time.Duration) {
	display.StyleColor(ggweb.Fill, color.Black)
	display.DrawRect(ggweb.Fill, display.Rect())

	p.drawPlayTime(display)

	display.SetCursor(ggweb.CursorDefault)
	for _, b := range p.Buttons {
		b.drawTo(display)
		if b.State == buttonHover {
			display.SetCursor(ggweb.CursorPointer)
		}
	}
}

func (p *PlayState) drawPlayTime(display *ggweb.Surface) {
	timeString := fmt.Sprintf("%.0f", p.SD.PlayTime.Seconds())
	display.SetFont(&ggweb.Font{
		Size:   20,
		Family: ggweb.FontFamilyMonospace,
	})
	display.StyleColor(ggweb.Fill, color.White)
	display.SetTextAlign(ggweb.TextAlignLeft)
	display.SetTextBaseline(ggweb.TextBaselineTop)
	display.DrawText(ggweb.Fill, timeString, 10, 10)
}

func (p *PlayState) saveState() {
	if e := setSaveData(&p.SD); e != nil {
		ggweb.Log("Unable to save progress:", e.Error())
	}
}

// // var playPauseBtn ui.BasicButton
// // var saveQuitBtn ui.BasicButton

// // const (
// // 	gameBtnStatePlay ui.ButtonState = iota
// // 	gameBtnStatePlayHover
// // 	gameBtnStatePlaySelect
// // 	gameBtnStatePause
// // 	gameBtnStatePauseHover
// // 	gameBtnStatePauseSelect
// // )

// // func makeButtons() {
// // 	makePlayPauseBtn()
// // 	makeSaveQuitBtn()
// // }

// // func makePlayPauseBtn() {
// // 	display := gogame.MainDisplay()

// // 	r := geo.Rect{X: 10, Y: float64(display.Height()) - 10 - 50, W: 50, H: 50}

// // 	fill := gogame.FillStyle{
// // 		Colorer: gogame.Color{R: 0.5, G: 0.5, B: 0.5, A: 0.8},
// // 	}
// // 	outline := gogame.StrokeStyle{
// // 		Colorer: gogame.Color{R: 0.5, G: 0.5, B: 0.5, A: 0.8},
// // 		Width:   5,
// // 	}
// // 	playSurf := gogame.NewSurface(int(r.W), int(r.H))
// // 	playSurf.Fill(gogame.FillBlack)
// // 	playSurf.DrawLines([][2]float64{
// // 		{10, 10},
// // 		{r.W - 10, r.H / 2},
// // 		{10, r.H - 10},
// // 	}, &fill)
// // 	playHoverSurf := playSurf.Copy()
// // 	playHoverSurf.DrawRect(geo.Rect{W: r.W, H: r.H}, &outline)
// // 	playSelectSurf := playSurf.Copy()
// // 	playSelectSurf.DrawRect(geo.Rect{W: r.W, H: r.H}, &fill)

// // 	pauseSurf := gogame.NewSurface(int(r.W), int(r.H))
// // 	pauseSurf.Fill(gogame.FillBlack)
// // 	pauseSurf.DrawRect(geo.Rect{X: 10, Y: 10, W: r.W / 4, H: r.H - 20}, &fill)
// // 	pauseSurf.DrawRect(geo.Rect{X: r.W - 10 - r.W/4, Y: 10, W: r.W / 4, H: r.H - 20}, &fill)
// // 	pauseHoverSurf := pauseSurf.Copy()
// // 	pauseHoverSurf.DrawRect(geo.Rect{W: r.W, H: r.H}, &outline)
// // 	pauseSelectSurf := pauseSurf.Copy()
// // 	pauseSelectSurf.DrawRect(geo.Rect{W: r.W, H: r.H}, &fill)

// // 	playPauseBtn = ui.BasicButton{
// // 		Rect:        r,
// // 		DefaultSurf: pauseSurf,
// // 		StateSurfs: map[ui.ButtonState]gogame.Surface{
// // 			gameBtnStatePlay:        playSurf,
// // 			gameBtnStatePlayHover:   playHoverSurf,
// // 			gameBtnStatePlaySelect:  playSelectSurf,
// // 			gameBtnStatePause:       pauseSurf,
// // 			gameBtnStatePauseHover:  pauseHoverSurf,
// // 			gameBtnStatePauseSelect: pauseSelectSurf,
// // 		},
// // 		Select: func() {
// // 			state.paused = !state.paused
// // 		},
// // 		State: gameBtnStatePause,
// // 	}
// // }

// // func makeSaveQuitBtn() {
// // 	// TODO: ...
// // 	saveQuitBtn = ui.BasicButton{}
// // }

// func (s *playState) Enter() {
// 	s.gameOver = false
// 	s.playTime = 0
// 	s.score = 0
// 	s.hoverIndex = -1
// 	s.numCards = 12
// 	s.maxActiveCards = 20
// 	s.cardRect = geo.Rect{X: 100, Y: 50, W: 70, H: 100}
// 	s.cardAreaWidth = 3
// 	s.cardGap = 10

// 	s.deck = struct {
// 		cards []card
// 		rect  geo.Rect
// 		hover bool
// 	}{
// 		rect: geo.Rect{X: 10, Y: 50, W: 70, H: 100},
// 	}
// 	s.activeCards = make([]card, 0)
// 	s.selectedCards = [3]int{-1, -1, -1}
// 	s.scalingCards = make(map[scaleAnim]bool)
// 	s.errorCards = make(map[int]time.Duration)

// 	s.makeAndShuffleDeck()
// 	s.drawCards(s.numCards)

// }

// func (s *playState) Exit() {
// 	gogame.MainDisplay().SetCursor(gogame.CursorDefault)
// }

// func (s *playState) Update(t, dt time.Duration) {
// 	if !s.paused {
// 		s.playTime += dt
// 	}

// 	if s.gameOver {
// 		// TODO: handle game over
// 	}

// 	s.handleEvents()
// 	s.checkForSet(t)
// 	s.moveCards(dt)

// 	s.draw(t)
// }

// func (s *playState) handleEvents() {
// 	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
// 		handleCommonEvents(evt)
// 		switch evt.Type {
// 		case event.KeyDown:
// 			// data := evt.Data.(event.KeyData)
// 		case event.KeyUp:
// 			data := evt.Data.(event.KeyData)
// 			switch data.Key {
// 			case key.P, key.Escape:
// 				s.paused = !s.paused
// 				// if s.paused {
// 				// 	playPauseBtn.State = gameBtnStatePlaySelect
// 				// } else {
// 				// 	playPauseBtn.State = gameBtnStatePauseSelect
// 				// }
// 			}
// 		case event.MouseButtonDown:
// 			// data := evt.Data.(event.MouseData)
// 			// if playPauseBtn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
// 			// 	if s.paused {
// 			// 		playPauseBtn.State = gameBtnStatePlaySelect
// 			// 	} else {
// 			// 		playPauseBtn.State = gameBtnStatePauseSelect
// 			// 	}
// 			// }
// 		case event.MouseButtonUp:
// 			data := evt.Data.(event.MouseData)
// 			if data.Button == 0 && !s.paused {
// 				// Handle clicking the deck button
// 				if s.deck.hover {
// 					s.drawCards(3)
// 					if len(s.activeCards) >= s.maxActiveCards {
// 						s.deck.hover = false
// 					}
// 				}
// 				// Handle clicking a card
// 				if s.hoverIndex >= 0 {
// 					// First remove the card from the selected list if already selected
// 					removed := false
// 					for i := 0; i < len(s.selectedCards); i++ {
// 						if s.selectedCards[i] == s.hoverIndex {
// 							s.selectedCards[i] = -1
// 							removed = true
// 							break
// 						}
// 					}
// 					if !removed {
// 						// If not unselecting a card then select it if there is room in the list
// 						for i := 0; i < len(s.selectedCards); i++ {
// 							if s.selectedCards[i] < 0 {
// 								s.selectedCards[i] = s.hoverIndex
// 								break
// 							}
// 						}
// 					}
// 				}
// 			}

// 			// if data.Button == 0 && playPauseBtn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
// 			// 	playPauseBtn.Select()
// 			// 	if playPauseBtn.State == gameBtnStatePauseSelect {
// 			// 		playPauseBtn.State = gameBtnStatePlayHover
// 			// 	} else if playPauseBtn.State == gameBtnStatePlaySelect {
// 			// 		playPauseBtn.State = gameBtnStatePauseHover
// 			// 	}
// 			// }
// 		case event.MouseMotion:
// 			data := evt.Data.(event.MouseMotionData)
// 			if !s.paused {
// 				// Hover over card
// 				found := false
// 				for i := range s.activeCards {
// 					cr := s.getCardRect(i)
// 					if cr.CollidePoint(data.Pos.X, data.Pos.Y) {
// 						s.hoverIndex = i
// 						found = true
// 						break
// 					}
// 				}
// 				if !found {
// 					s.hoverIndex = -1
// 				}
// 				// Hover over deck if there are cards left in it and there is room on the table
// 				s.deck.hover = (len(s.activeCards) < 20 &&
// 					s.deck.rect.CollidePoint(data.Pos.X, data.Pos.Y) &&
// 					len(s.deck.cards) > 0)
// 			}

// 			// if playPauseBtn.Rect.CollidePoint(data.Pos.X, data.Pos.Y) {
// 			// 	if s.paused {
// 			// 		if data.Buttons[0] {
// 			// 			playPauseBtn.State = gameBtnStatePlaySelect
// 			// 		} else {
// 			// 			playPauseBtn.State = gameBtnStatePlayHover
// 			// 		}
// 			// 	} else {
// 			// 		if data.Buttons[0] {
// 			// 			playPauseBtn.State = gameBtnStatePauseSelect
// 			// 		} else {
// 			// 			playPauseBtn.State = gameBtnStatePauseHover
// 			// 		}
// 			// 	}
// 			// } else {
// 			// 	if s.paused {
// 			// 		playPauseBtn.State = gameBtnStatePlay
// 			// 	} else {
// 			// 		playPauseBtn.State = gameBtnStatePause
// 			// 	}
// 			// }
// 		} // End switch
// 	} // End for
// }

// func (s *playState) checkForSet(t time.Duration) {
// 	numSelected := 0
// 	for i := 0; i < len(s.selectedCards); i++ {
// 		if s.selectedCards[i] >= 0 {
// 			numSelected++
// 		}
// 	}

// 	if numSelected == 3 {
// 		c1 := s.activeCards[s.selectedCards[0]]
// 		c2 := s.activeCards[s.selectedCards[1]]
// 		c3 := s.activeCards[s.selectedCards[2]]
// 		if isSet(c1, c2, c3) {
// 			s.score++
// 			s.lastScoreChange = t
// 			for i := 0; i < len(s.selectedCards); i++ {
// 				s.scalingCards[scaleAnim{
// 					cardSurf:  s.activeCards[s.selectedCards[i]].surface(s.cardRect.W, s.cardRect.H),
// 					pos:       s.getCardRect(s.selectedCards[i]),
// 					startTime: t,
// 				}] = true
// 			}
// 			is := []int{s.selectedCards[0], s.selectedCards[1], s.selectedCards[2]}
// 			sort.Ints(is)
// 			s.activeCards = append(s.activeCards[:is[2]], s.activeCards[is[2]+1:]...)
// 			s.activeCards = append(s.activeCards[:is[1]], s.activeCards[is[1]+1:]...)
// 			s.activeCards = append(s.activeCards[:is[0]], s.activeCards[is[0]+1:]...)
// 			s.cardPos = append(s.cardPos[:is[2]], s.cardPos[is[2]+1:]...)
// 			s.cardPos = append(s.cardPos[:is[1]], s.cardPos[is[1]+1:]...)
// 			s.cardPos = append(s.cardPos[:is[0]], s.cardPos[is[0]+1:]...)
// 			for len(s.activeCards) < s.numCards && len(s.deck.cards) > 0 {
// 				s.activeCards = append(s.activeCards, s.deck.cards[0])
// 				s.cardPos = append(s.cardPos, s.deck.rect)
// 				s.deck.cards = s.deck.cards[1:]
// 			}

// 			// Check for end game
// 			if len(s.deck.cards) == 0 && !s.setOnField() {
// 				s.gameOver = true
// 			}
// 		} else {
// 			s.errors++
// 			s.lastErrorChange = t
// 			for _, i := range s.selectedCards {
// 				s.errorCards[i] = t + time.Duration(50*time.Millisecond)
// 			}
// 		}
// 		for i := 0; i < len(s.selectedCards); i++ {
// 			s.selectedCards[i] = -1
// 		}
// 		s.hoverIndex = -1
// 	}
// }

// func (s *playState) moveCards(dt time.Duration) {
// 	for i := 0; i < len(s.activeCards); i++ {
// 		target := s.getCardRect(i)
// 		current := s.cardPos[i]
// 		dx := target.X - current.X
// 		dy := target.Y - current.Y
// 		dist := math.Sqrt(dx*dx + dy*dy)
// 		dirX := dx / dist
// 		dirY := dy / dist
// 		speed := math.Log(dist) * 250 * dt.Seconds()
// 		moveX := dirX * speed
// 		moveY := dirY * speed
// 		moveDist := math.Sqrt(moveX*moveX + moveY*moveY)
// 		if moveDist >= dist || dist <= 1 {
// 			// About to overshoot so just stop at target
// 			s.cardPos[i].X, s.cardPos[i].Y = target.X, target.Y
// 		} else {
// 			s.cardPos[i].X += moveX
// 			s.cardPos[i].Y += moveY
// 		}
// 	}
// }

// func (s *playState) draw(t time.Duration) {
// 	display := gogame.MainDisplay()
// 	display.Fill(gogame.FillBlack)
// 	s.drawPlayTime(display)
// 	s.drawScore(display, t)
// 	if s.paused {
// 		font := gogame.Font{Size: 50}
// 		style := gogame.TextStyle{
// 			Colorer:  gogame.White,
// 			Type:     gogame.Fill,
// 			Align:    gogame.TextAlignCenter,
// 			Baseline: gogame.TextBaselineHanging,
// 		}
// 		display.DrawText("Paused", float64(display.Width()/2), s.cardRect.Y, &font, &style)
// 		// saveQuitBtn.DrawTo(display)
// 		// font.Size = 20
// 		// style.Align = gogame.TextAlignLeft
// 		// style.Baseline = gogame.TextBaselineAlphabetic
// 		// display.DrawText("A set is any 3 cards where each of the 4 properties (count, fill, color, shape) are all the same or all different.",
// 		// 	50, s.cardRect.Y+70, &font, &style)
// 	} else {
// 		s.drawHoverHighlight(display, t)
// 		s.drawActiveCards(display)
// 		for i, endTime := range s.errorCards {
// 			display.DrawRect(s.getCardRect(i), &gogame.FillStyle{Colorer: gogame.Color{R: 1.0, A: 0.6}})
// 			if t > endTime {
// 				delete(s.errorCards, i)
// 			}
// 		}
// 		display.DrawRect(s.deck.rect, &gogame.FillStyle{Colorer: gogame.Color{R: 0.9, G: 0.4, B: 1.0, A: 1.0}})
// 		display.DrawText(fmt.Sprintf("%d", len(s.deck.cards)), s.deck.rect.CenterX(),
// 			s.deck.rect.CenterY(),
// 			&gogame.Font{
// 				Size:   20,
// 				Family: gogame.FontFamilyMonospace,
// 			}, &gogame.TextStyle{
// 				Colorer:  gogame.White,
// 				Align:    gogame.TextAlignCenter,
// 				Baseline: gogame.TextBaselineMiddle,
// 			})
// 		if s.deck.hover {
// 			display.DrawRect(s.deck.rect.Inflate(6, 6), &gogame.StrokeStyle{
// 				Colorer: gogame.Color{R: 1.0, G: 1.0, B: 0.0, A: 0.2*math.Sin(4*t.Seconds()) + 0.8},
// 				Width:   3,
// 				Join:    gogame.LineJoinRound,
// 			})
// 			display.DrawText("+3 cards", s.deck.rect.CenterX(), s.deck.rect.Bottom()+5,
// 				&gogame.Font{
// 					Size: 15,
// 				}, &gogame.TextStyle{
// 					Colorer:  gogame.White,
// 					Align:    gogame.TextAlignCenter,
// 					Baseline: gogame.TextBaselineTop,
// 				})
// 		}

// 		for sc, ok := range s.scalingCards {
// 			if !ok {
// 				continue
// 			}
// 			surf, done := sc.surface(t)
// 			r := surf.Rect()
// 			r.SetCenter(sc.pos.Center())
// 			display.Blit(surf, r.X, r.Y)
// 			if done {
// 				delete(s.scalingCards, sc)
// 			}
// 		}
// 	}

// 	// playPauseBtn.DrawTo(display)

// 	display.Flip()
// }

// func (s *playState) drawPlayTime(display gogame.Surface) {
// 	timeString := fmt.Sprintf("%.0f", s.playTime.Seconds())
// 	font := gogame.Font{
// 		Size:   20,
// 		Family: gogame.FontFamilyMonospace,
// 	}
// 	style := gogame.TextStyle{
// 		Colorer:  gogame.White,
// 		Type:     gogame.Fill,
// 		Align:    gogame.TextAlignLeft,
// 		Baseline: gogame.TextBaselineTop,
// 	}
// 	display.DrawText(timeString, 10, 10, &font, &style)
// }

// func (s *playState) drawScore(display gogame.Surface, t time.Duration) {
// 	scoreString := fmt.Sprintf("%d", s.score)
// 	errorString := fmt.Sprintf("-%d", s.errors)
// 	baseSize := 30
// 	ds := 15
// 	fadeTime := time.Duration(250 * time.Millisecond)
// 	font := gogame.Font{
// 		Size:   baseSize,
// 		Family: gogame.FontFamilyMonospace,
// 	}
// 	style := gogame.TextStyle{
// 		Colorer:  gogame.White,
// 		Type:     gogame.Fill,
// 		Align:    gogame.TextAlignCenter,
// 		Baseline: gogame.TextBaselineHanging,
// 	}
// 	font.Size = baseSize + ds - int(float64(ds)*math.Min(float64(t-s.lastScoreChange)/float64(fadeTime), 1.0))
// 	display.DrawText(scoreString, float64(display.Width()/2), 10, &font, &style)
// 	if s.errors > 0 {
// 		baseSize = 20
// 		style.Colorer = gogame.Color{R: 1.0, G: 0.9, B: 0.9, A: 1.0}
// 		font.Size = baseSize + ds - int(float64(ds)*math.Min(float64(t-s.lastErrorChange)/float64(fadeTime), 1.0))
// 		display.DrawText(errorString, float64(display.Width()/2)+50, 10, &font, &style)
// 	}
// }

// func (s *playState) drawHoverHighlight(display gogame.Surface, t time.Duration) {
// 	if s.hoverIndex < 0 {
// 		return
// 	}
// 	cr := s.getCardRect(s.hoverIndex)
// 	display.DrawRect(cr.Inflate(6, 6), &gogame.StrokeStyle{
// 		Colorer: gogame.Color{R: 1.0, G: 1.0, B: 0.0, A: 0.2*math.Sin(2*t.Seconds()) + 0.8},
// 		Width:   3,
// 		Join:    gogame.LineJoinRound,
// 	})
// }

// func (s *playState) drawActiveCards(display gogame.Surface) {
// 	selectSurf := gogame.NewSurface(int(s.cardRect.W), int(s.cardRect.H))
// 	selectSurf.Fill(&gogame.FillStyle{Colorer: gogame.Color{A: 0.3}})
// 	for i, card := range s.activeCards {
// 		cr := s.cardPos[i]
// 		display.Blit(card.surface(s.cardRect.W, s.cardRect.H), cr.X, cr.Y)
// 		// Slightly darken the card if selected
// 		for _, j := range s.selectedCards {
// 			if j == i {
// 				display.Blit(selectSurf, cr.X, cr.Y)
// 				break
// 			}
// 		}
// 	}
// }

// func (s *playState) makeAndShuffleDeck() {
// 	tmpDeck := []card{}
// 	for _, n := range []count{one, two, three} {
// 		for _, f := range []fill{empty, solid, line} {
// 			for _, c := range []color{red, green, purple} {
// 				for _, s := range []shape{oval, diamond, tilde} {
// 					tmpDeck = append(tmpDeck, card{count: n, fill: f, color: c, shape: s})
// 				}
// 			}
// 		}
// 	}
// 	s.deck.cards = make([]card, len(tmpDeck))
// 	order := rand.Perm(len(tmpDeck))
// 	for i, pos := range order {
// 		s.deck.cards[i] = tmpDeck[pos]
// 	}
// }

// func (s *playState) drawCards(count int) {
// 	for i := 0; i < count; i++ {
// 		s.activeCards = append(s.activeCards, s.deck.cards[i])
// 		s.cardPos = append(s.cardPos, s.deck.rect)
// 	}
// 	s.deck.cards = s.deck.cards[count:]
// }

// func (s *playState) setOnField() bool {
// 	for i1, c1 := range s.activeCards[:len(s.activeCards)-2] {
// 		for i2, c2 := range s.activeCards[i1+1 : len(s.activeCards)-1] {
// 			for _, c3 := range s.activeCards[i2+1 : len(s.activeCards)] {
// 				if isSet(c1, c2, c3) {
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

// // i is the index in s.activeCards
// func (s *playState) getCardRect(i int) geo.Rect {
// 	r := s.cardRect.Copy()
// 	r.X += (r.W + s.cardGap) * float64(i/s.cardAreaWidth)
// 	r.Y += (r.H + s.cardGap) * float64(i%s.cardAreaWidth)
// 	return r
// }

// type scaleAnim struct {
// 	cardSurf  gogame.Surface
// 	pos       geo.Rect
// 	startTime time.Duration
// }

// func (a *scaleAnim) surface(t time.Duration) (s gogame.Surface, done bool) {
// 	curTime := t - a.startTime
// 	animTime := time.Duration(400 * time.Millisecond)
// 	maxScale := 1.0 // 100% larger
// 	percent := float64(curTime) / float64(animTime)
// 	scale := 1.0 + (percent * maxScale)
// 	scaled := a.cardSurf.Scaled(scale, scale)
// 	fade := gogame.NewSurface(scaled.Width(), scaled.Height())
// 	fade.Fill(&gogame.FillStyle{Colorer: gogame.Color{A: 1.0 - percent}})
// 	scaled.BlitComp(fade, 0, 0, composite.DestinationIn)
// 	return scaled, curTime > animTime
// }
