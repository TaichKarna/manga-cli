package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a manga and save to local storage",
	Run: func(cmd *cobra.Command, args []string) {
		if title == "" {
			fmt.Println("Usage: manga-cli donwload --title 'One Piece' \n manga-cli download --title 'One Piece' --chapter 1")
			os.Exit(1)
		}
		fmt.Println("Downloading manga:", title)

	},

}

func init(){
	downloadCmd.Flags().StringVarP(&title, "title", "t", "", "Manga title (required)")
	downloadCmd.Flags().IntVarP(&chapter, "chapter", "c", 0, "Chapter number")

	downloadCmd.MarkFlagRequired("title")

	AddSubCommand(downloadCmd)
}