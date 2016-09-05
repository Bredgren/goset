package main

import (
	"time"

	"github.com/Bredgren/gogame"
	"github.com/Bredgren/gogame/event"
)

func main() {
	ready := gogame.Ready()
	go func() {
		<-ready
		setup()
		gogame.SetMainLoop(mainLoop)
	}()
}

func setup() {
	width, height := 500, 600
	display := gogame.MainDisplay()
	display.SetMode(width, height)
	display.Fill(gogame.FillBlack)
	display.Flip()
}

func mainLoop(t time.Duration) {
	// display := gogame.MainDisplay()
	// display.Fill(gogame.FillBlack)
	// display.Flip()
	for evt := event.Poll(); evt.Type != event.NoEvent; evt = event.Poll() {
		switch evt.Type {
		case event.Quit:
			gogame.UnsetMainLoop()
		case event.VideoResize:
			// data := evt.Data.(event.ResizeData)
		case event.KeyDown:
			// data := evt.Data.(event.KeyData)
		case event.KeyUp:
			// data := evt.Data.(event.KeyData)
		case event.MouseButtonDown:
			// data := evt.Data.(event.MouseData)
		case event.MouseButtonUp:
			// data := evt.Data.(event.MouseData)
		case event.MouseMotion:
			// data := evt.Data.(event.MouseMotionData)
		}
	}
}
