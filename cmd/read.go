package cmd

import (
	"fmt"
	"manga-cli/internals/config"
	readerUtil "manga-cli/internals/reader"
	"manga-cli/internals/utils"
	"os"

	"github.com/spf13/cobra"
)



var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a downloaded manga from local storage",
	Run: func(cmd *cobra.Command, args []string) {
		if title == "" || chapter == 0 {
			fmt.Println("Usage: manga-cli read --title 'One Piece' --chapter 1012")
			os.Exit(1)
		}

		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")

		if width == 0 || height == 0 {
			cfg, _ := config.GetAllConfig()

			if width == 0 {
				if val, ok := cfg["width"]; ok {
					if w, ok := val.(float64); ok {
						width = int(w)
					}
				}
			}

			if height == 0 {
				if val, ok := cfg["height"]; ok {
					if h, ok := val.(float64); ok {
						height = int(h)
					}
				}
			}
		}

		path, err := utils.GetPathByTitleAndChapter(title, chapter)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := readerUtil.StartReader(path, width, height); err != nil {
			fmt.Println("Failed to start reader:", err)
			os.Exit(1)
		}
	},

}

func init(){
	readCmd.Flags().StringVarP(&title, "title", "t", "", "Manga title (required)")
	readCmd.Flags().IntVarP(&chapter, "chapter", "c", 0, "Chapter number (required)")
	readCmd.Flags().Int("width", 0, "Width of image viewer")
	readCmd.Flags().Int("height", 0, "Height of image viewer")
	
	readCmd.MarkFlagRequired("title")
	readCmd.MarkFlagRequired("chapter")

	AddSubCommand(readCmd)
}