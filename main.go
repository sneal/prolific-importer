package main

/*
prolific-importer: import many tracker stories
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const TrackerBaseURL = "https://www.pivotaltracker.com/services"

func main() {
	if len(os.Args) < 3 {
		PrintUsageAndExit()
	}

	apiToken := os.Args[1]
	projectID, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Expected an integer Tracker project id, but got: %s", os.Args[2])
		os.Exit(1)
	}

	tracker := NewTrackerGateway(TrackerBaseURL, apiToken)
	importer := NewCsvImporter(tracker)

	content := readStdin()
	if len(os.Args) == 3 && content != nil {
		fmt.Fprintf(os.Stderr, "Importing STDIN\n")
		err := importer.ImportCsv(projectID, string(content))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(os.Args) != 4 || os.Args[1] == "help" {
		PrintUsageAndExit()
	}

	csvPath := os.Args[3]
	fmt.Fprintf(os.Stderr, "Importing %s\n", csvPath)
	err = importer.ImportCsvFromFile(projectID, csvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func readStdin() []byte {
	stat, err := os.Stdin.Stat()
	if err != nil || (stat.Mode()&os.ModeCharDevice) != 0 {
		return nil
	}

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil
	}
	return stdin
}
