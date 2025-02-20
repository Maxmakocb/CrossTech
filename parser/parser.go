package parser

import (
	"os"
	"encoding/json"

	"strings"

	"cross_tech/data_structures"
)

// ParseJsonFile parses a json file into a go datastructure.
func ParseJsonFile(path string) ([]data_structures.Track, error) {
	b, err := os.ReadFile("data/data.json")
    if err != nil {
		return nil, err
    }

	// If the input data has NaNs in unexpected fields of json, this will break,
	// however this works for the current objective of working with a given file, so, I think, this is good enough for now.
	b = []byte(strings.ReplaceAll(string(b), "NaN", "-1"))
	var tracks []data_structures.Track
	err = json.Unmarshal(b, &tracks)
	if err != nil {
		return nil, err
    }

	return tracks, nil
}