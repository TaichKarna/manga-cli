package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	title string 
	chapter int
)

var rootCmd = &cobra.Command{
	Use:   "manga-cli",
	Short: "A terminal manga reader and downloader",
	Long:  `manga-cli is a terminal tool for downloading and reading manga from supported sources.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use a subcommand: download, read, list, etc.")
		_ = cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func AddSubCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
