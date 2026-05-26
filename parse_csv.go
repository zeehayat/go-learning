package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type MeterReading struct {
	ConnectionID   string
	MeterSerial    string
	ConsumerID     string
	CurrentReading float64
}

func parseCSV(filePath string) ([]MeterReading, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var records []MeterReading
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		readingValue, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		m := MeterReading{
			ConnectionID:   record[0],
			MeterSerial:    record[1],
			ConsumerID:     record[2],
			CurrentReading: readingValue,
		}
		records = append(records, m)
	}

	return records, nil
}
