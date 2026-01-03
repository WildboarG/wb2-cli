package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Component 表示一个可用的组件
type Component struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Category     string   `yaml:"category,omitempty"` // 组件分类（如：network, peripheral, 3rdparty 等）
	Dependencies []string `yaml:"dependencies,omitempty"`
	// 组件在 SDK 中的路径（用于 Makefile）
	SDKComponents []string `yaml:"sdk_components,omitempty"`
	// 需要添加到 INCLUDE_COMPONENTS 的组件
	IncludeComponents []string `yaml:"include_components,omitempty"`
	// 需要添加到 COMPONENTS_NETWORK 的组件
	NetworkComponents []string `yaml:"network_components,omitempty"`
	// 需要添加到 COMPONENTS_BLSYS 的组件
	BLSysComponents []string `yaml:"blsys_components,omitempty"`
	// 需要添加到 COMPONENTS_VFS 的组件
	VFSComponents []string `yaml:"vfs_components,omitempty"`
	// 需要添加到 COMPONENTS_MQTT 的组件（如果有）
	MQTTComponents []string `yaml:"mqtt_components,omitempty"`
	// proj_config.mk 中需要设置的配置项
	ConfigFlags map[string]string `yaml:"config_flags,omitempty"`
	// 模板文件路径（相对于 templates/components/）
	TemplateFiles []string `yaml:"template_files,omitempty"`
}

// ComponentsConfig 组件配置文件结构
type ComponentsConfig struct {
	Components []Component `yaml:"components"`
}

// UserConfig 用户配置文件结构
type UserConfig struct {
	SDKPath string `yaml:"sdk_path"`
}

// LoadComponents 从 assets/components.yaml 加载组件配置
func LoadComponents() ([]Component, error) {
	// 尝试从多个位置查找配置文件
	var configPath string
	var found bool

	// 1. 尝试从可执行文件所在目录（安装模式）
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		path := filepath.Join(exeDir, "assets", "components.yaml")
		if _, err := os.Stat(path); err == nil {
			configPath = path
			found = true
		}
	}

	// 2. 尝试从当前工作目录（开发模式）
	if !found {
		cwd, err := os.Getwd()
		if err == nil {
			// 尝试直接在 cwd
			path := filepath.Join(cwd, "assets", "components.yaml")
			if _, err := os.Stat(path); err == nil {
				configPath = path
				found = true
			} else {
				// 尝试在 wb2-cli 子目录
				path = filepath.Join(cwd, "wb2-cli", "assets", "components.yaml")
				if _, err := os.Stat(path); err == nil {
					configPath = path
					found = true
				}
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("找不到组件配置文件 components.yaml，请确保文件存在于 assets/ 目录下")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取组件配置文件失败: %v", err)
	}

	var config ComponentsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析组件配置文件失败: %v", err)
	}

	return config.Components, nil
}

// LoadConfig 加载用户配置文件
func LoadConfig() (*UserConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户主目录失败: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "wb2-cli")
	configPath := filepath.Join(configDir, "config.yaml")

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &UserConfig{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取用户配置文件失败: %v", err)
	}

	var config UserConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析用户配置文件失败: %v", err)
	}

	return &config, nil
}

// SaveConfig 保存用户配置文件
func SaveConfig(config *UserConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "wb2-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}
