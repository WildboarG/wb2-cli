package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComponentStruct(t *testing.T) {
	// Test Component struct creation
	component := Component{
		Name:        "test_component",
		Description: "Test component description",
		Category:    "network",
		Dependencies: []string{"wifi"},
		SDKComponents: []string{"component1"},
		ConfigFlags: map[string]string{
			"CONFIG_TEST": "1",
		},
		TemplateFiles: []string{"test.c.tmpl"},
	}

	if component.Name != "test_component" {
		t.Errorf("Expected name 'test_component', got '%s'", component.Name)
	}

	if component.Description != "Test component description" {
		t.Errorf("Expected description 'Test component description', got '%s'", component.Description)
	}

	if component.Category != "network" {
		t.Errorf("Expected category 'network', got '%s'", component.Category)
	}

	if len(component.Dependencies) != 1 || component.Dependencies[0] != "wifi" {
		t.Errorf("Expected dependencies ['wifi'], got %v", component.Dependencies)
	}

	if component.ConfigFlags["CONFIG_TEST"] != "1" {
		t.Errorf("Expected CONFIG_TEST=1, got %s", component.ConfigFlags["CONFIG_TEST"])
	}
}

func TestComponentsConfigStruct(t *testing.T) {
	// Test ComponentsConfig struct
	config := ComponentsConfig{
		Components: []Component{
			{Name: "comp1", Description: "Component 1"},
			{Name: "comp2", Description: "Component 2"},
		},
	}

	if len(config.Components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(config.Components))
	}

	if config.Components[0].Name != "comp1" {
		t.Errorf("Expected first component name 'comp1', got '%s'", config.Components[0].Name)
	}
}

func TestLoadComponents(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create assets directory
	assetsDir := filepath.Join(tempDir, "assets")
	err := os.MkdirAll(assetsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create assets directory: %v", err)
	}

	// Create a test components.yaml file
	testConfig := `components:
  - name: test_wifi
    description: Test WiFi component
    category: network
  - name: test_gpio
    description: Test GPIO component
    category: peripheral
    dependencies:
      - test_wifi`

	configPath := filepath.Join(assetsDir, "components.yaml")
	err = os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Change to temp directory to simulate the loading logic
	oldWd, _ := os.Getwd()
	defer func() {
		os.Chdir(oldWd)
	}()

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test loading components
	components, err := LoadComponents()
	if err != nil {
		t.Fatalf("LoadComponents failed: %v", err)
	}

	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}

	// Check first component
	if components[0].Name != "test_wifi" {
		t.Errorf("Expected first component name 'test_wifi', got '%s'", components[0].Name)
	}
	if components[0].Category != "network" {
		t.Errorf("Expected first component category 'network', got '%s'", components[0].Category)
	}

	// Check second component
	if components[1].Name != "test_gpio" {
		t.Errorf("Expected second component name 'test_gpio', got '%s'", components[1].Name)
	}
	if components[1].Category != "peripheral" {
		t.Errorf("Expected second component category 'peripheral', got '%s'", components[1].Category)
	}
	if len(components[1].Dependencies) != 1 || components[1].Dependencies[0] != "test_wifi" {
		t.Errorf("Expected dependencies ['test_wifi'], got %v", components[1].Dependencies)
	}
}

func TestLoadComponentsFileNotFound(t *testing.T) {
	// Change to a directory without components.yaml
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)

	err := os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test loading components when file doesn't exist
	_, err = LoadComponents()
	if err == nil {
		t.Error("Expected error when components.yaml not found, but got nil")
	}

	if !strings.Contains(err.Error(), "找不到组件配置文件") {
		t.Errorf("Expected error message about missing config file, got: %v", err)
	}
}

func TestUserConfigStruct(t *testing.T) {
	// Test UserConfig struct
	config := UserConfig{
		SDKPath: "/path/to/sdk",
	}

	if config.SDKPath != "/path/to/sdk" {
		t.Errorf("Expected SDK path '/path/to/sdk', got '%s'", config.SDKPath)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary home directory for testing
	tempDir := t.TempDir()

	// Mock the user home directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tempDir)

	// Create config directory and file
	configDir := filepath.Join(tempDir, ".config", "wb2-cli")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	testConfig := `sdk_path: /custom/sdk/path`
	configPath := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.SDKPath != "/custom/sdk/path" {
		t.Errorf("Expected SDK path '/custom/sdk/path', got '%s'", config.SDKPath)
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	// Create a temporary home directory without config file
	tempDir := t.TempDir()

	// Mock the user home directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tempDir)

	// Test loading config when file doesn't exist
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should not fail when config doesn't exist, got: %v", err)
	}

	// Should return default config (empty SDK path)
	if config.SDKPath != "" {
		t.Errorf("Expected empty SDK path for missing config, got '%s'", config.SDKPath)
	}
}