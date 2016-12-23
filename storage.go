package main

import (
	"encoding/json"
	"time"

	"github.com/Bredgren/gogame/ggweb"
)

// import (
// 	"encoding/json"
// 	"time"

// 	"github.com/Bredgren/gogame"
// )

// Local storage keys
const (
	saveData   = "saveData"
	leaderData = "leaderData"
)

type SaveData struct {
	PlayTime    time.Duration `json:"play_time"`
	ActiveCards []card        `json:"active_cards"`
	Deck        []card        `json:"deck"`
	Sets        int           `json:"sets"`
	Errors      int           `json:"errors"`
}

func getSaveData() (SaveData, bool) {
	strData, ok := ggweb.LocalStorageGet(saveData)
	if !ok {
		return SaveData{}, false
	}
	data := SaveData{}
	e := json.Unmarshal([]byte(strData), &data)
	if e != nil {
		// Must be saved in improper format somehow, just delete it.
		ggweb.Log("Error reading saved game data:", e.Error())
		ggweb.LocalStorageRemove(saveData)
		return SaveData{}, false
	}
	return data, true
}

func setSaveData(d *SaveData) error {
	strData, e := json.Marshal(d)
	if e != nil {
		return e
	}
	ggweb.LocalStorageSet(saveData, string(strData))
	return nil
}

func clearSaveData() {
	ggweb.LocalStorageRemove(saveData)
}

// type leaderboardEntry struct {
// 	numSets   int
// 	numErrors int
// 	time      time.Duration
// 	date      time.Time
// }

// func getLeaderboardData() []leaderboardEntry {
// 	d, ok := gogame.LocalStorageGet(leaderData)
// 	if !ok {
// 		return []leaderboardEntry{}
// 	}
// 	entries := []leaderboardEntry{}
// 	e := json.Unmarshal([]byte(d), entries)
// 	if e != nil {
// 		// Must be saved in improper format somehow, just delete it.
// 		gogame.Log("Error reading leaderboard data:", e.Error())
// 		gogame.LocalStorageRemove(leaderData)
// 		return []leaderboardEntry{}
// 	}
// 	return entries
// }

// func saveLeaderboardData(data []leaderboardEntry) {
// 	str, e := json.Marshal(data)
// 	if e != nil {
// 		// This actually a problem, it shouldn't fail to marshal.
// 		panic(e)
// 	}
// 	gogame.LocalStorageSet(leaderData, string(str))
// }
