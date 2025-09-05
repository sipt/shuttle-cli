package apis

import (
	"fmt"
	"strings"

	"github.com/sipt/shuttle/controller/model"

	apipkg "github.com/sipt/shuttle/events/record"
)

// GetRecords gets connection records list
func (c *APIClient) GetRecords(filter string) error {
	result := &model.Response[[]apipkg.RecordEntity]{}
	resp, err := c.client.R().
		SetResult(result).
		Get("/api/records")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}

	fmt.Println("Connection Records:")
	if len(result.Data) == 0 {
		fmt.Println("  No records found")
		return nil
	}

	fmt.Printf("%-6s %-50s %-30s %-10s %-10s %-12s %-8s\n",
		"ID", "Destination", "Policy", "Upload", "Download", "Status", "Protocol")
	fmt.Println("---------------------------------------------------------------------------------------------------------------")

	for _, record := range result.Data {
		if filter != "" && !strings.Contains(record.DestAddr, filter) {
			continue
		}

		fmt.Printf("%d %-50s %-30s %-10s %-10s %-12s %-8s\n",
			record.ID,
			record.DestAddr,
			record.Policy,
			formatBytes(record.Up),
			formatBytes(record.Down),
			record.Status,
			record.Protocol,
		)
	}

	fmt.Printf("\nTotal records: %d\n", len(result.Data))
	return nil
}

// ClearRecords clears connection records
func (c *APIClient) ClearRecords() error {
	result := &model.Response[any]{}
	resp, err := c.client.R().
		SetResult(result).
		Delete("/api/records")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code == 0 {
		green.Println("Connection records cleared successfully")
	} else {
		fmt.Print("Clear records failed: ")
		red.Println(result.Message)
	}

	return nil
}

// formatBytes formats bytes to human readable format
func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1fGB", float64(bytes)/(1024*1024*1024))
	}
}
