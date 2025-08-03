package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// CycleString handles both string and numeric cycle values from the API
type CycleString string

func (c *CycleString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*c = CycleString(s)
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*c = CycleString(strconv.FormatFloat(f, 'f', -1, 64))
		return nil
	}
	return fmt.Errorf("invalid cycle format: %s", string(data))
}

// EOLCycle represents a product lifecycle from the endoflife.date API
type EOLCycle struct {
	Cycle        CycleString `json:"cycle"`
	ReleaseDate  string      `json:"releaseDate"`
	EOL          interface{} `json:"eol"`
	Latest       string      `json:"latest"`
	Link         *string     `json:"link"`
	LTS          interface{} `json:"lts"`
	Support      interface{} `json:"support"`
	Discontinued interface{} `json:"discontinued"`
}

// EOLResult represents the analysis result for a container image
type EOLResult struct {
	Product        string `json:"product"`
	Version        string `json:"version"`
	Status         string `json:"status"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
	Link           string `json:"link"`
	EOLDate        string `json:"eolDate"`
	SupportEndDate string `json:"supportEndDate"`
	DaysRemaining  int    `json:"daysRemaining"`
	Latest         string `json:"latest"`
}

// Status constants
const (
	StatusCritical = "CRITICAL"
	StatusWarning  = "WARNING"
	StatusInfo     = "INFO"
	StatusOK       = "OK"
	StatusUnknown  = "UNKNOWN"
)
