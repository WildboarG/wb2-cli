package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"wb2-cli/internal/config"
	"wb2-cli/internal/generator"
)

var (
	projectPath string
	interactive bool
)

// clearScreen 跨平台清屏函数
func clearScreen() {
	fmt.Print("\033[2J\033[H") // ANSI转义序列，在现代终端中都支持
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "创建一个新的 WB2 项目",
	Long: `创建一个新的 WB2 项目，支持交互式选择组件。

示例:
  wb2-cli new my_project
  wb2-cli new my_project --path ./projects
  wb2-cli new my_project --sdk-path /path/to/sdk`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "项目创建路径（默认为当前目录）")
	newCmd.Flags().BoolVarP(&interactive, "interactive", "i", true, "交互式选择组件（默认启用）")
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// 验证项目名称
	if !isValidProjectName(projectName) {
		return fmt.Errorf("无效的项目名称: %s (只能包含字母、数字、下划线和连字符)", projectName)
	}

	// 获取 SDK 路径
	sdkPath, err := getSDKPath()
	if err != nil {
		return fmt.Errorf("获取 SDK 路径失败: %v", err)
	}

	// 验证 SDK 路径
	if !isValidSDKPath(sdkPath) {
		return fmt.Errorf("无效的 SDK 路径: %s", sdkPath)
	}

	// 加载组件配置
	components, err := config.LoadComponents()
	if err != nil {
		return fmt.Errorf("加载组件配置失败: %v", err)
	}

	// 交互式选择组件
	selectedComponents, err := selectComponents(components)
	if err != nil {
		return fmt.Errorf("选择组件失败: %v", err)
	}

	// 解析组件依赖
	resolvedComponents, err := resolveDependencies(components, selectedComponents)
	if err != nil {
		return fmt.Errorf("解析组件依赖失败: %v", err)
	}

	// 生成项目路径
	fullProjectPath := filepath.Join(projectPath, projectName)

	// 检查目录是否已存在
	if _, err := os.Stat(fullProjectPath); err == nil {
		return fmt.Errorf("项目目录已存在: %s", fullProjectPath)
	}

	// 创建项目
	gen := generator.New(sdkPath)
	err = gen.GenerateProject(projectName, fullProjectPath, resolvedComponents)
	if err != nil {
		return fmt.Errorf("生成项目失败: %v", err)
	}

	fmt.Printf("\n✅ 项目创建成功！\n")
	fmt.Printf("📁 项目路径: %s\n", fullProjectPath)
	fmt.Printf("📦 已选择组件: %s\n", strings.Join(selectedComponents, ", "))
	fmt.Printf("\n下一步:\n")
	fmt.Printf("  cd %s\n", fullProjectPath)
	fmt.Printf("  make -j8\n")

	return nil
}

func isValidProjectName(name string) bool {
	// 只允许字母、数字、下划线和连字符
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-') {
			return false
		}
	}
	return len(name) > 0
}

func getSDKPath() (string, error) {
	// 如果命令行指定了 SDK 路径，优先使用
	if sdkPath != "" {
		return sdkPath, nil
	}

	// 从配置文件读取
	cfg, err := config.LoadConfig()
	if err != nil {
		// 如果配置文件不存在，尝试自动检测
		return autoDetectSDKPath()
	}

	if cfg.SDKPath != "" {
		return cfg.SDKPath, nil
	}

	// 如果配置文件中也没有，尝试自动检测
	return autoDetectSDKPath()
}

func autoDetectSDKPath() (string, error) {
	// 尝试从当前工作目录向上查找 SDK
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 检查当前目录是否是 SDK 根目录
	if isValidSDKPath(cwd) {
		return cwd, nil
	}

	// 检查父目录
	parent := filepath.Dir(cwd)
	if isValidSDKPath(parent) {
		return parent, nil
	}

	return "", fmt.Errorf("无法自动检测 SDK 路径，请使用 --sdk-path 参数指定")
}

func isValidSDKPath(path string) bool {
	// 检查是否存在必要的 SDK 目录和文件
	requiredPaths := []string{
		"components",
		"applications",
		"make_scripts_riscv",
		"version.mk",
	}

	for _, req := range requiredPaths {
		fullPath := filepath.Join(path, req)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// selectComponentsWindows Windows版本的组件选择（简化版）
func selectComponentsWindows(allComponents []config.Component) ([]string, error) {
	fmt.Println("🌟 wb2-cli - 组件选择器")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("在Windows环境下，我们提供简化的组件选择方式。")
	fmt.Println("您可以输入组件名称（多个用逗号分隔），或者输入'all'选择所有组件。")
	fmt.Println()

	// 按分类显示可用组件
	categories := make(map[string][]config.Component)
	for _, comp := range allComponents {
		categories[comp.Category] = append(categories[comp.Category], comp)
	}

	for category, comps := range categories {
		fmt.Printf("📁 %s:\n", category)
		for _, comp := range comps {
			fmt.Printf("  - %s: %s\n", comp.Name, comp.Description)
		}
		fmt.Println()
	}

	fmt.Print("请输入要选择的组件（用逗号分隔，或输入'all'选择全部，或按回车跳过）: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}, nil
	}

	if input == "all" {
		var allNames []string
		for _, comp := range allComponents {
			allNames = append(allNames, comp.Name)
		}
		return allNames, nil
	}

	// 解析用户输入的组件名称
	selectedNames := strings.Split(input, ",")
	var validSelections []string

	for _, name := range selectedNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		// 检查组件是否存在
		found := false
		for _, comp := range allComponents {
			if comp.Name == name {
				validSelections = append(validSelections, name)
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("⚠️  警告: 组件 '%s' 不存在，已跳过\n", name)
		}
	}

	return validSelections, nil
}

func selectComponents(allComponents []config.Component) ([]string, error) {
	if !interactive {
		// 非交互模式，返回空列表（只包含基础组件）
		return []string{}, nil
	}

	// 根据操作系统选择不同的交互方式
	if runtime.GOOS == "windows" {
		return selectComponentsWindows(allComponents)
	}

	// Unix/Linux 版本使用原始终端交互

	// 按分类组织组件
	componentsByCategory := make(map[string][]config.Component)
	categoryNames := map[string]string{
		"network":    "🌐 网络组件",
		"peripheral": "🔌 外设组件",
		"3rdparty":   "📦 第三方组件",
		"audio":      "🔊 音频组件",
		"fs":         "💾 文件系统组件",
		"multimedia": "🎬 多媒体组件",
		"system":     "⚙️  系统组件",
		"other":      "📋 其他组件",
	}

	// 将组件按分类分组
	for _, comp := range allComponents {
		category := comp.Category
		if category == "" {
			category = "other"
		}
		componentsByCategory[category] = append(componentsByCategory[category], comp)
	}

	// 已选择的组件集合
	selectedSet := make(map[string]bool)
	categoryOrder := []string{"network", "peripheral", "3rdparty", "audio", "fs", "multimedia", "system", "other"}

	// 当前菜单状态
	currentCategory := ""
	selectedIndex := 0
	componentIndex := 0

	// 计算有效分类列表（只计算一次，避免重复计算）
	validCategories := []string{}
	for _, cat := range categoryOrder {
		if len(componentsByCategory[cat]) > 0 {
			validCategories = append(validCategories, cat)
		}
	}

	// 主循环
	for {
		// 清屏
		clearScreen()

		if currentCategory == "" {
			// 显示主菜单（分类列表）
			fmt.Println(strings.Repeat("=", 70))
			fmt.Println("           WB2 组件选择菜单 (类似 menuconfig)")
			fmt.Println(strings.Repeat("=", 70))
			fmt.Println()

			// 显示已选择的组件
			if len(selectedSet) > 0 {
				fmt.Println("已选择的组件:")
				selectedList := make([]string, 0, len(selectedSet))
				for name := range selectedSet {
					selectedList = append(selectedList, name)
				}
				// 排序
				for i := 0; i < len(selectedList)-1; i++ {
					for j := i + 1; j < len(selectedList); j++ {
						if selectedList[i] > selectedList[j] {
							selectedList[i], selectedList[j] = selectedList[j], selectedList[i]
						}
					}
				}
				for _, name := range selectedList {
					fmt.Printf("  ✓ %s\n", name)
				}
				fmt.Println()
			}

			// 显示分类列表
			fmt.Println("请选择分类:")
			for i, cat := range validCategories {
				comps := componentsByCategory[cat]
				catDisplayName := categoryNames[cat]
				if catDisplayName == "" {
					catDisplayName = cat
				}

				prefix := "  "
				if i == selectedIndex {
					prefix = "> "
				}
				fmt.Printf("%s▶ %s (%d)\n", prefix, catDisplayName, len(comps))
			}
			fmt.Println()
			fmt.Println("操作: ↑↓ 导航 | → 进入 | 回车 完成选择")

		} else {
			// 显示分类内的组件列表
			comps := componentsByCategory[currentCategory]
			catDisplayName := categoryNames[currentCategory]

			fmt.Println(strings.Repeat("=", 70))
			fmt.Printf("  %s\n", catDisplayName)
			fmt.Println(strings.Repeat("=", 70))
			fmt.Println()

			fmt.Println("请选择组件:")
			for i, comp := range comps {
				prefix := "  "
				if i == componentIndex {
					prefix = "> "
				}

				status := " "
				if selectedSet[comp.Name] {
					status = "✓"
				}

				fmt.Printf("%s[%s] %s - %s\n", prefix, status, comp.Name, comp.Description)
			}
			fmt.Println()
			fmt.Println("操作: ↑↓ 导航 | 空格 选择/取消 | ← 返回 | 回车 返回")
		}

		// 读取按键
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		
		if err != nil {
			term.Restore(int(os.Stdin.Fd()), oldState)
			return nil, err
		}

		// 处理按键
		if currentCategory == "" {
			// 主菜单模式
			compsCount := len(validCategories)

			switch char {
			case 27: // ESC 序列
				buf := make([]byte, 2)
				reader.Read(buf)
				term.Restore(int(os.Stdin.Fd()), oldState)
				if buf[0] == '[' {
					switch buf[1] {
					case 'A': // ↑
						if selectedIndex > 0 {
							selectedIndex--
						}
					case 'B': // ↓
						if selectedIndex < compsCount-1 {
							selectedIndex++
						}
					case 'C': // →
						// 进入选中的分类
						if selectedIndex < len(validCategories) {
							currentCategory = validCategories[selectedIndex]
							componentIndex = 0
						}
					}
				}
				continue
			case '\n', '\r': // 回车 - 完成选择
				// 转换为组件名称列表
				selectedComponents := make([]string, 0, len(selectedSet))
				for name := range selectedSet {
					selectedComponents = append(selectedComponents, name)
				}
				term.Restore(int(os.Stdin.Fd()), oldState)
				clearScreen() // 清屏
				return selectedComponents, nil
			case 'q', 'Q':
				// 退出
				term.Restore(int(os.Stdin.Fd()), oldState)
				clearScreen() // 清屏
				return nil, fmt.Errorf("用户取消")
			default:
				term.Restore(int(os.Stdin.Fd()), oldState)
			}
		} else {
			// 组件列表模式
			comps := componentsByCategory[currentCategory]

			switch char {
			case 27: // ESC 序列
				buf := make([]byte, 2)
				reader.Read(buf)
				term.Restore(int(os.Stdin.Fd()), oldState)
				if buf[0] == '[' {
					switch buf[1] {
					case 'A': // ↑
						if componentIndex > 0 {
							componentIndex--
						}
					case 'B': // ↓
						if componentIndex < len(comps)-1 {
							componentIndex++
						}
					case 'D': // ←
						// 返回主菜单
						currentCategory = ""
						selectedIndex = 0
					}
				}
				continue
			case ' ': // 空格 - 切换选择状态
				if componentIndex < len(comps) {
					comp := comps[componentIndex]
					selectedSet[comp.Name] = !selectedSet[comp.Name]
				}
				term.Restore(int(os.Stdin.Fd()), oldState)
				continue
			case '\n', '\r': // 回车 - 返回上一级
				currentCategory = ""
				selectedIndex = 0
				term.Restore(int(os.Stdin.Fd()), oldState)
				continue
			case 'q', 'Q':
				// 退出
				term.Restore(int(os.Stdin.Fd()), oldState)
				clearScreen() // 清屏
				return nil, fmt.Errorf("用户取消")
			default:
				term.Restore(int(os.Stdin.Fd()), oldState)
			}
		}
	}
}

func resolveDependencies(allComponents []config.Component, selected []string) ([]config.Component, error) {
	// 创建组件映射
	componentMap := make(map[string]config.Component)
	for _, comp := range allComponents {
		componentMap[comp.Name] = comp
	}

	// 解析依赖
	resolved := make(map[string]bool)
	toResolve := make([]string, len(selected))
	copy(toResolve, selected)

	for len(toResolve) > 0 {
		compName := toResolve[0]
		toResolve = toResolve[1:]

		if resolved[compName] {
			continue
		}

		comp, ok := componentMap[compName]
		if !ok {
			return nil, fmt.Errorf("未知的组件: %s", compName)
		}

		resolved[compName] = true

		// 添加依赖
		for _, dep := range comp.Dependencies {
			if !resolved[dep] {
				toResolve = append(toResolve, dep)
			}
		}
	}

	// 转换为组件列表
	result := make([]config.Component, 0, len(resolved))
	for name := range resolved {
		result = append(result, componentMap[name])
	}

	return result, nil
}
