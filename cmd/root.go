package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	sdkPath string
	version = "dev" // This will be set during build
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wb2-cli",
	Short: "Ai-Thinker WB2 项目生成工具",
	Long: `wb2-cli 是一个用于快速创建 Ai-Thinker WB2 芯片项目的命令行工具。
它可以帮助您快速生成项目框架，并选择需要的组件。`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 全局 flags
	rootCmd.PersistentFlags().StringVar(&sdkPath, "sdk-path", "", "SDK 根目录路径（如果未设置，将从配置文件读取）")
	rootCmd.Flags().BoolP("version", "v", false, "显示版本信息")

	// Override the default version flag behavior
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			cmd.Printf("wb2-cli version %s\n", version)
			return
		}
		cmd.Help()
	}
}
