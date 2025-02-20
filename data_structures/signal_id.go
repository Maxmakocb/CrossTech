package data_structures

// SignalId is one of the datastructures, stored in the original json file.
type SignalId struct {
	SignalId    int `json:"signal_id"`
	SignalName  string `json:"signal_name"`
	Elr         string `json:"elr"`
	Mileage     float64 `json:"mileage"`
}
