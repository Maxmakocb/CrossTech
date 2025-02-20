package data_structures

// Track is one of the datastructures, stored in the original json file.
type Track struct {
	TrackId   int `json:"track_id"`
	Source 	  string `json:"source"`
	Target 	  string `json:"target"`
	SignalIds []SignalId `json:"signal_ids"`
}

type TrackStandalone struct{
	TrackId   int `json:"track_id"`
	Source 	  string `json:"source"`
	Target 	  string `json:"target"`
}

func TrackToStandalone(t Track) *TrackStandalone {
	return &TrackStandalone{
		TrackId: t.TrackId,
		Source:  t.Source,
		Target:  t.Target,
	}
}
