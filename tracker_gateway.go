package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

const ContentTypeJSON = "application/json"

type TrackerGateway struct {
	url      string
	apiToken string
	client   http.Client
}

func NewTrackerGateway(url, apiToken string) *TrackerGateway {
	return &TrackerGateway{
		url:      url,
		apiToken: apiToken,
		client: http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (t *TrackerGateway) CreateStory(story Story) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(story)
	if err != nil {
		return fmt.Errorf("Could not create Tracker story: %w", err)
	}

	apiPath := fmt.Sprintf("v5/projects/%d/stories", story.ProjectID)
	req, err := t.newRequest(apiPath, b)
	if err != nil {
		return fmt.Errorf("Could not create Tracker story: %w", err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("Could not create Tracker story: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf(
			"Could not create Tracker story, recieved a non successful HTTP status: %d (%s)",
			resp.StatusCode, resp.Status)
	}

	return nil
}

func (t *TrackerGateway) newRequest(apiPath string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(t.url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, apiPath)
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-TrackerToken", t.apiToken)
	req.Header.Add("Content-Type", ContentTypeJSON)
	req.Header.Add("Accept", ContentTypeJSON)

	return req, nil
}
