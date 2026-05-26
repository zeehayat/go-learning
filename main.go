package main

import (
	"fmt"
)

type Reading struct {
	MeterID   string
	Kwh       float64
	Timestamp int64
}
type Validator interface {
	Validate() error
}

func (r *Reading) Validate() error {
	if r.Kwh < 0 {
		return fmt.Errorf("Negative Reading not allowed")
	}
	if r.MeterID == "" {
		return fmt.Errorf("MeterId is required")
	}
	return nil
}

func main() {
	validators := []Validator{
		&Reading{MeterID: "", Kwh: 1, Timestamp: 0},
		&Reading{MeterID: "meter1", Kwh: 100, Timestamp: 0},
		&Reading{MeterID: "XYZ", Kwh: -1, Timestamp: 0},
	}
	for _, v := range validators {
		if err := v.Validate(); err != nil {
			fmt.Println(err)
		}
	}
}
