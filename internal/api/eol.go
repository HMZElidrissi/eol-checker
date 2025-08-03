package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/HMZElidrissi/eol-checker/internal/models"
)

const (
	EOLAPIBaseURL  = "https://endoflife.date/api"
	RequestTimeout = 10 * time.Second
)

// Client represents an API client for endoflife.date
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new EOL API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
		baseURL: EOLAPIBaseURL,
	}
}

// GetProductCycles fetches EOL cycles for a given product
func (c *Client) GetProductCycles(product string) ([]models.EOLCycle, error) {
	url := fmt.Sprintf("%s/%s.json", c.baseURL, product)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Product not found, not an error
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var cycles []models.EOLCycle
	if err := json.NewDecoder(resp.Body).Decode(&cycles); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return cycles, nil
}
