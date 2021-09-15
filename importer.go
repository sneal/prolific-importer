package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Label struct {
	Name string `json:"name"`
}

type Story struct {
	ProjectID   int     `json:"project_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	StoryType   string  `json:"story_type"`
	Labels      []Label `json:"labels"`
}

type Importer struct {
	tracker *TrackerGateway
}

func NewCsvImporter(tracker *TrackerGateway) *Importer {
	return &Importer{
		tracker: tracker,
	}
}

func (i *Importer) ImportCsvFromFile(trackerProjectID int, csvFilePath string) error {
	content, err := ioutil.ReadFile(csvFilePath)
	if err != nil {
		return fmt.Errorf("Couldn't load file %s: %w", csvFilePath, err)
	}
	return i.ImportCsv(trackerProjectID, string(content))
}

func (i *Importer) ImportCsv(trackerProjectID int, csvContent string) error {
	stories, err := parseIntoStories(csvContent)
	if err != nil {
		return fmt.Errorf("Couldn't parse CSV into stories: %w", err)
	}
	return i.CreateStories(trackerProjectID, stories)
}

func parseIntoStories(csvContent string) ([]Story, error) {
	var stories []Story

	r := strings.NewReader(csvContent)
	csvReader := csv.NewReader(r)
	header, err := csvReader.Read()
	if err != nil && err != io.EOF {
		return stories, err
	}
	if !isCSVHeadersValid(header) {
		return stories,
			fmt.Errorf("Expected CSV in the following format: Title,Type,Description,Labels but instead got %v", header)
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return stories, err
		}
		story := Story{
			Name:        record[0],
			StoryType:   record[1],
			Description: record[2],
			Labels:      parseIntoLabels(record[3]),
		}
		stories = append(stories, story)
	}

	return stories, nil
}

func (i *Importer) CreateStories(trackerProjectID int, stories []Story) error {
	for _, story := range stories {
		story.ProjectID = trackerProjectID
		err := i.tracker.CreateStory(story)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseIntoLabels(labelRecord string) []Label {
	rawLabels := strings.Split(labelRecord, ",")
	var labels = make([]Label, len(rawLabels))
	for i, l := range rawLabels {
		labels[i] = Label{
			Name: strings.TrimSpace(l),
		}
	}
	return labels
}

func isCSVHeadersValid(header []string) bool {
	return header[0] == "Title" && header[1] == "Type" && header[2] == "Description" && header[3] == "Labels"
}
