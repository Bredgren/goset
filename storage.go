package main

import (
	"encoding/json"
	"time"

	"github.com/Bredgren/gogame"
)

// Local storage keys
const (
	saveData   = "saveData"
	leaderData = "leaderData"
)

type leaderboardEntry struct {
	numSets   int
	numErrors int
	time      time.Duration
	date      time.Time
}

func getLeaderboardData() []leaderboardEntry {
	d, ok := gogame.LocalStorageGet(leaderData)
	if !ok {
		return []leaderboardEntry{}
	}
	entries := []leaderboardEntry{}
	e := json.Unmarshal([]byte(d), entries)
	if e != nil {
		// Must be saved in improper format somehow, just delete it.
		gogame.Log("Error reading leaderboard data:", e.Error())
		gogame.LocalStorageRemove(leaderData)
		return []leaderboardEntry{}
	}
	return entries
}

func saveLeaderboardData(data []leaderboardEntry) {
	str, e := json.Marshal(data)
	if e != nil {
		// This actually a problem, it shouldn't fail to marshal.
		panic(e)
	}
	gogame.LocalStorageSet(leaderData, string(str))
}
