package apis

import (
	"fmt"
	"strings"

	"github.com/sipt/shuttle/controller/model"

	apipkg "github.com/sipt/shuttle/group/api"
)

type GroupServer struct {
	Name     string `json:"name"`
	Typ      string `json:"typ"`
	RTT      string `json:"rtt"`
	Selected bool   `json:"selected"`
}

// 服务器编号映射，用于存储编号到服务器名称的映射
var serverIndexMap = make(map[string]string)

// 生成服务器编号 a-z, aa-zz, aaa-zzz 等
func generateServerIndex(index int) string {
	if index < 0 {
		return ""
	}

	result := ""
	for {
		result = string(rune('a'+index%26)) + result
		index = index/26 - 1
		if index < 0 {
			break
		}
	}
	return result
}

// 根据编号获取服务器名称
func GetServerNameByIndex(index string) (string, bool) {
	name, exists := serverIndexMap[index]
	return name, exists
}

// 清空服务器编号映射
func ClearServerIndexMap() {
	serverIndexMap = make(map[string]string)
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

	// 清空之前的映射并重新建立
	ClearServerIndexMap()

	for i, server := range group.Servers {
		serverIndex := generateServerIndex(i)
		serverIndexMap[serverIndex] = server.Name

		fmt.Print("    ")
		if server.Selected {
			green.Print("✔")
		} else {
			white.Print(" ")
		}
		fmt.Printf(" [%s] %s - ", serverIndex, server.Name)
		PrintRtt(server.RTT)
		fmt.Println()
	}
	return nil
}

// 选择服务器组中的服务器
func (c *APIClient) SelectGroupServer(group, server string) error {
	// 检查是否是编号选择
	actualServerName := server
	if serverName, exists := GetServerNameByIndex(strings.ToLower(server)); exists {
		actualServerName = serverName
		fmt.Printf("Selecting server by index [%s]: %s\n", server, actualServerName)
	} else {
		fmt.Println("Selecting server...", group, server)
	}

	result := &model.Response[apipkg.Group]{}
	resp, err := c.client.R().
		SetResult(result).
		SetQueryParams(map[string]string{"group": group, "server": actualServerName}).
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

	// 清空之前的映射并重新建立
	ClearServerIndexMap()

	for i, server := range group.Servers {
		serverIndex := generateServerIndex(i)
		serverIndexMap[serverIndex] = server.Name

		fmt.Print("    ")
		if server.Selected {
			green.Print("✔")
		} else {
			white.Print(" ")
		}
		fmt.Printf(" [%s] %s - ", serverIndex, server.Name)
		PrintRtt(server.RTT)
		fmt.Println()
	}
	return nil
}
