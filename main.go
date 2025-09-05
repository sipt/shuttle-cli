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

// Navigation state
type NavigationState struct {
	Path      []string // Current path segments ["groups", "IEPL"]
	Context   string   // Current context: "root", "groups", "group"
	GroupName string   // Current group name when in group context
}

var navState = &NavigationState{
	Path:    []string{},
	Context: "root",
}

// 命令建议 - 根据上下文动态生成
var rootSuggestions = []prompt.Suggest{
	{Text: "status", Description: "Get system status"},
	{Text: "inbounds", Description: "Get inbound listeners list"},
	{Text: "reload", Description: "Reload configuration"},
	{Text: "servers", Description: "Get server list"},
	{Text: "server", Description: "Get specific server info, usage: server <name>"},
	{Text: "server-rtt", Description: "Test server RTT, usage: server-rtt <name>"},
	{Text: "cd", Description: "Change directory, usage: cd <path>"},
	{Text: "ll", Description: "List current directory contents"},
	{Text: "pwd", Description: "Print working directory"},
	{Text: "help", Description: "Show help information"},
	{Text: "exit", Description: "Exit program"},
	{Text: "quit", Description: "Exit program"},
}

var groupsSuggestions = []prompt.Suggest{
	{Text: "ll", Description: "List all groups"},
	{Text: "cd", Description: "Change to group directory, usage: cd <group-name>"},
	{Text: "pwd", Description: "Print working directory"},
	{Text: "help", Description: "Show help information"},
	{Text: "..", Description: "Go back to parent directory"},
	{Text: "exit", Description: "Exit program"},
	{Text: "quit", Description: "Exit program"},
}

var groupSuggestions = []prompt.Suggest{
	{Text: "ll", Description: "List servers in current group"},
	{Text: "rtt", Description: "Test RTT for all servers in group"},
	{Text: "select", Description: "Select server by name or index (a, b, c...)"},
	{Text: "cd", Description: "Change directory, usage: cd <path>"},
	{Text: "pwd", Description: "Print working directory"},
	{Text: "help", Description: "Show help information"},
	{Text: "..", Description: "Go back to parent directory"},
	{Text: "exit", Description: "Exit program"},
	{Text: "quit", Description: "Exit program"},
}

func completer(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest
	switch navState.Context {
	case "root":
		suggestions = rootSuggestions
	case "groups":
		suggestions = groupsSuggestions
	case "group":
		suggestions = groupSuggestions
	default:
		suggestions = rootSuggestions
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

// Helper functions for navigation
func getCurrentPath() string {
	if len(navState.Path) == 0 {
		return "shuttle"
	}
	return "shuttle/" + strings.Join(navState.Path, "/")
}

func getCurrentPrompt() string {
	return getCurrentPath() + "> "
}

func navigateTo(path string) error {
	if path == ".." {
		if len(navState.Path) > 0 {
			navState.Path = navState.Path[:len(navState.Path)-1]
		}
		updateContext()
		return nil
	}

	if path == "" || path == "/" {
		navState.Path = []string{}
		updateContext()
		return nil
	}

	// Handle absolute paths
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
		navState.Path = []string{}
	}

	if path != "" {
		parts := strings.Split(path, "/")
		for _, part := range parts {
			if part == ".." {
				if len(navState.Path) > 0 {
					navState.Path = navState.Path[:len(navState.Path)-1]
				}
			} else if part != "" {
				navState.Path = append(navState.Path, part)
			}
		}
	}

	updateContext()
	return nil
}

func updateContext() {
	if len(navState.Path) == 0 {
		navState.Context = "root"
		navState.GroupName = ""
	} else if len(navState.Path) == 1 && navState.Path[0] == "groups" {
		navState.Context = "groups"
		navState.GroupName = ""
	} else if len(navState.Path) == 2 && navState.Path[0] == "groups" {
		navState.Context = "group"
		navState.GroupName = navState.Path[1]
	} else {
		// Invalid path, reset to root
		navState.Path = []string{}
		navState.Context = "root"
		navState.GroupName = ""
	}
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

	// Handle navigation commands first
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			navigateToRoot()
		} else {
			if err := navigateTo(args[1]); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
		return
	case "ll":
		handleListCommand()
		return
	case "pwd":
		fmt.Println(getCurrentPath())
		return
	case "rtt":
		if navState.Context == "group" {
			if err := apiClient.TestGroupRTT(navState.GroupName); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Println("RTT command only available in group context")
		}
		return
	case "select":
		if navState.Context == "group" {
			if len(args) < 2 {
				fmt.Println("Usage: select <server-name-or-index>")
				fmt.Println("  Examples: select a, select b, select \"My Server Name\"")
				return
			}
			server := strings.Join(args[1:], " ")
			if err := apiClient.SelectGroupServer(navState.GroupName, server); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Println("Select command only available in group context")
		}
		return
	}

	// Handle legacy commands
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

func navigateToRoot() {
	navState.Path = []string{}
	updateContext()
}

func handleListCommand() {
	switch navState.Context {
	case "root":
		fmt.Println("Available directories:")
		fmt.Println("  groups/")
		fmt.Println("\nAvailable commands:")
		fmt.Println("  status, inbounds, reload, servers")
	case "groups":
		if err := apiClient.GetGroups(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "group":
		if err := apiClient.GetGroup(navState.GroupName); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func showHelp() {
	fmt.Printf("\nCurrent location: %s\n", getCurrentPath())
	fmt.Println("\nNavigation commands:")
	fmt.Println("  cd <path>             - Change directory (cd groups, cd .., cd /)")
	fmt.Println("  ll                    - List current directory contents")
	fmt.Println("  pwd                   - Print working directory")

	switch navState.Context {
	case "root":
		fmt.Println("\nAvailable commands:")
		fmt.Println("  status                - Get system status")
		fmt.Println("  inbounds              - Get inbound listeners list")
		fmt.Println("  reload                - Reload configuration")
		fmt.Println("  servers               - Get server list")
		fmt.Println("  server <name>         - Get specific server info")
		fmt.Println("  server-rtt <name>     - Test server RTT")
		fmt.Println("\nAvailable directories:")
		fmt.Println("  groups/               - Server groups management")
	case "groups":
		fmt.Println("\nAvailable commands:")
		fmt.Println("  ll                    - List all groups")
		fmt.Println("  cd <group-name>       - Enter specific group")
		fmt.Println("  cd ..                 - Go back to root")
	case "group":
		fmt.Printf("\nGroup: %s\n", navState.GroupName)
		fmt.Println("Available commands:")
		fmt.Println("  ll                    - List servers in group")
		fmt.Println("  rtt                   - Test RTT for all servers")
		fmt.Println("  select <server>       - Select server by name or index (a, b, c...)")
		fmt.Println("  cd ..                 - Go back to groups")
	}

	fmt.Println("\nGeneral commands:")
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
	fmt.Println("Use 'cd groups' to navigate to groups, 'll' to list, 'pwd' to show current path")

	// 启动交互式 prompt
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("shuttle-cli"),
		prompt.OptionLivePrefix(func() (string, bool) {
			return getCurrentPrompt(), true
		}),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(" "),
	)
	p.Run()
}
