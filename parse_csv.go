package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
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

func parseCSV(filePath string) ([]MeterReading, int, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to open log file: %w", err)
	}
	defer logFile.Close()
	logWriter := bufio.NewWriter(logFile)
	defer logWriter.Flush()

	var skippedRows int = 0

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return nil, 0, err
	}

	var records []MeterReading

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		if len(record) < 4 {
			return nil, 0, fmt.Errorf("Invalid row :%v", record)
		}
		if record[3] == "" {
			skippedRows++
			logSkippedRow(logWriter, "invalid reading(empty)", record)
			continue
		}
		readingValue, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			fmt.Println("Invalid reading(non numeric)", record, err)
			continue
		}
		m := MeterReading{
			ConnectionID:   record[0],
			MeterSerial:    record[1],
			ConsumerID:     record[2],
			CurrentReading: readingValue,
		}
		records = append(records, m)
	}

	return records, skippedRows, nil
}
func logSkippedRow(writer *bufio.Writer, reason string, row []string) {

	line := fmt.Sprintf("%s: %v\n", reason, row)
	writer.WriteString(line)
}

func main() {
	readings, skippedRows, err := parseCSV("DATA/ashuran_april_readings.csv")
	if err != nil {
		fmt.Println("The file couldn't be parsed:", err)
		return
	}
	fmt.Println("Total Readings:", len(readings))

	fmt.Println("Skipped Rows:", skippedRows)
}
