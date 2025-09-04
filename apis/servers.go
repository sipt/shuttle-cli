package apis

import (
	"fmt"

	"github.com/sipt/shuttle/controller/model"

	serverpkg "github.com/sipt/shuttle/server"
	apipkg "github.com/sipt/shuttle/server/api"
)

// Server server structure
type Server struct {
	Name string `json:"name"`
	Type string `json:"typ"`
	RTT  string `json:"rtt"`
}

// GetServers gets server list
func (c *APIClient) GetServers() error {
	result := &model.Response[[]apipkg.ItemResponse]{}
	resp, err := c.client.R().
		SetResult(result).
		Get("/api/servers")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}
	fmt.Println("Servers:")
	for _, server := range result.Data {
		if server.Name == serverpkg.Direct || server.Name == serverpkg.Reject {
			continue
		}
		fmt.Printf("  %s - ", server.Name)
		// Color code RTT based on value
		PrintRtt(server.RTT)
		fmt.Println()
	}

	return nil
}

// GetServer gets specific server information
func (c *APIClient) GetServer(name string) error {
	result := &model.Response[apipkg.ItemResponse]{}
	resp, err := c.client.R().
		SetResult(result).
		Get(fmt.Sprintf("/api/servers/%s", name))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}

	server := result.Data

	fmt.Printf("Server Information:\n")
	fmt.Printf("  Name: %s\n", server.Name)
	fmt.Printf("  Type: %s\n", server.Typ)
	fmt.Printf("  RTT: ")
	PrintRtt(server.RTT)
	fmt.Println()

	return nil
}

// TestServerRTT tests server latency
func (c *APIClient) TestServerRTT(name string) error {
	result := &model.Response[apipkg.ItemResponse]{}
	resp, err := c.client.R().
		SetResult(result).
		Put(fmt.Sprintf("/api/servers/%s/rtt", name))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}
	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}

	server := result.Data

	fmt.Printf("Server Information:\n")
	fmt.Printf("  Name: %s\n", server.Name)
	fmt.Printf("  Type: %s\n", server.Typ)
	fmt.Printf("  RTT: ")
	PrintRtt(server.RTT)
	fmt.Println()

	return nil
}
