package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	// Test that root command is properly configured
	if rootCmd.Use != "wb2-cli" {
		t.Errorf("Expected command use 'wb2-cli', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Root command should have a short description")
	}

	if rootCmd.Long == "" {
		t.Error("Root command should have a long description")
	}
}

func TestRootCmdVersionFlag(t *testing.T) {
	// Test version flag
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Set version flag
	rootCmd.Flags().Set("version", "true")

	// Reset after test
	defer rootCmd.Flags().Set("version", "false")

	// The version flag should be handled in the Run function
	// This is a basic test to ensure the flag exists
	if !rootCmd.Flags().HasAvailableFlags() {
		t.Error("Root command should have flags")
	}

	versionFlag := rootCmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Error("Version flag should exist")
	}

	if versionFlag.Shorthand != "v" {
		t.Errorf("Expected version flag shorthand 'v', got '%s'", versionFlag.Shorthand)
	}
}

func TestExecute(t *testing.T) {
	// Test that Execute function exists and doesn't panic immediately
	// Note: This is a smoke test since Execute calls os.Exit on error
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute panicked: %v", r)
		}
	}()

	// We can't easily test Execute directly without causing the program to exit
	// In a real scenario, you might want to test command execution through integration tests
}

func TestSDKPathFlag(t *testing.T) {
	// Test that sdk-path flag exists
	// Note: Flags are initialized in init() function
	// We need to ensure the init function has run

	sdkPathFlag := rootCmd.PersistentFlags().Lookup("sdk-path")
	if sdkPathFlag == nil {
		t.Error("sdk-path persistent flag should exist")
	}

	if sdkPathFlag != nil && sdkPathFlag.Usage == "" {
		t.Error("sdk-path flag should have usage description")
	}
}

func TestInitFunction(t *testing.T) {
	// Test that init function has set up flags properly
	// This is tested indirectly through the flag existence tests above

	// Verify that the persistent flags are set
	if rootCmd.PersistentFlags().Lookup("sdk-path") == nil {
		t.Error("sdk-path should be a persistent flag")
	}
}