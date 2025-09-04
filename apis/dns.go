package apis

import (
	"fmt"

	"github.com/sipt/shuttle/controller/model"
)

// DNSCache DNS cache structure
type DNSCache struct {
	Type        string `json:"typ"`
	Domain      string `json:"domain"`
	IP          string `json:"ip"`
	DNSServer   string `json:"dns_server"`
	CountryCode string `json:"country_code"`
}

// GetDNSCache gets DNS cache list
func (c *APIClient) GetDNSCache() error {
	result := &model.Response[[]DNSCache]{}
	resp, err := c.client.R().
		SetResult(result).
		Get("/api/dns/cache")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}

	fmt.Println("DNS Cache:")
	if len(result.Data) == 0 {
		fmt.Println("  No cache entries found")
		return nil
	}

	for _, cache := range result.Data {
		fmt.Printf("  %s -> %s (%s) [%s] [%s]\n",
			cache.Domain, cache.IP, cache.Type, cache.DNSServer, cache.CountryCode)
	}

	return nil
}

// ClearDNSCache clears DNS cache
func (c *APIClient) ClearDNSCache() error {
	result := &model.Response[any]{}
	resp, err := c.client.R().
		SetResult(result).
		Delete("/api/dns/cache")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code == 0 {
		green.Println("DNS cache cleared successfully")
	} else {
		fmt.Print("Clear DNS cache failed: ")
		red.Println(result.Message)
	}

	return nil
}
