package main

import (
	"fmt"
	"os"
)

func PrintUsageAndExit() {
	fmt.Println(`prolific-importer v0.1

Usage:
    prolific-importer [API_TOKEN] [PROJECT_ID] [FILE]
        Imports the prolific generated Pivotal Tracker CSV file FILE to the Pivotal Tracker project PROJECT_ID.
        If FILE is not specified, content is read from STDIN.

        A prolific file can be converted and imported into Pivotal Tracker in one step:

            prolific stories.prolific | prolific-importer "$API_TOKEN" 12345

        is a useful one-liner.

    prolific-importer help
        You're looking at it!`)
	os.Exit(1)
}
