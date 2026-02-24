package match

import (
	"context"

	// "encoding/json"
	// "errors"
	// "fmt"
	// "log"
	"net/http"
	"time"
)

type Client interface {
	GetMatch(ctx context.Context, matchID uint64) (*MatchResponse, error)
}

type httpClient struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) Client {
	return &httpClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *httpClient) GetMatch(ctx context.Context, matchID uint64) (*MatchResponse, error) {
	/*url := fmt.Sprintf("%s/matches/%d", c.baseURL, matchID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("match not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("match service returned %d", resp.StatusCode)
	}

	var match MatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
		return nil, err
	}
	*/

	match := MatchResponse{ID: 1, Status: "live"}
	return &match, nil
}
