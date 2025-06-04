package cmd

import (
	"fmt"
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

		path, err := utils.GetPathByTitleAndChapter(title, chapter)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := readerUtil.StartReader(path); err != nil {
			fmt.Println("Failed to start reader:", err)
			os.Exit(1)
		}
	},

}

func init(){
	readCmd.Flags().StringVarP(&title, "title", "t", "", "Manga title (required)")
	readCmd.Flags().IntVarP(&chapter, "chapter", "c", 0, "Chapter number (required)")

	readCmd.MarkFlagRequired("title")
	readCmd.MarkFlagRequired("chapter")

	AddSubCommand(readCmd)
}