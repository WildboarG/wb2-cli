package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"wb2-cli/internal/config"
)

// Generator 项目生成器
type Generator struct {
	sdkPath string
}

// New 创建新的生成器实例
func New(sdkPath string) *Generator {
	return &Generator{
		sdkPath: sdkPath,
	}
}

// ProjectData 传递给模板的数据结构
type ProjectData struct {
	ProjectName    string
	SDKPath        string
	Components     []config.Component
	HasWifi        bool
	HasMQTT        bool
	HasBLE         bool
	HasHTTP        bool
	HasStorage     bool
	HasGPIO        bool
	HasUART        bool
	HasI2C         bool
	HasSPI         bool
	HasPWM         bool
	HasADC         bool
	HasTimer       bool
	HasSmartconfig bool
	HasBlufi       bool
	HasHTTPS       bool
	HasLWIPTLS     bool
	IncludeComps   []string
	NetworkComps   []string
	BLSysComps     []string
	VFSComps       []string
	MQTTComps      []string
	ConfigFlags    map[string]string
}

// GenerateProject 生成项目
func (g *Generator) GenerateProject(projectName, projectPath string, components []config.Component) error {
	// 创建项目目录
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("创建项目目录失败: %v", err)
	}

	// 创建项目子目录（与项目同名）
	projectSubDir := filepath.Join(projectPath, projectName)
	if err := os.MkdirAll(projectSubDir, 0755); err != nil {
		return fmt.Errorf("创建项目子目录失败: %v", err)
	}

	// 准备模板数据
	data := g.prepareProjectData(projectName, components)

	// 生成基础文件
	if err := g.generateBaseFiles(projectPath, data); err != nil {
		return err
	}

	// 生成组件相关文件
	if err := g.generateComponentFiles(projectSubDir, data); err != nil {
		return err
	}

	return nil
}

func (g *Generator) prepareProjectData(projectName string, components []config.Component) *ProjectData {
	data := &ProjectData{
		ProjectName:  projectName,
		SDKPath:      g.sdkPath,
		Components:   components,
		IncludeComps: []string{},
		NetworkComps: []string{},
		BLSysComps:   []string{},
		VFSComps:     []string{},
		MQTTComps:    []string{},
		ConfigFlags:  make(map[string]string),
	}

	// 基础组件（所有项目都需要）
	baseIncludeComps := []string{
		"freertos_riscv_ram", "bl602", "bl602_std", "newlibc", "hosal",
		"mbedtls_lts", "lwip", "vfs", "yloop", "utils", "cli",
		"blog", "blog_testc", "coredump",
	}
	data.IncludeComps = append(data.IncludeComps, baseIncludeComps...)

	// 基础 BLSYS 组件
	baseBLSysComps := []string{"bltime", "blfdt", "blmtd", "bloop", "looprt", "loopset"}
	data.BLSysComps = append(data.BLSysComps, baseBLSysComps...)

	// 基础 VFS 组件
	baseVFSComps := []string{"romfs"}
	data.VFSComps = append(data.VFSComps, baseVFSComps...)

	// 处理每个组件
	for _, comp := range components {
		// 检查组件类型
		compNameLower := strings.ToLower(comp.Name)
		if strings.Contains(compNameLower, "wifi") {
			data.HasWifi = true
		}
		if strings.Contains(compNameLower, "mqtt") {
			data.HasMQTT = true
		}
		if strings.Contains(compNameLower, "ble") || strings.Contains(compNameLower, "bluetooth") {
			data.HasBLE = true
		}
		if strings.Contains(compNameLower, "http") {
			data.HasHTTP = true
		}
		if strings.Contains(compNameLower, "storage") || strings.Contains(compNameLower, "flash") {
			data.HasStorage = true
		}
		if strings.Contains(compNameLower, "gpio") {
			data.HasGPIO = true
		}
		if strings.Contains(compNameLower, "uart") {
			data.HasUART = true
		}
		if strings.Contains(compNameLower, "i2c") {
			data.HasI2C = true
		}
		if strings.Contains(compNameLower, "spi") {
			data.HasSPI = true
		}
		if strings.Contains(compNameLower, "pwm") {
			data.HasPWM = true
		}
		if strings.Contains(compNameLower, "adc") {
			data.HasADC = true
		}
		if strings.Contains(compNameLower, "timer") {
			data.HasTimer = true
		}
		if strings.Contains(compNameLower, "smartconfig") {
			data.HasSmartconfig = true
		}
		if strings.Contains(compNameLower, "blufi") {
			data.HasBlufi = true
		}
		if strings.Contains(compNameLower, "https") && !strings.Contains(compNameLower, "lwip") {
			data.HasHTTPS = true
		}
		if strings.Contains(compNameLower, "lwip_tls") || strings.Contains(compNameLower, "lwip_tls") {
			data.HasLWIPTLS = true
		}

		// 合并组件配置
		data.IncludeComps = append(data.IncludeComps, comp.IncludeComponents...)
		data.NetworkComps = append(data.NetworkComps, comp.NetworkComponents...)
		data.BLSysComps = append(data.BLSysComps, comp.BLSysComponents...)
		data.VFSComps = append(data.VFSComps, comp.VFSComponents...)
		data.MQTTComps = append(data.MQTTComps, comp.MQTTComponents...)

		// 合并配置标志
		for k, v := range comp.ConfigFlags {
			data.ConfigFlags[k] = v
		}
	}

	// 如果有 Wi-Fi，添加必要的组件
	if data.HasWifi {
		wifiComps := []string{
			"wifi", "wifi_manager", "wpa_supplicant", "bl_os_adapter",
			"wifi_hosal", "lwip_dhcpd", "netutils", "blcrypto_suite",
			"rfparam_adapter_tmp",
		}
		data.IncludeComps = append(data.IncludeComps, wifiComps...)
		data.NetworkComps = append(data.NetworkComps, "sntp", "dns_server")
		data.BLSysComps = append(data.BLSysComps, "blota", "loopadc")
		data.ConfigFlags["CONFIG_WIFI"] = "1"
	}

	// 如果有 SmartConfig，添加 smartconfig_airkiss 组件
	if data.HasSmartconfig {
		data.IncludeComps = append(data.IncludeComps, "smartconfig_airkiss")
	}

	// 如果有 MQTT，添加必要的组件
	if data.HasMQTT {
		data.IncludeComps = append(data.IncludeComps, "httpc")
	}

	// 去重
	data.IncludeComps = uniqueStrings(data.IncludeComps)
	data.NetworkComps = uniqueStrings(data.NetworkComps)
	data.BLSysComps = uniqueStrings(data.BLSysComps)
	data.VFSComps = uniqueStrings(data.VFSComps)
	data.MQTTComps = uniqueStrings(data.MQTTComps)

	return data
}

func (g *Generator) generateBaseFiles(projectPath string, data *ProjectData) error {
	// 生成 Makefile
	if err := g.generateFileFromTemplate(
		"Makefile.tmpl",
		filepath.Join(projectPath, "Makefile"),
		data,
	); err != nil {
		return fmt.Errorf("生成 Makefile 失败: %v", err)
	}

	// 生成 proj_config.mk
	if err := g.generateFileFromTemplate(
		"proj_config.mk.tmpl",
		filepath.Join(projectPath, "proj_config.mk"),
		data,
	); err != nil {
		return fmt.Errorf("生成 proj_config.mk 失败: %v", err)
	}

	// 生成 README.md
	if err := g.generateFileFromTemplate(
		"README.md.tmpl",
		filepath.Join(projectPath, "README.md"),
		data,
	); err != nil {
		return fmt.Errorf("生成 README.md 失败: %v", err)
	}

	return nil
}

func (g *Generator) generateComponentFiles(projectSubDir string, data *ProjectData) error {
	// 创建 include 目录
	includeDir := filepath.Join(projectSubDir, "include")
	if err := os.MkdirAll(includeDir, 0755); err != nil {
		return fmt.Errorf("创建 include 目录失败: %v", err)
	}

	// 生成 main.c
	if err := g.generateFileFromTemplate(
		"main.c.tmpl",
		filepath.Join(projectSubDir, "main.c"),
		data,
	); err != nil {
		return fmt.Errorf("生成 main.c 失败: %v", err)
	}

	// 生成 bouffalo.mk
	if err := g.generateFileFromTemplate(
		"bouffalo.mk.tmpl",
		filepath.Join(projectSubDir, "bouffalo.mk"),
		data,
	); err != nil {
		return fmt.Errorf("生成 bouffalo.mk 失败: %v", err)
	}

	// 生成 main_board.h
	if err := g.generateFileFromTemplate(
		"main_board.h.tmpl",
		filepath.Join(includeDir, "main_board.h"),
		data,
	); err != nil {
		return fmt.Errorf("生成 main_board.h 失败: %v", err)
	}

	// 根据组件生成特定文件
	for _, comp := range data.Components {
		if err := g.generateComponentSpecificFiles(projectSubDir, comp, data); err != nil {
			return fmt.Errorf("生成组件 %s 文件失败: %v", comp.Name, err)
		}
	}

	return nil
}

func (g *Generator) generateComponentSpecificFiles(projectSubDir string, comp config.Component, data *ProjectData) error {
	// 根据组件的 template_files 配置生成特定文件
	if len(comp.TemplateFiles) == 0 {
		return nil
	}

	for _, tmplFile := range comp.TemplateFiles {
		// 获取模板文件路径
		tmplPath := g.getTemplatePath(tmplFile)
		if tmplPath == "" {
			// 如果模板文件不存在，跳过（不报错，因为某些组件可能不需要特定模板）
			continue
		}

		// 确定输出文件路径
		// template_files 格式：component_name/file.c.tmpl
		// 输出格式：component_name/file.c
		outputFileName := strings.TrimSuffix(filepath.Base(tmplFile), ".tmpl")
		outputDir := filepath.Join(projectSubDir, filepath.Dir(tmplFile))
		
		// 创建输出目录
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("创建组件目录失败: %v", err)
		}

		outputPath := filepath.Join(outputDir, outputFileName)

		// 生成文件
		if err := g.generateFileFromTemplate(tmplFile, outputPath, data); err != nil {
			// 如果模板文件不存在，只记录警告，不中断流程
			fmt.Printf("警告: 无法生成组件文件 %s: %v\n", outputPath, err)
			continue
		}
	}

	return nil
}

func (g *Generator) generateFileFromTemplate(templateName, outputPath string, data interface{}) error {
	// 获取模板文件路径
	tmplPath := g.getTemplatePath(templateName)
	if tmplPath == "" {
		return fmt.Errorf("找不到模板文件: %s", templateName)
	}

	// 读取模板内容
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %v", err)
	}

	// 解析模板
	tmpl, err := template.New(templateName).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer file.Close()

	// 渲染模板
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("渲染模板失败: %v", err)
	}

	return nil
}

func (g *Generator) getTemplatePath(templateName string) string {
	// 尝试从多个位置查找模板文件
	var possiblePaths []string

	// 1. 从当前工作目录查找（开发模式）
	cwd, _ := os.Getwd()
	possiblePaths = append(possiblePaths,
		filepath.Join(cwd, "internal", "generator", "templates", templateName),
		filepath.Join(cwd, "wb2-cli", "internal", "generator", "templates", templateName),
	)

	// 2. 从可执行文件所在目录查找（安装模式）
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		possiblePaths = append(possiblePaths,
			filepath.Join(exeDir, "templates", templateName),
			filepath.Join(exeDir, "internal", "generator", "templates", templateName),
		)
	}

	// 3. 尝试从源码目录查找（如果从 SDK 根目录运行）
	if cwd != "" {
		possiblePaths = append(possiblePaths,
			filepath.Join(cwd, "wb2-cli", "internal", "generator", "templates", templateName),
		)
	}

	// 查找第一个存在的文件
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, s := range strs {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
