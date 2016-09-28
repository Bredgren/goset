package main

import "time"

// StateMgr keeps track of a current and previous GameState.
type StateMgr struct {
	prevState    GameState
	currentState GameState
}

// Goto transitions to the given state. If the current state is valid then it's Exit function
// is called first. The Enter function for the new state is called last.
func (s *StateMgr) Goto(gs GameState) {
	if s.currentState != nil {
		s.currentState.Exit()
		s.prevState = s.currentState
	}
	s.currentState = gs
	gs.Enter()
}

// Current returns the current state, or nil if there is none.
func (s *StateMgr) Current() GameState {
	return s.currentState
}

// Previous returns the previous state, or nill if there is none.
func (s *StateMgr) Previous() GameState {
	return s.prevState
}

// GameState reperents an arbitrary state used in a game.
type GameState interface {
	Enter()
	Update(t, dt time.Duration)
	Exit()
}
