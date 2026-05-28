/*
Package main acts as the orchestration point for the application.
This file opens a log file, sets up a channel for communication, spawns a concurrent
goroutine to parse the CSV records using `parseCSVV2`, and reads items off the channel
to print them out and count total processed records.
*/
package main

// Import necessary core library packages
import (
	// Used for buffered writing to log files
	"fmt" // Used for basic printing to the console terminal
)

// parseResult wraps the final outcome of the background CSV processing routine
type parseResult struct {
	skippedRows int   // Total lines skipped during the parse
	err         error // Holds any critical failure error encountered
}

func main() {
	// Open or create a log.txt file in append-only mode for main execution logging
	//	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	// If opening the file fails, display the error and terminate execution early
	//fmt.Println("Failed to open log file: ", err)
	//return
	//}
	// Ensure the file handle is safely closed when main finishes execution
	//defer logFile.Close()

	// Create a buffered writer wrapping the file handle for high performance
	//logWriter := bufio.NewWriter(logFile)
	// Ensure any data waiting in the buffer gets pushed to disk before the app exits
	//defer logWriter.Flush()

	// Instantiate an unbuffered channel capable of transmitting MeterReading structures
	readings := make(chan MeterReading)
	results := make(chan parseResult, 1)
	// Spawn a concurrent background thread (Goroutine) to handle file parsing
	// so the main thread can read data instantly as it becomes available.
	go func() {
		// Execute the CSV parsing function in the background.
		// It will stream parsed entries directly into the 'readings' channel.
		skippedRows, err := parseCSVV2("DATA/ashuran_april_readings.csv", readings)
		results <- parseResult{skippedRows: skippedRows, err: err}
	}()

	// Initialize a counter to track the total number of successful records processed
	count := 0

	// Loop continuously reads from the 'readings' channel.
	// The loop terminates automatically when the channel is closed inside parseCSVV2.
	for reading := range readings {
		// Print each received meter reading structure to the console
		fmt.Println(reading)
		// Increment the total successfully read counter
		count++
	}

	// Output the total count of successfully processed readings to the console
	fmt.Println("Total Readings: ", count)
}
