# Shuttle CLI 使用说明

这是一个基于 go-prompt 的交互式命令行工具，用于快速请求 Shuttle API。

## 功能特性

- 交互式命令行界面
- Tab 键自动补全
- 支持多个快速命令
- 友好的错误提示
- 可配置的服务器地址

## 安装和运行

### 编译并运行
```bash
go build
./shuttle-cli
```

### 直接运行
```bash
go run main.go
```

### 指定服务器地址
```bash
# 通过命令行参数
go run main.go -host=192.168.1.100:8080

# 通过环境变量
export SHUTTLE_CTL_HOST=192.168.1.100:8080
go run main.go
```

## 可用命令

| 命令 | 描述 | API 端点 |
|------|------|----------|
| `status` | 获取系统状态 | GET /api/status |
| `inbounds` | 获取入站监听器列表 | GET /api/inbounds |
| `reload` | 重新加载配置 | PUT /api/reload |
| `help` | 显示帮助信息 | - |
| `exit` / `quit` | 退出程序 | - |

## 使用示例

```bash
# 启动 CLI
$ go run main.go
Shuttle CLI - 连接到 localhost:8080
输入 'help' 查看可用命令，或直接输入命令
使用 Tab 键自动补全命令

# 查看帮助
shuttle> help

可用命令:
  status    - 获取系统状态
  inbounds  - 获取入站监听器列表
  reload    - 重新加载配置
  help      - 显示此帮助信息
  exit/quit - 退出程序

当前连接到: localhost:8080

# 获取系统状态
shuttle> status
正在获取系统状态...
系统状态: running

# 获取入站监听器列表
shuttle> inbounds
正在获取入站监听器列表...
入站监听器列表:
  名称: http-proxy, 类型: http, 地址: :8080

# 重新加载配置
shuttle> reload
正在重新加载配置...
配置重新加载成功: success

# 退出程序
shuttle> exit
再见!
```

## 特性说明

### 自动补全
- 使用 Tab 键可以自动补全命令
- 支持模糊匹配

### 错误处理
- 网络错误会显示友好的错误信息
- HTTP 状态码错误会被正确捕获
- API 响应错误会显示详细信息

### 配置
- 默认连接到 `localhost:8080`
- 可通过 `-host` 参数指定其他地址
- 支持 `SHUTTLE_CTL_HOST` 环境变量

## 项目结构

```
shuttle-cli/
├── main.go           # 主程序入口，交互式界面
├── apis/
│   └── api.go       # API 客户端实现
├── go.mod           # Go 模块依赖
├── API.md           # Shuttle API 文档
└── CLI_USAGE.md     # 本使用说明
```

## 依赖库

- `github.com/c-bata/go-prompt` - 交互式命令行界面
- `resty.dev/v3` - HTTP 客户端库

## 扩展指南

要添加新的 API 命令：

1. 在 `apis/api.go` 中添加新的方法
2. 在 `main.go` 的 `suggestions` 中添加命令建议
3. 在 `executor` 函数中添加命令处理逻辑
4. 在 `showHelp` 函数中添加帮助信息

示例：
```go
// 在 apis/api.go 中添加
func (c *APIClient) GetServers() error {
    // 实现逻辑
}

// 在 main.go 中添加
var suggestions = []prompt.Suggest{
    // 现有命令...
    {Text: "servers", Description: "获取服务器列表"},
}

// 在 executor 中添加
case "servers":
    if err := apiClient.GetServers(); err != nil {
        fmt.Printf("错误: %v\n", err)
    }
```
