package cmd

import (
	"fmt"
	"manga-cli/internals/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set or view configuration options (e.g., download path, viewer, language)",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() 
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration option",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]
		err := config.SetConfigOption(key, value)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Set %s = %s\n", key, value)
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration option",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		val, err := config.GetConfigOption(key)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("%s = %v\n", key, val)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration options",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetAllConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Configurable options:")
		for k, meta := range config.ValidConfigOptions {
			val := cfg[k]
			if val == nil {
				val = meta.Default
			}
			fmt.Printf("  %-15s = %-20v # %s\n", k, val, meta.Description)
		}
	},
}

func init() {
	configCmd.AddCommand(configSetCmd, configGetCmd, configListCmd)
	AddSubCommand(configCmd)
}
