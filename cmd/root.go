package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	title string 
	chapter int
	width int 
	height int
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

func init() {
    rootCmd.PersistentFlags().IntVar(&width, "width", 0, "Width of image viewer")
    rootCmd.PersistentFlags().IntVar(&height, "height", 0, "Height of image viewer")
	rootCmd.Flags().IntVar(&from, "from", 0, "Start of chapter range")
	rootCmd.Flags().IntVar(&to, "to", 0, "End of chapter range")	

}

func Execute() error {
	return rootCmd.Execute()
}

func AddSubCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
