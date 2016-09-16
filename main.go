package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/Bredgren/gogame"
	"github.com/Bredgren/gogame/composite"
	"github.com/Bredgren/gogame/event"
	"github.com/Bredgren/gogame/geo"
	"github.com/Bredgren/gogame/key"
)

func main() {
	ready := gogame.Ready()
	go func() {
		<-ready
		setup()
		gogame.SetMainLoop(mainLoop)
	}()
}

const (
	displayW, displayH = 680, 500
)

var state = struct {
	deck struct {
		cards []card
		rect  geo.Rect
		hover bool
	}
	gameState     gameState
	paused        bool
	lastTime      time.Duration
	playTime      time.Duration
	activeCards   []card
	hoverIndex    int
	selectedCards [3]int
	cardRect      geo.Rect
	cardAreaWidth int
	cardGap       float64
	numCards      int
	score         int
}{
	deck: struct {
		cards []card
		rect  geo.Rect
		hover bool
	}{
		rect: geo.Rect{X: 10, Y: 50, W: 70, H: 100},
	},
	gameState:     menuState,
	activeCards:   make([]card, 0),
	hoverIndex:    -1,
	selectedCards: [3]int{-1, -1, -1},
	cardRect:      geo.Rect{X: 100, Y: 50, W: 70, H: 100}, // Location and size of top-left card
	cardAreaWidth: 3,
	cardGap:       10,
	numCards:      12, // Target number of cards on table
}

func setup() {
	display := gogame.MainDisplay()
	display.SetMode(displayW, displayH)
	display.Fill(gogame.FillBlack)
	display.Flip()
	rand.Seed(time.Now().UnixNano())
}

func makeAndShuffleDeck() {
	tmpDeck := []card{}
	for _, n := range []count{one, two, three} {
		for _, f := range []fill{empty, solid, line} {
			for _, c := range []color{red, green, purple} {
				for _, s := range []shape{oval, dimond, tilde} {
					tmpDeck = append(tmpDeck, card{count: n, fill: f, color: c, shape: s})
				}
			}
		}
	}
	state.deck.cards = make([]card, len(tmpDeck))
	order := rand.Perm(len(tmpDeck))
	for i, pos := range order {
		state.deck.cards[i] = tmpDeck[pos]
	}
}

func mainLoop(t time.Duration) {
	dt := t - state.lastTime
	state.lastTime = t

	switch state.gameState {
	case menuState:
		gotoPlayState()
	case playState:
		handlePlayStateLoop(t, dt)
	case gameOverState:
	case leaderboardState:
	}
}

func gotoPlayState() {
	state.playTime = 0
	state.gameState = playState
	makeAndShuffleDeck()
	for i := 0; i < state.numCards; i++ {
		state.activeCards = append(state.activeCards, state.deck.cards[i])
	}
	state.deck.cards = state.deck.cards[len(state.activeCards):]
	state.score = 0
}

func handlePlayStateLoop(t time.Duration, dt time.Duration) {
	if !state.paused {
		state.playTime += dt
	}

	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
		switch evt.Type {
		case event.Quit:
			gogame.UnsetMainLoop()
		case event.VideoResize:
			// data := evt.Data.(event.ResizeData)
		case event.KeyDown:
			// data := evt.Data.(event.KeyData)
		case event.KeyUp:
			data := evt.Data.(event.KeyData)
			switch data.Key {
			case key.P:
				state.paused = !state.paused
				if !state.paused {
					state.lastTime = t
				}
			}
		case event.MouseButtonDown:
			// data := evt.Data.(event.MouseData)
		case event.MouseButtonUp:
			data := evt.Data.(event.MouseData)
			if data.Button == 0 {
				if state.deck.hover {
					numCards := 3
					for i := 0; i < numCards; i++ {
						state.activeCards = append(state.activeCards, state.deck.cards[i])
					}
					state.deck.cards = state.deck.cards[numCards:]
					if len(state.activeCards) >= 20 {
						state.deck.hover = false
					}
				}
				if state.hoverIndex >= 0 {
					// First removed the card from the selected list if already selected
					removed := false
					for i := 0; i < len(state.selectedCards); i++ {
						if state.selectedCards[i] == state.hoverIndex {
							state.selectedCards[i] = -1
							removed = true
							break
						}
					}
					if !removed {
						// If not unselecting a card then select it if there is room in the list
						for i := 0; i < len(state.selectedCards); i++ {
							if state.selectedCards[i] < 0 {
								state.selectedCards[i] = state.hoverIndex
								break
							}
						}
					}
				}
			}
		case event.MouseMotion:
			data := evt.Data.(event.MouseMotionData)
			found := false
			for i := range state.activeCards {
				cr := getCardRect(i)
				if cr.CollidePoint(data.Pos.X, data.Pos.Y) {
					state.hoverIndex = i
					found = true
					break
				}
			}
			if !found {
				state.hoverIndex = -1
			}
			state.deck.hover = len(state.activeCards) < 20 && state.deck.rect.CollidePoint(data.Pos.X, data.Pos.Y)
		}
	}

	numSelected := 0
	for i := 0; i < len(state.selectedCards); i++ {
		if state.selectedCards[i] >= 0 {
			numSelected++
		}
	}

	if numSelected == 3 {
		// TODO:
		//   Animate cards out and new cards in
		//   Check for end game
		c1 := state.activeCards[state.selectedCards[0]]
		c2 := state.activeCards[state.selectedCards[1]]
		c3 := state.activeCards[state.selectedCards[2]]
		if ((c1.count == c2.count && c1.count == c3.count) || (c1.count != c2.count && c1.count != c3.count && c2.count != c3.count)) &&
			((c1.fill == c2.fill && c1.fill == c3.fill) || (c1.fill != c2.fill && c1.fill != c3.fill && c2.fill != c3.fill)) &&
			((c1.color == c2.color && c1.color == c3.color) || (c1.color != c2.color && c1.color != c3.color && c2.color != c3.color)) &&
			((c1.shape == c2.shape && c1.shape == c3.shape) || (c1.shape != c2.shape && c1.shape != c3.shape && c2.shape != c3.shape)) {
			state.score++
			is := []int{state.selectedCards[0], state.selectedCards[1], state.selectedCards[2]}
			sort.Ints(is)
			state.activeCards = append(state.activeCards[:is[2]], state.activeCards[is[2]+1:]...)
			state.activeCards = append(state.activeCards[:is[1]], state.activeCards[is[1]+1:]...)
			state.activeCards = append(state.activeCards[:is[0]], state.activeCards[is[0]+1:]...)
			for len(state.activeCards) < 12 && len(state.deck.cards) > 0 {
				state.activeCards = append(state.activeCards, state.deck.cards[0])
				state.deck.cards = state.deck.cards[1:]
			}
		}
		for i := 0; i < len(state.selectedCards); i++ {
			state.selectedCards[i] = -1
		}
	}

	display := gogame.MainDisplay()
	display.Fill(gogame.FillBlack)
	drawPlayTime(display)
	drawScore(display)
	if state.paused {
		font := gogame.Font{Size: 50}
		style := gogame.TextStyle{
			Colorer:  gogame.White,
			Type:     gogame.Fill,
			Align:    gogame.TextAlignCenter,
			Baseline: gogame.TextBaselineHanging,
		}
		display.DrawText("Paused", float64(display.Width()/2), state.cardRect.Y, &font, &style)
	} else {
		drawHoverHighlight(display, t)
		drawActiveCards(display)
		display.DrawRect(state.deck.rect, &gogame.FillStyle{Colorer: gogame.Color{R: 0.9, G: 0.4, B: 1.0, A: 1.0}})
		display.DrawText(fmt.Sprintf("%d", len(state.deck.cards)), state.deck.rect.CenterX(),
			state.deck.rect.CenterY(),
			&gogame.Font{
				Size:   20,
				Family: gogame.FontFamilyMonospace,
			}, &gogame.TextStyle{
				Colorer:  gogame.White,
				Align:    gogame.TextAlignCenter,
				Baseline: gogame.TextBaselineMiddle,
			})
		if state.deck.hover {
			display.DrawRect(state.deck.rect.Inflate(6, 6), &gogame.StrokeStyle{
				Colorer: gogame.Color{R: 1.0, G: 1.0, B: 0.0, A: 0.2*math.Sin(4*t.Seconds()) + 0.8},
				Width:   3,
				Join:    gogame.LineJoinRound,
			})
			display.DrawText("+3 cards", state.deck.rect.CenterX(), state.deck.rect.Bottom()+5,
				&gogame.Font{
					Size: 15,
				}, &gogame.TextStyle{
					Colorer:  gogame.White,
					Align:    gogame.TextAlignCenter,
					Baseline: gogame.TextBaselineTop,
				})
		}
	}

	display.Flip()
}

// i is the index in state.activeCards
func getCardRect(i int) geo.Rect {
	r := state.cardRect.Copy()
	r.X += (r.W + state.cardGap) * float64(i/state.cardAreaWidth)
	r.Y += (r.H + state.cardGap) * float64(i%state.cardAreaWidth)
	return r
}

func drawPlayTime(display gogame.Surface) {
	// TODO: option to hide the time
	timeString := fmt.Sprintf("%.0f", state.playTime.Seconds())
	font := gogame.Font{
		Size:   20,
		Family: gogame.FontFamilyMonospace,
	}
	style := gogame.TextStyle{
		Colorer:  gogame.White,
		Type:     gogame.Fill,
		Align:    gogame.TextAlignLeft,
		Baseline: gogame.TextBaselineTop,
	}
	display.DrawText(timeString, 10, 10, &font, &style)
}

func drawScore(display gogame.Surface) {
	scoreString := fmt.Sprintf("%d", state.score)
	font := gogame.Font{
		Size:   30,
		Family: gogame.FontFamilyMonospace,
	}
	style := gogame.TextStyle{
		Colorer:  gogame.White,
		Type:     gogame.Fill,
		Align:    gogame.TextAlignCenter,
		Baseline: gogame.TextBaselineHanging,
	}
	display.DrawText(scoreString, float64(display.Width()/2), 10, &font, &style)
}

func drawHoverHighlight(display gogame.Surface, t time.Duration) {
	if state.hoverIndex < 0 {
		return
	}
	cr := getCardRect(state.hoverIndex)
	display.DrawRect(cr.Inflate(6, 6), &gogame.StrokeStyle{
		Colorer: gogame.Color{R: 1.0, G: 1.0, B: 0.0, A: 0.2*math.Sin(2*t.Seconds()) + 0.8},
		Width:   3,
		Join:    gogame.LineJoinRound,
	})
}

func drawActiveCards(display gogame.Surface) {
	selectSurf := gogame.NewSurface(int(state.cardRect.W), int(state.cardRect.H))
	selectSurf.Fill(&gogame.FillStyle{Colorer: gogame.Color{A: 0.2}})
	for i, card := range state.activeCards {
		cr := getCardRect(i)
		display.Blit(card.surface(state.cardRect.W, state.cardRect.H), cr.X, cr.Y)
		// Slightly darken the card if selected
		for _, j := range state.selectedCards {
			if j == i {
				display.Blit(selectSurf, cr.X, cr.Y)
				break
			}
		}
	}
}

func makeOval(w, h float64, c color, f fill) gogame.Surface {
	color := gogame.Color(c)
	s := gogame.NewSurface(int(w), int(h))
	if f == solid {
		style := &gogame.FillStyle{
			Colorer: color,
		}
		s.DrawRect(geo.Rect{X: h / 2, Y: 0, W: w - h, H: h}, style)
		s.DrawArc(geo.Rect{X: 0, Y: 0, W: h, H: h}, math.Pi/2, 3*math.Pi/2, style)
		s.DrawArc(geo.Rect{X: w - h, Y: 0, W: h, H: h}, -math.Pi/2, math.Pi/2, style)
	} else {
		style := &gogame.StrokeStyle{
			Colorer: color,
			Width:   0.03 * float64(w),
		}
		if f == line {
			s.Blit(makeLines(w, h, style), 0, 0)
			mask := gogame.NewSurface(int(w), int(h))
			mask.DrawRect(geo.Rect{X: h / 2, Y: 0, W: w - h, H: h}, gogame.FillWhite)
			mask.DrawArc(geo.Rect{X: 0, Y: 0, W: h, H: h}, math.Pi/2, 3*math.Pi/2, gogame.FillWhite)
			mask.DrawArc(geo.Rect{X: w - h, Y: 0, W: h, H: h}, -math.Pi/2, math.Pi/2, gogame.FillWhite)
			s.BlitComp(mask, 0, 0, composite.DestinationIn)
		}
		s.DrawLine(h/2, 1, w-h/2, 1, style)
		s.DrawLine(h/2, h-1, w-h/2, h-1, style)
		s.DrawArc(geo.Rect{X: 1, Y: 1, W: h, H: h - 2}, math.Pi/2, 3*math.Pi/2, style)
		s.DrawArc(geo.Rect{X: w - h - 1, Y: 1, W: h, H: h - 2}, -math.Pi/2, math.Pi/2, style)
	}
	return s
}

func makeDimond(w, h float64, c color, f fill) gogame.Surface {
	color := gogame.Color(c)
	s := gogame.NewSurface(int(w), int(h))
	points := [][2]float64{
		{w / 2, 0},
		{w, h / 2},
		{w / 2, h},
		{0, h / 2},
		{w / 2, 0},
	}
	if f == solid {
		s.DrawLines(points, &gogame.FillStyle{Colorer: color})
	} else {
		style := &gogame.StrokeStyle{
			Colorer: color,
			Width:   0.03 * w,
		}
		if f == line {
			s.Blit(makeLines(w, h, style), 0, 0)
			mask := gogame.NewSurface(int(w), int(h))
			mask.DrawLines(points, gogame.FillWhite)
			s.BlitComp(mask, 0, 0, composite.DestinationIn)
		}
		points[0][1] += style.Width / 2
		points[1][0] -= style.Width / 2
		points[2][1] -= style.Width / 2
		points[3][0] += style.Width / 2
		points[4][1] += style.Width / 2
		s.DrawLines(points, style)
	}
	return s
}

func makeTilde(w, h float64, c color, f fill) gogame.Surface {
	color := gogame.Color(c)
	s := gogame.NewSurface(int(w), int(h))
	p1x, p1y := w*0.04, h*0.5
	p2x, p2y := w*0.77, h*0.125
	cp1x, cp1y := w*0.2, h*0.5-h
	cp1x2, cp1y2 := p2x-w*0.2, h*0.5
	cp2x, cp2y := p2x+w*0.1, 0.0
	cp2x2, cp2y2 := w+w*0.01, 0.0
	points := [][2]float64{
		{p1x, p1y}, {cp1x, cp1y}, {cp1x2, cp1y2},
		{p2x, p2y}, {cp2x, cp2y}, {cp2x2, cp2y2},
		{w - p1x, h - p1y}, {w - cp1x, h - cp1y}, {w - cp1x2, h - cp1y2},
		{w - p2x, h - p2y}, {w - cp2x, h - cp2y}, {w - cp2x2, h - cp2y2},
		{p1x, p1y},
	}
	if f == solid {
		s.DrawBezierCurves(points, &gogame.FillStyle{Colorer: color})
	} else {
		style := &gogame.StrokeStyle{
			Colorer: color,
			Width:   0.03 * float64(w),
		}
		if f == line {
			s.Blit(makeLines(w, h, style), 0, 0)
			mask := gogame.NewSurface(int(w), int(h))
			mask.DrawBezierCurves(points, gogame.FillWhite)
			s.BlitComp(mask, 0, 0, composite.DestinationIn)
		}
		s.DrawBezierCurves(points, style)
	}
	return s
}

func makeLines(w, h float64, style *gogame.StrokeStyle) gogame.Surface {
	s := gogame.NewSurface(int(w), int(h))
	wf := float64(w)
	for x := 0.1 * wf; x < wf; x += 0.1 * wf {
		xx := x
		if int(math.Ceil(style.Width))%2 != 0 {
			xx = math.Floor(x) + 0.5
		}
		s.DrawLine(xx, 0, xx, h, style)
	}
	return s
}

type count int

const (
	one count = iota
	two
	three
)

type shape int

const (
	oval shape = iota
	dimond
	tilde
)

type color gogame.Color

var (
	red    = color(gogame.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0})
	green  = color(gogame.Color{R: 0.0, G: 0.9, B: 0.0, A: 1.0})
	purple = color(gogame.Color{R: 1.0, G: 0.0, B: 1.0, A: 1.0})
)

type fill int

const (
	empty fill = iota
	solid
	line
)

type card struct {
	count count
	fill  fill
	color color
	shape shape
}

func (c *card) surface(w, h float64) gogame.Surface {
	surf := gogame.NewSurface(int(w), int(h))
	surf.Fill(gogame.FillWhite)
	var shapeSurf gogame.Surface
	shapeW := w * 0.75
	shapeH := shapeW / 2
	switch c.shape {
	case oval:
		shapeSurf = makeOval(shapeW, shapeH, c.color, c.fill)
	case dimond:
		shapeSurf = makeDimond(shapeW, shapeH, c.color, c.fill)
	case tilde:
		shapeSurf = makeTilde(shapeW, shapeH, c.color, c.fill)
	}
	centerRect := geo.Rect{W: shapeW, H: shapeH}
	centerRect.SetCenterX(w / 2)
	centerRect.SetCenterY(h / 2)
	switch c.count {
	case one:
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y)
	case two:
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y-shapeH*0.75)
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y+shapeH*0.75)
	case three:
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y-shapeH*1.25)
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y)
		surf.Blit(shapeSurf, centerRect.X, centerRect.Y+shapeH*1.25)
	}
	return surf
}

type gameState int

const (
	menuState gameState = iota
	playState
	gameOverState
	leaderboardState
)
