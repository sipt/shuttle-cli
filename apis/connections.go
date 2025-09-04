package apis

import (
	"fmt"
)

// Connection connection structure
type Connection struct {
	Input  ConnectionInfo `json:"input"`
	Output ConnectionInfo `json:"output"`
}

type ConnectionInfo struct {
	ID         int64  `json:"id"`
	LocalAddr  string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
}

// GetConnections gets current connections list
func (c *APIClient) GetConnections() error {
	result := &[]Connection{}
	resp, err := c.client.R().
		SetResult(result).
		Get("/api/conns")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	fmt.Println("Current Connections:")
	if len(*result) == 0 {
		fmt.Println("  No active connections")
		return nil
	}

	fmt.Printf("%-8s %-25s %-25s %-8s %-25s %-25s\n",
		"InputID", "Input Local", "Input Remote", "OutputID", "Output Local", "Output Remote")
	fmt.Println("---------------------------------------------------------------------------------------------------------------")

	for _, conn := range *result {
		fmt.Printf("%-8d %-25s %-25s %-8d %-25s %-25s\n",
			conn.Input.ID,
			conn.Input.LocalAddr,
			conn.Input.RemoteAddr,
			conn.Output.ID,
			conn.Output.LocalAddr,
			conn.Output.RemoteAddr,
		)
	}

	fmt.Printf("\nTotal connections: %d\n", len(*result))
	return nil
}
