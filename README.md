# wb2-cli

Ai-Thinker WB2 芯片项目快速生成工具

[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/WildboarG/wb2-cli.svg)](https://github.com/WildboarG/wb2-cli/releases)
[![Build](https://img.shields.io/github/actions/workflow/status/WildboarG/wb2-cli/release.yml)](https://github.com/WildboarG/wb2-cli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/WildboarG/wb2-cli)](https://goreportcard.com/report/github.com/WildboarG/wb2-cli)
[![License](https://img.shields.io/github/license/WildboarG/wb2-cli.svg)](https://github.com/WildboarG/wb2-cli/blob/main/LICENSE)

## 功能特性

- 🚀 快速创建 WB2 项目框架
- 📦 交互式菜单选择组件（类似 menuconfig）
- 🔗 自动解析组件依赖关系
- 📝 自动生成项目文件（Makefile、proj_config.mk、main.c 等）
- ⚙️ 支持自定义 SDK 路径
- 🎯 支持多种 Wi-Fi 配网方式（静态连接、SmartConfig、BluFi）
- 🖥️ 跨平台支持（Linux、Windows）

## 安装

### 从源码编译

**Linux/macOS:**
```bash
cd wb2-cli
go build -o wb2-cli .
```

**Windows:**
```bash
cd wb2-cli
go build -o wb2-cli.exe .
```

将编译好的可执行文件放到系统 PATH 中，或直接使用相对路径运行。

### 🧪 测试

运行测试套件：

```bash
# 运行所有测试
go test ./...

# 带覆盖率报告
go test -cover ./...

# 详细测试输出
go test -v ./...
```

### 交叉编译

您也可以在Linux上为Windows编译：

```bash
cd wb2-cli
GOOS=windows GOARCH=amd64 go build -o wb2-cli.exe .
```

## 快速开始

```bash
# 创建新项目
wb2-cli new my_project

# 指定 SDK 路径（推荐首次使用）
wb2-cli new my_project --sdk-path /path/to/Ai-Thinker-WB2
```

## 组件选择菜单

工具采用类似 `menuconfig` 的交互式菜单，支持键盘导航：

### Linux/macOS 版本

- **主菜单**：使用 ↑↓ 键浏览分类，→ 键进入分类，回车键完成选择
- **组件列表**：空格键选中/取消，← 键返回主菜单
- **导航**：Q 键退出程序

### Windows 版本

Windows环境下使用简化的文本界面：

```
🌟 wb2-cli - 组件选择器
========================

📁 network:
  - wifi: Wi-Fi 连接功能（Station/AP 模式）
  - mqtt: MQTT 客户端功能
  - ...

请输入要选择的组件（用逗号分隔，或输入'all'选择全部，或按回车跳过）:
```

### 支持的组件分类

- 🌐 **网络组件**：Wi-Fi、MQTT、HTTP、BLE、SmartConfig、BluFi 等
- 🔌 **外设组件**：GPIO、UART、I2C、SPI、PWM、ADC、Timer
- 💾 **存储组件**：Flash、EasyFlash、ROMFS、SPIFFS
- 📱 **多媒体组件**：LVGL、JPEG编解码
- 🔧 **系统组件**：cJSON、CLI、日志、OTA 等

## Wi-Fi 配网方式

### 静态连接（仅选择 wifi）

```c
#define ROUTER_SSID "your_wifi_ssid"
#define ROUTER_PWD "your_wifi_password"

static void wifi_sta_connect(char* ssid, char* password) {
    wifi_interface_t wifi_interface = wifi_mgmr_sta_enable();
    wifi_mgmr_sta_connect(wifi_interface, ssid, password, NULL, NULL, 0, 0);
}
```

### SmartConfig 配网（wifi + smartconfig）

```c
#include <smartconfig.h>

// 在 Wi-Fi 就绪事件中
blog_info("Starting smartconfig...");
wifi_smartconfig_v1_start();
```

### BluFi 配网（wifi + ble + blufi）

```c
#include <blufi.h>

// BLE 栈会自动处理 BluFi 配网
```

## 项目结构

```
my_project/
├── Makefile              # 项目构建文件
├── proj_config.mk        # 项目配置文件
├── README.md             # 项目说明文件
└── my_project/           # 源代码目录
    ├── main.c            # 主程序入口
    ├── bouffalo.mk       # 组件构建配置
    └── include/
        └── main_board.h  # 硬件配置头文件
```

## 编译和烧录

```bash
cd my_project

# 编译项目
make -j8

# 烧录到开发板
make flash p=/dev/ttyUSB0 b=921600
```

## SDK 路径配置

工具按以下优先级查找 SDK：

1. 命令行参数 `--sdk-path`
2. 配置文件 `~/.config/wb2-cli/config.yaml`
3. 自动检测（向上查找目录）

## 添加新组件

### 1. 编辑组件配置

在 `assets/components.yaml` 中添加组件定义：

```yaml
- name: my_component
  description: 我的组件描述
  category: network  # 分类：network, peripheral, storage, multimedia, system
  dependencies:      # 依赖组件（可选）
    - wifi
  sdk_components:    # SDK 组件列表
    - component1
    - component2
  config_flags:      # 配置标志（可选）
    CONFIG_MY_FLAG: "1"
```

### 2. 添加模板文件（可选）

如果组件需要生成特定代码，在 `internal/generator/templates/components/` 下创建模板文件。

### 3. 更新生成逻辑

在 `internal/generator/generator.go` 中添加组件特定的生成逻辑。

## 开发说明

### 项目架构

```c
wb2-cli/
├── cmd/                  # CLI 命令定义
├── internal/
│   ├── config/          # 组件配置管理
│   └── generator/       # 项目文件生成器
│       └── templates/   # 模板文件
├── assets/
│   └── components.yaml  # 组件定义文件
└── main.go
```

### 模板系统

- **主模板**：`main.c.tmpl` - 生成主程序文件
- **构建模板**：`Makefile.tmpl`, `proj_config.mk.tmpl` - 生成构建配置
- **组件模板**：`components/` 目录下的特定组件模板

## 发布新版本

### 创建 Release

1. 更新 `CHANGELOG.md` 文件，添加新版本的变更记录
2. 提交更改并推送
3. 在 GitHub 上创建新的 tag：

```bash
# 创建带注解的tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# 推送tag到GitHub
git push origin v1.0.0
```

4. GitHub Actions 会自动：
   - 编译所有平台的二进制文件
   - 生成源码包和校验文件
   - 创建 GitHub Release 并上传所有文件

### 版本号规范

项目遵循 [Semantic Versioning](https://semver.org/)：

- **MAJOR.MINOR.PATCH** (例如: 1.0.0)
- **MAJOR**: 不兼容的 API 变更
- **MINOR**: 向后兼容的新功能
- **PATCH**: 向后兼容的 bug 修复

## 开发

### 本地测试发布流程

```bash
# 运行测试
go test -v ./...

# 交叉编译测试
GOOS=linux GOARCH=amd64 go build -o wb2-cli-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o wb2-cli-windows-amd64.exe .
GOOS=darwin GOARCH=amd64 go build -o wb2-cli-darwin-amd64 .
```

### 代码质量

- 使用 `go vet` 检查代码
- 使用 `go fmt` 格式化代码
- 运行测试覆盖率检查

## 许可证

遵循与 Ai-Thinker-WB2 SDK 相同的许可证。

## 贡献

欢迎提交 Issue 和 Pull Request！

### 贡献者指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

请确保：
- 通过所有测试
- 更新相关文档
- 遵循现有的代码风格
