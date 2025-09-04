package apis

import (
	"fmt"

	"github.com/sipt/shuttle/controller/model"
	"resty.dev/v3"

	apipkg "github.com/sipt/shuttle/cmd/api"
)

// Inbound inbound listener structure
type Inbound struct {
	Name string `json:"name"`
	Type string `json:"typ"`
	Addr string `json:"addr"`
}

// APIClient API client
type APIClient struct {
	client *resty.Client
	host   string
}

// NewAPIClient creates a new API client
func NewAPIClient(host string) *APIClient {
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://%s", host))

	return &APIClient{
		client: client,
		host:   host,
	}
}

// GetStatus gets system status
func (c *APIClient) GetStatus() error {
	result := &model.Response[string]{}
	resp, err := c.client.R().
		SetResult(result).
		Get("/api/status")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}

	fmt.Print("System status: ")
	switch status := result.Data; status {
	case apipkg.StatusStarting, apipkg.StatusRunning:
		green.Println(status)
	case apipkg.StatusStopped:
		red.Println(status)
	default:
		white.Println(status)
	}
	return nil
}

// GetInbounds gets inbound listeners list
func (c *APIClient) GetInbounds() error {
	result := &model.Response[[]Inbound]{}
	resp, err := c.client.R().
		SetResult(result).
		Get("/api/inbounds")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	fmt.Println("Inbound listeners:")
	for _, inbound := range result.Data {
		fmt.Printf("  Name: %s, Type: %s, Address: %s\n", inbound.Name, inbound.Type, inbound.Addr)
	}

	return nil
}

// Reload reloads configuration
func (c *APIClient) Reload() error {
	result := &model.Response[any]{}
	resp, err := c.client.R().
		SetResult(result).
		Put("/api/reload")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code == 0 {
		green.Println("success")
	} else {
		fmt.Print("Reload failed: ")
		red.Println(result.Message)
	}
	return nil
}
