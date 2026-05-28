package main

import (
	"context"
	"fmt"
	"time"
)

type Reading struct {
	MeterID   string
	Kwh       float64
	Timestamp int64
}
type Validator interface {
	Validate() error
}

type DBRepository struct {
	connString string
}

func (r *DBRepository) saveReading(ctx context.Context, reading *Reading) error {
	if err := reading.Validate(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout")

	}
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

func generateReadings(ch chan *Reading) {
	for i := 0; i < 5; i++ {
		ch <- &Reading{MeterID: fmt.Sprintf("MTR-%d", i), Kwh: float64(10*i + i), Timestamp: time.Now().Unix()}
	}
	close(ch)
}

func mainT() {
	readings := make(chan *Reading, 10)
	go generateReadings(readings)

	for r := range readings {
		if err := r.Validate(); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Processed Meter %s with %f kWh\n", r.MeterID, r.Kwh)
	}
}
