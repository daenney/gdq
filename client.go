package gdq

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const trackerV1 = "https://gamesdonequick.com/tracker/api/v1"

// Client is a GDQ API client
type Client struct {
	ctx  context.Context
	c    *http.Client
	base string
}

// New creates a new GDQ API client
func New(ctx context.Context, client *http.Client) *Client {
	return &Client{
		ctx:  ctx,
		c:    client,
		base: trackerV1,
	}
}

// Latest returns the latest event
func (c *Client) Latest() (Event, error) {
	body, err := getWithCtx(c.ctx, c.c, fmt.Sprintf("%s/search?type=event", c.base))
	if err != nil {
		return Event{}, err
	}

	var resp = eventsResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return Event{}, err
	}

	if len(resp) == 0 {
		return Event{}, fmt.Errorf("there are no known events")
	}

	ev := resp[len(resp)-1]
	return Event{ID: ev.PK, Short: ev.Fields.Short, Name: ev.Fields.Name, Year: ev.Fields.Date.Year()}, nil
}

// GetSchedule returns the Schedule for a GDQ event
//
// A client has to be passed in. Please make sure to configure your client
// correctly, so not http.DefaultClient. Be nice to server admins and make
// sure your client sets a User-Agent header.
func (c *Client) GetSchedule(ev Event) (*Schedule, error) {
	grp, ctx := errgroup.WithContext(c.ctx)
	queries := []string{
		fmt.Sprintf("%s/search?type=run&event=%d", c.base, ev.ID),
		fmt.Sprintf("%s/search?type=runner&event=%d", c.base, ev.ID),
		fmt.Sprintf("%s/hosts/%d", c.base, ev.ID),
	}

	results := [3][]byte{}
	for i, query := range queries {
		i, query := i, query
		grp.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
			defer cancel()
			res, err := getWithCtx(ctx, c.c, query)
			if err == nil {
				results[i] = res
			}
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		return nil, err
	}

	var runs = runsResp{}
	err := json.Unmarshal(results[0], &runs)
	if err != nil {
		return nil, err
	}

	var runners = runnersResp{}
	err = json.Unmarshal(results[1], &runners)
	if err != nil {
		return nil, err
	}

	var hosts = hostsResp{}
	err = json.Unmarshal(results[2], &hosts)
	if err != nil {
		return nil, err
	}

	runnersByID := map[uint]runnerResp{}
	for _, runner := range runners {
		runnersByID[runner.ID] = runner
	}

	runsByID := map[uint]*Run{}
	finalRuns := []*Run{}

	for _, run := range runs {
		r := run.toRun()
		for _, runner := range run.Fields.Runners {
			r.Runners = append(r.Runners, runnersByID[runner].toRunner())
		}
		finalRuns = append(finalRuns, &r)
		runsByID[run.ID] = &r
	}

	for _, host := range hosts {
		r, ok := runsByID[host.Fields.Run]
		if !ok {
			continue
		}
		r.Hosts = append(r.Hosts, host.Fields.Handle)
	}

	schedule := NewScheduleFrom(finalRuns)

	return schedule, nil
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
		return ioutil.ReadAll(resp.Body)
	case http.StatusBadRequest:
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("client error: %w", err)
		}
		var cerr = struct {
			Error   string `json:"error"`
			Message string `json:"exception"`
		}{}
		err = json.Unmarshal(msg, &cerr)
		if err != nil {
			return nil, fmt.Errorf("client error, failed to unmarshal body: %w", err)
		}
		if cerr.Error == "" || cerr.Message == "" {
			return nil, fmt.Errorf("client error, unexpected body: %s", string(msg))
		}
		return nil, fmt.Errorf("client error: %s %s", cerr.Error, cerr.Message)
	case http.StatusNotFound:
		return nil, fmt.Errorf("resource not found, you probably have an error in your request URL")
	default:
		return nil, fmt.Errorf("received unexpected status code: %d (%s)", resp.StatusCode, resp.Status)
	}
}
