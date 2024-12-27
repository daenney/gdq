package gdq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	trackerV2 = "https://gamesdonequick.com/tracker/api/v2"
)

// Client is a GDQ API client.
type Client struct {
	c  *http.Client
	v2 string
}

// New creates a new GDQ API client.
func New(client *http.Client) *Client {
	return &Client{
		c:  client,
		v2: trackerV2,
	}
}

// Events returns all events, sorted by start date.
func (c *Client) Events(ctx context.Context) ([]*Event, error) {
	resp, err := fromJSON[eventsResp](ctx, c.c, fmt.Sprintf("%s/events/", c.v2))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve events: %w", err)
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("there are no known events")
	}

	return resp.toEvents(), nil
}

// Event retrieves event information for the event ID.
//
// Contrary to [Client.Events], this includes donation information.
func (c *Client) Event(ctx context.Context, ev uint) (*Event, error) {
	resp, err := fromJSON[eventResp](ctx, c.c, fmt.Sprintf("%s/events/%d/?totals=true", c.v2, ev))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data for event %d: %w", ev, err)
	}

	return resp.toEvent(), nil
}

// Runs returns all runs for an Event.
func (c *Client) Runs(ctx context.Context, ev uint) ([]*Run, error) {
	resp, err := fromJSON[runResp](ctx, c.c, fmt.Sprintf("%s/events/%d/runs/", c.v2, ev))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve runs for event %d: %w", ev, err)
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("there are no interviews")
	}

	return resp.toRuns(), nil
}

// Schedule returns the [Schedule] for a GDQ event.
//
// This schedule only contains runs, not interviews.
func (c *Client) Schedule(ctx context.Context, ev uint) (*Schedule, error) {
	runs, err := c.Runs(ctx, ev)
	if err != nil {
		return nil, err
	}
	return NewScheduleFrom(runs), nil
}

func getWithCtx(ctx context.Context, c *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return io.ReadAll(resp.Body)
	case http.StatusBadRequest, http.StatusNotFound:
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("client error: %w", err)
		}
		cerr := struct {
			Detail string `json:"detail"`
		}{}
		err = json.Unmarshal(msg, &cerr)
		if err != nil {
			return nil, fmt.Errorf("client error, failed to unmarshal body: %w", err)
		}
		if cerr.Detail == "" {
			return nil, fmt.Errorf("client error, unexpected body: %s", string(msg))
		}
		return nil, fmt.Errorf("client error: %s ", cerr.Detail)
	default:
		return nil, fmt.Errorf("received unexpected status code: %s", resp.Status)
	}
}

func fromJSON[T any](ctx context.Context, c *http.Client, endpoint string) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	data, err := getWithCtx(ctx, c, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from: %s, %w", endpoint, err)
	}

	var resp T
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &resp, nil
}
