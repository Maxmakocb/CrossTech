package cmd

import (
	"fmt"
	"cross_tech/parser"
	"cross_tech/database"
)

var data_path = "data/data.json"

func populate() {
	tracks, err := parser.ParseJsonFile("data/data.json")
	if err != nil {
		fmt.Printf("failed, parsing the json: %v", err)
		return
	}

	db_session, err := database.New()
	defer db_session.Close()
	if err != nil {
		fmt.Printf("failed, opening a new db connection session: %v", err)
		return
	}

	tracksAdded := 0
	tracksFailedToAdd := 0
	signalsAdded := 0
	signalsFailedToAdd := 0
	
	for _, track := range tracks {
		err = db_session.CreateTrack(track)
		if err != nil {
			fmt.Printf("failed, adding a track: %v\n", err)
			tracksFailedToAdd += 1
		} else {
			tracksAdded += 1
		}

		for _, signal := range track.SignalIds {
			err = db_session.CreateSignal(signal, track.TrackId)
			if err != nil {
				fmt.Printf("failed, adding a signal: %v\n", err)
				signalsFailedToAdd += 1
				continue
			} else {
				signalsAdded += 1
			}
		}
	}

	fmt.Println("Data insertion report:")
	fmt.Printf("Tracks added: %d\n", tracksAdded)
	fmt.Printf("Tracks failed to add: %d\n", tracksFailedToAdd)
	fmt.Printf("Signals added: %d\n", 	signalsAdded)
	fmt.Printf("Signals failed to add: %d\n", signalsFailedToAdd)
}