package apis

import (
	"fmt"

	"github.com/sipt/shuttle/controller/model"

	apipkg "github.com/sipt/shuttle/group/api"
)

type GroupServer struct {
	Name     string `json:"name"`
	Typ      string `json:"typ"`
	RTT      string `json:"rtt"`
	Selected bool   `json:"selected"`
}

// 获取服务器组列表
func (c *APIClient) GetGroups() error {
	result := &model.Response[[]apipkg.Group]{}
	resp, err := c.client.R().
		SetResult(result).
		Get("/api/groups")
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}
	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}
	fmt.Println("Groups:")
	for _, group := range result.Data {
		fmt.Printf("  %s - %s\n", group.Name, group.Selected.Name)
	}
	return nil
}

// 获取指定服务器组信息
func (c *APIClient) GetGroup(name string) error {
	result := &model.Response[apipkg.Group]{}
	resp, err := c.client.R().
		SetResult(result).
		SetQueryParam("group", name).
		Get("/api/group")
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}
	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}
	group := result.Data
	fmt.Printf("Group: %s (%s)\n", group.Name, group.Typ)
	fmt.Println("  Servers:")
	for _, server := range group.Servers {
		fmt.Print("      ")
		if server.Selected {
			green.Print("✔")
		} else {
			white.Print(" ")
		}
		fmt.Printf(" %s - ", server.Name)
		PrintRtt(server.RTT)
		fmt.Println()
	}
	return nil
}

// 选择服务器组中的服务器
func (c *APIClient) SelectGroupServer(group, server string) error {
	fmt.Println("Selecting server...", group, server)
	result := &model.Response[apipkg.Group]{}
	resp, err := c.client.R().
		SetResult(result).
		SetQueryParams(map[string]string{"group": group, "server": server}).
		Put("/api/group/select")
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}
	if result.Code == 0 {
		fmt.Printf("  %s - %s\n", result.Data.Name, result.Data.Selected.Name)
	} else {
		fmt.Print("Select failed: ")
		red.Println(result.Message)
	}
	return nil
}

// 重置服务器组延迟测试
func (c *APIClient) TestGroupRTT(groupName string) error {
	result := &model.Response[apipkg.Group]{}
	resp, err := c.client.R().
		SetResult(result).
		SetQueryParam("group", groupName).
		Put("/api/group/rtt")
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
	}
	if result.Code != 0 {
		fmt.Print("RTT test failed: ")
		red.Println(result.Message)
	}
	group := result.Data
	fmt.Printf("Group: %s (%s)\n", group.Name, group.Typ)
	fmt.Println("  Servers:")
	for _, server := range group.Servers {
		fmt.Print("      ")
		if server.Selected {
			green.Print("✔")
		}
		fmt.Printf(" %s - ", server.Name)
		PrintRtt(server.RTT)
		fmt.Println()
	}
	return nil
}
