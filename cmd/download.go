package cmd

import (
	"fmt"
	"manga-cli/internals/api"
	"manga-cli/internals/downloader"
	"os"

	"github.com/spf13/cobra"
)


var from, to int

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a specific manga chapter or range",
	Run: func(cmd *cobra.Command, args []string) {
	title, _ := cmd.Flags().GetString("title")
	chapterNum, _ := cmd.Flags().GetInt("chapter")
	from, _ := cmd.Flags().GetInt("from")
	to, _ := cmd.Flags().GetInt("to")
	dataSaver, _ := cmd.Flags().GetBool("data-saver")

	if title == "" {
		fmt.Println("Please specify --title")
		os.Exit(1)
	}

	if chapterNum > 0 {
		chData, err := api.GetChapterIDByNumber(title, chapterNum)
		if err != nil {
			fmt.Println("Failed to get chapter:", err)
			os.Exit(1)
		}

		err = downloader.DownloadChapter(title, chData.ID, chData.Attributes.Chapter, dataSaver)
		if err != nil {
			fmt.Println("Download error:", err)
			os.Exit(1)
		}
		fmt.Println("Downloaded chapter", chapterNum)
		return
	}

	if from > 0 && to > 0 && to >= from {
		chMap, err := api.GetChapterIDsByRange(title, from, to)
		if err != nil {
			fmt.Println("Failed to fetch chapter range:", err)
			os.Exit(1)
		}

		for i := from; i <= to; i++ {
			chID, ok := chMap[i]
			if !ok {
				fmt.Printf("Chapter %d not found\n", i)
				continue
			}

			err := downloader.DownloadChapter(title, chID, fmt.Sprintf("%d", i), dataSaver)
			if err != nil {
				fmt.Printf("Error downloading chapter %d: %v\n", i, err)
			} else {
				fmt.Printf("✅ Downloaded chapter %d\n", i)
			}
		}
		return
	}

	fmt.Println("Please specify either --chapter or --from and --to")
},

}


func init(){
	downloadCmd.Flags().Int("from", 0, "Start of chapter range")
	downloadCmd.Flags().Int("to", 0, "End of chapter range")
	downloadCmd.Flags().IntVarP(&chapter, "chapter", "c", 0, "Specific chapter number")
	downloadCmd.Flags().StringVarP(&title, "title", "t", "", "Manga title (required)")
	downloadCmd.Flags().Bool("data-saver", false, "Use data-saver mode for lower quality images")

	downloadCmd.MarkFlagRequired("title")

	AddSubCommand(downloadCmd)
}

