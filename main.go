package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/sipt/shuttle-cli/apis"
)

var host = flag.String("host", os.Getenv("SHUTTLE_CTL_HOST"), "The host of the controller")

var controllerHost = "localhost:8080"
var apiClient *apis.APIClient

// 命令建议
var suggestions = []prompt.Suggest{
	{Text: "status", Description: "Get system status"},
	{Text: "inbounds", Description: "Get inbound listeners list"},
	{Text: "reload", Description: "Reload configuration"},
	{Text: "servers", Description: "Get server list"},
	{Text: "server", Description: "Get specific server info, usage: server <name>"},
	{Text: "server-rtt", Description: "Test server RTT, usage: server-rtt <name>"},
	{Text: "groups", Description: "Get server group list"},
	{Text: "group", Description: "Get group info, usage: group <name>"},
	{Text: "group-select", Description: "Select server in group, usage: group-select <group> <server>"},
	{Text: "group-rtt", Description: "Test group RTT, usage: group-rtt <group>"},
	{Text: "help", Description: "Show help information"},
	{Text: "exit", Description: "Exit program"},
	{Text: "quit", Description: "Exit program"},
}

func completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "status":
		if err := apiClient.GetStatus(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "inbounds":
		if err := apiClient.GetInbounds(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "reload":
		if err := apiClient.Reload(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "servers":
		if err := apiClient.GetServers(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "server":
		if len(args) < 2 {
			fmt.Println("Usage: server <name>")
			return
		}
		if err := apiClient.GetServer(args[1]); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "server-rtt":
		if len(args) < 2 {
			fmt.Println("Usage: server-rtt <name>")
			return
		}
		if err := apiClient.TestServerRTT(args[1]); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "groups":
		if err := apiClient.GetGroups(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "group":
		if len(args) < 2 {
			fmt.Println("Usage: group <name>")
			return
		}
		if err := apiClient.GetGroup(args[1]); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "group-select":
		if len(args) < 3 {
			fmt.Println("Usage: group-select <group> <server>")
			return
		}
		server := strings.Join(args[2:], " ")
		if err := apiClient.SelectGroupServer(args[1], server); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "group-rtt":
		if len(args) < 2 {
			fmt.Println("Usage: group-rtt <group>")
			return
		}
		if err := apiClient.TestGroupRTT(args[1]); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "help":
		showHelp()
	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", input)
		fmt.Println("Type 'help' to see available commands")
	}
}

func showHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  status                - Get system status")
	fmt.Println("  inbounds              - Get inbound listeners list")
	fmt.Println("  reload                - Reload configuration")
	fmt.Println("  servers               - Get server list")
	fmt.Println("  server <name>         - Get specific server info")
	fmt.Println("  server-rtt <name>     - Test server RTT")
	fmt.Println("  groups                - Get server group list")
	fmt.Println("  group <name>          - Get group info")
	fmt.Println("  group-select <g> <s>  - Select server in group")
	fmt.Println("  group-rtt <group>     - Test group RTT")
	fmt.Println("  help                  - Show this help information")
	fmt.Println("  exit/quit             - Exit program")
	fmt.Printf("\nCurrently connected to: %s\n", controllerHost)
	fmt.Println()
}

func main() {
	flag.Parse()
	if host := *host; host != "" {
		controllerHost = host
	}

	// 初始化 API 客户端
	apiClient = apis.NewAPIClient(controllerHost)

	fmt.Printf("Shuttle CLI - Connected to %s\n", controllerHost)
	fmt.Println("Type 'help' to see available commands, or enter commands directly")
	fmt.Println("Use Tab key for command auto-completion")

	// 启动交互式 prompt
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("shuttle-cli"),
		prompt.OptionPrefix("shuttle> "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(" "),
	)
	p.Run()
}
