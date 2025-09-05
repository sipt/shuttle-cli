package apis

import (
	"fmt"
	"strings"

	"github.com/sipt/shuttle/controller/model"

	apipkg "github.com/sipt/shuttle/dns/api"
)

// GetDNSCache gets DNS cache list
func (c *APIClient) GetDNSCache(filter string) error {
	result := &model.Response[[]apipkg.DNS]{}
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
		if filter != "" && !strings.Contains(cache.Domain, filter) {
			continue
		}
		fmt.Printf("  %-50s -> %-20s (%s) [%s] [%s]\n",
			cache.Domain, cache.IP, cache.Typ, cache.DNSServer, cache.CountryCode)
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
