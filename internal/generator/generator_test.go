package generator

import (
	"testing"

	"wb2-cli/internal/config"
)

func TestNewGenerator(t *testing.T) {
	sdkPath := "/path/to/sdk"
	gen := New(sdkPath)

	if gen == nil {
		t.Fatal("New() returned nil")
	}

	if gen.sdkPath != sdkPath {
		t.Errorf("Expected SDK path '%s', got '%s'", sdkPath, gen.sdkPath)
	}
}

func TestProjectDataStruct(t *testing.T) {
	// Test ProjectData struct creation
	data := &ProjectData{
		ProjectName: "test_project",
		SDKPath:     "/sdk/path",
		Components: []config.Component{
			{Name: "wifi", Description: "WiFi component"},
		},
		HasWifi: true,
	}

	if data.ProjectName != "test_project" {
		t.Errorf("Expected project name 'test_project', got '%s'", data.ProjectName)
	}

	if data.SDKPath != "/sdk/path" {
		t.Errorf("Expected SDK path '/sdk/path', got '%s'", data.SDKPath)
	}

	if len(data.Components) != 1 {
		t.Errorf("Expected 1 component, got %d", len(data.Components))
	}

	if !data.HasWifi {
		t.Error("Expected HasWifi to be true")
	}
}

func TestPrepareProjectData(t *testing.T) {
	gen := New("/test/sdk/path")

	components := []config.Component{
		{
			Name:             "wifi",
			Description:      "WiFi component",
			Category:         "network",
			IncludeComponents: []string{"wifi_station", "wifi_softap"},
			ConfigFlags: map[string]string{
				"CONFIG_WIFI_ENABLE": "1",
			},
		},
		{
			Name:        "gpio",
			Description: "GPIO component",
			Category:    "peripheral",
		},
	}

	data := gen.prepareProjectData("test_project", components)

	// Test basic fields
	if data.ProjectName != "test_project" {
		t.Errorf("Expected project name 'test_project', got '%s'", data.ProjectName)
	}

	if data.SDKPath != "/test/sdk/path" {
		t.Errorf("Expected SDK path '/test/sdk/path', got '%s'", data.SDKPath)
	}

	// Test component detection
	if !data.HasWifi {
		t.Error("Expected HasWifi to be true when wifi component is present")
	}

	if !data.HasGPIO {
		t.Error("Expected HasGPIO to be true when gpio component is present")
	}

	// Test SDK components inclusion
	if len(data.IncludeComps) == 0 {
		t.Error("Expected IncludeComps to contain components")
	}

	found := false
	for _, comp := range data.IncludeComps {
		if comp == "wifi_station" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'wifi_station' in IncludeComps")
	}

	// Test config flags
	if data.ConfigFlags["CONFIG_WIFI_ENABLE"] != "1" {
		t.Errorf("Expected CONFIG_WIFI_ENABLE=1, got %s", data.ConfigFlags["CONFIG_WIFI_ENABLE"])
	}
}

func TestPrepareProjectDataEmptyComponents(t *testing.T) {
	gen := New("/sdk/path")
	data := gen.prepareProjectData("empty_project", []config.Component{})

	if data.ProjectName != "empty_project" {
		t.Errorf("Expected project name 'empty_project', got '%s'", data.ProjectName)
	}

	// Test that all HasXXX flags are false for empty components
	if data.HasWifi || data.HasGPIO || data.HasBLE {
		t.Error("Expected all HasXXX flags to be false for empty components")
	}
}


func TestGeneratorCreation(t *testing.T) {
	// Test basic generator creation and methods
	gen := New("/test/sdk/path")

	// Test that generator has expected SDK path
	if gen.sdkPath != "/test/sdk/path" {
		t.Errorf("Expected SDK path '/test/sdk/path', got '%s'", gen.sdkPath)
	}

	// Test that we can call public methods without panic
	data := gen.prepareProjectData("test", []config.Component{})
	if data.ProjectName != "test" {
		t.Errorf("Expected project name 'test', got '%s'", data.ProjectName)
	}
}