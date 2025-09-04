package apis

import (
	"fmt"

	"github.com/sipt/shuttle/controller/model"
)

// GetRuleMode gets current routing mode
func (c *APIClient) GetRuleMode() error {
	result := &model.Response[string]{}
	resp, err := c.client.R().
		SetResult(result).
		Get("/api/rule/mode")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}

	fmt.Print("Current routing mode: ")
	switch mode := result.Data; mode {
	case "ModeDirect":
		green.Println("Direct")
	case "ModeGlobal":
		orange.Println("Global")
	case "ModeRule":
		white.Println("Rule")
	default:
		white.Println(mode)
	}

	return nil
}

// SetRuleMode sets routing mode
func (c *APIClient) SetRuleMode(mode string) error {
	result := &model.Response[any]{}
	resp, err := c.client.R().
		SetResult(result).
		Put(fmt.Sprintf("/api/rule/mode/%s", mode))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}

	if result.Code == 0 {
		fmt.Printf("Routing mode set to: ")
		switch mode {
		case "ModeDirect":
			green.Println("Direct")
		case "ModeGlobal":
			orange.Println("Global")
		case "ModeRule":
			white.Println("Rule")
		default:
			white.Println(mode)
		}
	} else {
		fmt.Print("Set routing mode failed: ")
		red.Println(result.Message)
	}

	return nil
}
