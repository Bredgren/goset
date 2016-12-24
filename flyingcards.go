package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/Bredgren/gogame/geo"
	"github.com/Bredgren/gogame/ggweb"
	"github.com/Bredgren/gogame/particle"
)

// import (
// 	"math"
// 	"math/rand"
// 	"time"

// 	"github.com/Bredgren/gogame"
// 	"github.com/Bredgren/gogame/composite"
// 	"github.com/Bredgren/wrand"
// )

// type flyingCardBg struct {
// 	surf   gogame.Surface
// 	groups [3]*flyingCardGroup
// }

type FlyingCard struct {
	particle.Particle
	Surf  *ggweb.Surface
	Group *FlyingCardGroup
}

func (c *FlyingCard) Draw(display *ggweb.Surface, t time.Duration) {
	if c.Group == nil || !c.Group.Active {
		return
	}
	r := c.Surf.Rect()
	r.SetCenter(c.Pos.X, c.Pos.Y)
	if !c.Group.Fading {
		display.Blit(c.Surf, r.X, r.Y)
	} else {
		display.Save()
		a := 1 - math.Min((t-c.Group.FadeStart).Seconds()/c.Group.FadeTime.Seconds(), 1)
		display.SetAlpha(a)
		display.Blit(c.Surf, r.X, r.Y)
		display.Restore()
	}
}

type FlyingCardGroup struct {
	Cards      [3]FlyingCard
	TargetDist float64
	K          float64
	Dampening  float64
	Fading     bool
	FadeStart  time.Duration
	FadeTime   time.Duration
	Active     bool
}

// func (b *flyingCardBg) update(t, dt time.Duration) {
// 	for i, g := range b.groups {
// 		if g == nil {
// 			g = &flyingCardGroup{}
// 			b.groups[i] = g
// 		}
// 		if !g.active && rand.Intn(100) < 2 {
// 			g.activate()
// 		}
// 		g.update(t, dt)
// 	}

// 	b.surf.Fill(gogame.FillBlack)
// 	for _, g := range b.groups {
// 		if !g.active {
// 			continue
// 		}
// 		for _, c := range g.cards {
// 			r := c.surf.Rect()
// 			r.SetCenter(c.posX, c.posY)
// 			if !g.fading {
// 				b.surf.Blit(c.surf, r.X, r.Y)
// 			} else {
// 				percent := (t - g.fadeStart).Seconds() / g.fadeTime.Seconds()
// 				if percent >= 1.0 {
// 					g.active = false
// 				} else {
// 					fade := gogame.NewSurface(int(r.W), int(r.H))
// 					fade.Fill(&gogame.FillStyle{Colorer: gogame.Color{A: 1.0 - percent}})
// 					fade.BlitComp(c.surf, 0, 0, composite.SourceIn)
// 					b.surf.Blit(fade, r.X, r.Y)
// 				}
// 			}
// 		}
// 	}
// }

// func (b *flyingCardBg) numActiveGroups() int {
// 	i := 0
// 	for _, g := range b.groups {
// 		if g.active {
// 			i++
// 		}
// 	}
// 	return i
// }

func (g *FlyingCardGroup) Activate(startPos geo.VecGen, mass, targetDist, k, dampening geo.NumGen) {
	g.TargetDist = targetDist()
	g.K = k()
	g.Dampening = dampening()
	c1, c2, c3 := randomSet()
	g.Cards = [3]FlyingCard{
		{Surf: c1.surface(70, 100).Scaled(0.5, 0.5), Particle: particle.Particle{Pos: startPos(), Mass: mass()}, Group: g},
		{Surf: c2.surface(70, 100).Scaled(0.5, 0.5), Particle: particle.Particle{Pos: startPos(), Mass: mass()}, Group: g},
		{Surf: c3.surface(70, 100).Scaled(0.5, 0.5), Particle: particle.Particle{Pos: startPos(), Mass: mass()}, Group: g},
	}
	g.Fading = false
	g.FadeTime = time.Duration(2500 * time.Millisecond)
	g.Active = true
}

func (g *FlyingCardGroup) Update(t, dt time.Duration) {
	if !g.Active {
		return
	}

	// Returns the force that should be applied to c1 because of c2.
	forceBetween := func(c1, c2 *FlyingCard) (f geo.Vec) {
		dir := c2.Pos.Minus(c1.Pos).Normalized()
		spring := (c2.Pos.Dist(c1.Pos) - g.TargetDist) * g.K
		v1 := c1.Vel.Dot(dir)
		v2 := -c2.Vel.Dot(dir)
		v := v1 + v2
		damp := g.Dampening * v
		return dir.Times(spring - damp)
	}

	force01 := forceBetween(&g.Cards[0], &g.Cards[1])
	force02 := forceBetween(&g.Cards[0], &g.Cards[2])
	force12 := forceBetween(&g.Cards[1], &g.Cards[2])

	g.Cards[0].ApplyForce(force01)
	g.Cards[0].ApplyForce(force02)

	g.Cards[1].ApplyForce(force01.Times(-1))
	g.Cards[1].ApplyForce(force12)

	g.Cards[2].ApplyForce(force02.Times(-1))
	g.Cards[2].ApplyForce(force12.Times(-1))

	// Try to prevent large speed increases due to leaving the browser tab in suspension
	if dt > time.Duration(33)*time.Millisecond {
		dt = time.Duration(33) * time.Millisecond
	}
	for i := range g.Cards {
		g.Cards[i].Update(dt)
	}

	if force01.Plus(force02).Plus(force12).Len() < 10 && !g.Fading {
		g.Fading = true
		g.FadeStart = t
	}

	if g.Fading && (t-g.FadeStart).Seconds()/g.FadeTime.Seconds() >= 1.0 {
		g.Active = false
	}
}

func (g *FlyingCardGroup) Draw(display *ggweb.Surface, t time.Duration) {
	for _, c := range g.Cards {
		c.Draw(display, t)
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
		color: map[int]color.Color{0: red, 1: green, 2: purple}[rand.Intn(3)],
		shape: map[int]shape{0: oval, 1: diamond, 2: tilde}[rand.Intn(3)],
	}
}

func missingCard(c1, c2 card) card {
	counts := map[count]bool{one: true, two: true, three: true}
	fills := map[fill]bool{empty: true, solid: true, line: true}
	colors := map[color.Color]bool{red: true, green: true, purple: true}
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

// // randomPointInRectShell returns a uniformly random point inside a rectangluar shell.
// // w and h are the outer width and height of the rectangle, thickness is the size of the shell.
// func randomPointInRectShell(w, h, thickness int) (x, y int) {
// 	weights := make([]float64, w)
// 	for x := range weights {
// 		weight := 1.0
// 		if x <= thickness || x >= w-thickness {
// 			weight = float64(h) / float64(thickness*2)
// 		}
// 		weights[x] = weight
// 	}
// 	x = wrand.SelectIndex(weights)
// 	if x <= thickness || x >= w-thickness {
// 		y = rand.Intn(h)
// 	} else {
// 		y = rand.Intn(thickness * 2)
// 		if y > thickness {
// 			y += h - (thickness * 2)
// 		}
// 	}
// 	return
// }
