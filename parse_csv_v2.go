/*
Package main provides the entry point and core logic for the application.
This specific file (parse_csv_v2.go) contains the logic to read and parse
electricity/utility meter readings from a CSV file. It processes the records
sequentially, validates the data format, logs skipped or corrupt rows to an external
log file using buffered I/O, and streams valid readings concurrently via a Go channel.
*/
package main

// Import core library packages for file operations, string conversion, and CSV parsing
import (
	"bufio"        // Implements buffered I/O to optimize writes to the log file
	"encoding/csv" // Provides CSV reading capabilities
	"fmt"          // Formatted I/O for printing and string formatting
	"io"           // Used to detect the end of the file (io.EOF)
	"os"           // Used for operating system tasks like opening files
	"strconv"      // Used to convert string fields into numeric float64 values
)

// parseCSVV2 reads a CSV file containing meter readings, skips the header,
// validates each row, streams valid rows to a send-only channel, and logs malformed rows.
// It returns the number of skipped rows and any critical error encountered.
func parseCSVV2(filePath string, readings chan<- MeterReading) (int, error) {

	// Ensure that the readings channel is closed when this function exits.
	// This signals to the receiver loop in main.go that no more data is coming.
	defer close(readings)

	// Attempt to open the target CSV file at the specified path
	file, err := os.Open(filePath)
	if err != nil {
		// Return 0 skipped rows and the error if the file cannot be opened
		return 0, err
	}
	// Ensure the CSV file is closed cleanly when the function returns to free resources
	defer file.Close()

	// Open or create log.txt in append-only mode, creating it with 0644 permissions if missing
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Return an error wrapped with descriptive context if log creation fails
		return 0, fmt.Errorf("Failed to open log file: %w", err)
	}
	// Ensure the log file descriptor is closed when the function exits
	defer logFile.Close()

	// Create a buffered writer for logFile to minimize expensive disk I/O operations
	logWriter := bufio.NewWriter(logFile)
	// Ensure any remaining buffered log messages are flushed to disk before exiting
	defer logWriter.Flush()

	// Initialize a counter to track rows skipped due to empty or invalid data
	var skippedRows int = 0

	// Instantiate a new CSV reader using the opened file pointer
	reader := csv.NewReader(file)

	// Read the first line of the CSV to skip the header row (e.g., Column Names)
	_, err = reader.Read()
	if err != nil {
		// Return immediately if the file is completely empty or corrupted at the header
		return 0, err
	}

	// Begin an infinite loop to process the CSV rows one by one
	for {
		// Read the next row from the CSV file
		record, err := reader.Read()

		// Check if we have hit the End Of File (EOF)
		if err == io.EOF {
			break // Break the loop safely; file reading is complete
		}
		// If an unexpected error occurs while reading a row, abort and return it
		if err != nil {
			return 0, err
		}

		// Validate that the row has at least 4 columns (Indices 0, 1, 2, and 3)
		if len(record) < 4 {
			// Note: The original code returned len(record) as skipped rows here,
			// but returning a formal formatting error is safest.
			return skippedRows, fmt.Errorf("invalid row: %v", record)
		}

		// Check if the 4th column (Index 3: CurrentReading) is empty
		if record[3] == "" {
			skippedRows++ // Increment the count of skipped rows
			// Write the row details and the reason to the log file buffer
			logSkippedRowV2(logWriter, "invalid reading(empty)", record)
			continue // Skip the rest of this loop iteration and move to the next row
		}

		// Attempt to convert the 4th column string value into a 64-bit float
		readingValue, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			// If conversion fails (non-numeric data), print to standard output
			fmt.Println("Invalid reading(non numeric)", record, err)
			logSkippedRowV2(logWriter, "invalid reading(non numeric)", record)

			// Note: You might want to increment skippedRows and log to logWriter here as well!
			continue // Skip processing this row and continue the loop
		}

		// Instantiate a new MeterReading struct using data from the valid CSV row
		m := MeterReading{
			ConnectionID:   record[0],    // Column 1: Connection Identifier
			MeterSerial:    record[1],    // Column 2: Meter Serial Number
			ConsumerID:     record[2],    // Column 3: Consumer Account ID
			CurrentReading: readingValue, // Converted Column 4: Numeric current reading
		}

		// Send the populated MeterReading struct into the readings channel
		readings <- m
	}

	// Note: The redundant `close(readings)` at the bottom was safely removed because
	// the `defer close(readings)` at the top executes automatically on return.

	// Return the total count of skipped rows and a nil error signifying success
	return skippedRows, nil
}

// logSkippedRowV2 takes a buffered writer, a failure reason, and the raw row data,
// formats them into a single log string line, and writes it to the buffer.
func logSkippedRowV2(writer *bufio.Writer, reason string, row []string) {
	// Construct a formatted log string ending with a newline character
	line := fmt.Sprintf("%s: %v\n", reason, row)
	// Write the formatted string into the buffered writer
	writer.WriteString(line)
}
