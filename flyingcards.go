package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/Bredgren/gogame"
	"github.com/Bredgren/gogame/composite"
	"github.com/Bredgren/wrand"
)

type flyingCardBg struct {
	surf   gogame.Surface
	groups [3]*flyingCardGroup
}

type flyingCardGroup struct {
	cards      [3]*flyingCard
	targetDist float64
	k          float64
	dampening  float64
	fading     bool
	fadeStart  time.Duration
	fadeTime   time.Duration
	active     bool
}

type flyingCard struct {
	surf       gogame.Surface
	posX, posY float64
	velX, velY float64
	mass       float64
}

func (b *flyingCardBg) update(t, dt time.Duration) {
	for i, g := range b.groups {
		if g == nil {
			g = &flyingCardGroup{}
			b.groups[i] = g
		}
		if !g.active && rand.Intn(100) < 2 {
			g.activate()
		}
		g.update(t, dt)
	}

	b.surf.Fill(gogame.FillBlack)
	for _, g := range b.groups {
		if !g.active {
			continue
		}
		for _, c := range g.cards {
			r := c.surf.Rect()
			r.SetCenter(c.posX, c.posY)
			if !g.fading {
				b.surf.Blit(c.surf, r.X, r.Y)
			} else {
				percent := (t - g.fadeStart).Seconds() / g.fadeTime.Seconds()
				if percent >= 1.0 {
					g.active = false
				} else {
					fade := gogame.NewSurface(int(r.W), int(r.H))
					fade.Fill(&gogame.FillStyle{Colorer: gogame.Color{A: 1.0 - percent}})
					fade.BlitComp(c.surf, 0, 0, composite.SourceIn)
					b.surf.Blit(fade, r.X, r.Y)
				}
			}
		}
	}
}

func (b *flyingCardBg) numActiveGroups() int {
	i := 0
	for _, g := range b.groups {
		if g.active {
			i++
		}
	}
	return i
}

func (g *flyingCardGroup) activate() {
	g.targetDist = rand.Float64()*10 + 50
	g.k = rand.Float64()*5 + 15
	g.dampening = rand.Float64()*0.07 + 0.85
	c1, c2, c3 := randomSet()
	r := gogame.MainDisplay().Rect().Inflate(100, 100)
	thickness := 75
	x1, y1 := randomPointInRectShell(int(r.W)+200, int(r.H)+200, thickness)
	x2, y2 := randomPointInRectShell(int(r.W)+200, int(r.H)+200, thickness)
	x3, y3 := randomPointInRectShell(int(r.W)+200, int(r.H)+200, thickness)
	g.cards = [3]*flyingCard{
		{surf: c1.surface(70, 100).Scaled(0.5, 0.5), posX: float64(x1 - 100), posY: float64(y1 - 100), mass: rand.Float64()*50 + 75},
		{surf: c2.surface(70, 100).Scaled(0.5, 0.5), posX: float64(x2 - 100), posY: float64(y2 - 100), mass: rand.Float64()*50 + 75},
		{surf: c3.surface(70, 100).Scaled(0.5, 0.5), posX: float64(x3 - 100), posY: float64(y3 - 100), mass: rand.Float64()*50 + 75},
	}
	g.fading = false
	g.fadeTime = time.Duration(2500 * time.Millisecond)
	g.active = true
}

func (g *flyingCardGroup) update(t, dt time.Duration) {
	if !g.active {
		return
	}

	// Returns the force that should be applied to c1 because of c2.
	forceBetween := func(c1, c2 *flyingCard) (fx, fy float64) {
		distX := c2.posX - c1.posX
		distY := c2.posY - c1.posY
		dist := math.Sqrt(distX*distX + distY*distY)
		f := (dist - g.targetDist) * g.k
		return f * (distX / dist), f * (distY / dist)
	}

	forceX01, forceY01 := forceBetween(g.cards[0], g.cards[1])
	forceX02, forceY02 := forceBetween(g.cards[0], g.cards[2])
	forceX12, forceY12 := forceBetween(g.cards[1], g.cards[2])

	g.cards[0].velX += (forceX01 + forceX02) / g.cards[0].mass
	g.cards[0].velX *= g.dampening
	g.cards[0].velY += (forceY01 + forceY02) / g.cards[0].mass
	g.cards[0].velY *= g.dampening

	g.cards[1].velX += (-forceX01 + forceX12) / g.cards[1].mass
	g.cards[1].velX *= g.dampening
	g.cards[1].velY += (-forceY01 + forceY12) / g.cards[1].mass
	g.cards[1].velY *= g.dampening

	g.cards[2].velX += (-forceX12 - forceX02) / g.cards[2].mass
	g.cards[2].velX *= g.dampening
	g.cards[2].velY += (-forceY12 - forceY02) / g.cards[2].mass
	g.cards[2].velY *= g.dampening

	for _, c := range g.cards {
		c.posX += c.velX * dt.Seconds()
		c.posY += c.velY * dt.Seconds()
	}

	dx01, dy01 := (g.cards[0].posX - g.cards[1].posX), (g.cards[0].posY - g.cards[1].posY)
	d01 := dx01*dx01 + dy01*dy01
	dx02, dy02 := (g.cards[0].posX - g.cards[2].posX), (g.cards[0].posY - g.cards[2].posY)
	d02 := dx02*dx02 + dy02*dy02
	dx12, dy12 := (g.cards[1].posX - g.cards[2].posX), (g.cards[1].posY - g.cards[2].posY)
	d12 := dx12*dx12 + dy12*dy12

	if d01+d02+d12 < math.Pow(g.targetDist, 3) && !g.fading {
		g.fading = true
		g.fadeStart = t
	}
}

func randomSet() (c1, c2, c3 card) {
	c1 = randomCard()
	c2 = randomCard()
	for c1 == c2 {
		c2 = randomCard()
	}
	c3 = missingCard(c1, c2)
	return
}

func randomCard() card {
	return card{
		count: map[int]count{0: one, 1: two, 2: three}[rand.Intn(3)],
		fill:  map[int]fill{0: empty, 1: solid, 2: line}[rand.Intn(3)],
		color: map[int]color{0: red, 1: green, 2: purple}[rand.Intn(3)],
		shape: map[int]shape{0: oval, 1: diamond, 2: tilde}[rand.Intn(3)],
	}
}

func missingCard(c1, c2 card) card {
	counts := map[count]bool{one: true, two: true, three: true}
	fills := map[fill]bool{empty: true, solid: true, line: true}
	colors := map[color]bool{red: true, green: true, purple: true}
	shapes := map[shape]bool{oval: true, diamond: true, tilde: true}

	counts[c1.count] = false
	counts[c2.count] = false

	fills[c1.fill] = false
	fills[c2.fill] = false

	colors[c1.color] = false
	colors[c2.color] = false

	shapes[c1.shape] = false
	shapes[c2.shape] = false

	c := card{}
	if c1.count == c2.count {
		c.count = c1.count
	} else {
		for i, ok := range counts {
			if ok {
				c.count = i
			}
		}
	}

	if c1.fill == c2.fill {
		c.fill = c1.fill
	} else {
		for i, ok := range fills {
			if ok {
				c.fill = i
			}
		}
	}

	if c1.color == c2.color {
		c.color = c1.color
	} else {
		for i, ok := range colors {
			if ok {
				c.color = i
			}
		}
	}

	if c1.shape == c2.shape {
		c.shape = c1.shape
	} else {
		for i, ok := range shapes {
			if ok {
				c.shape = i
			}
		}
	}

	return c
}

// randomPointInRectShell returns a uniformly random point inside a rectangluar shell.
// w and h are the outer width and height of the rectangle, thickness is the size of the shell.
func randomPointInRectShell(w, h, thickness int) (x, y int) {
	weights := make([]float64, w)
	for x := range weights {
		weight := 1.0
		if x <= thickness || x >= w-thickness {
			weight = float64(h) / float64(thickness*2)
		}
		weights[x] = weight
	}
	x = wrand.SelectIndex(weights)
	if x <= thickness || x >= w-thickness {
		y = rand.Intn(h)
	} else {
		y = rand.Intn(thickness * 2)
		if y > thickness {
			y += h - (thickness * 2)
		}
	}
	return
}
