package cmd

import (
	"fmt"
	"manga-cli/internals/api"
	"manga-cli/internals/config"
	"manga-cli/internals/downloader"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)


var from, to int

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a specific manga chapter or range",
	Run: func(cmd *cobra.Command, args []string) {
		if title == "" {
			fmt.Println("Error: --title is required")
			os.Exit(1)
		}

		if chapter == 0 && (from == 0 || to == 0) {
			fmt.Println("Error: Either --chapter or both --from and --to must be provided")
			os.Exit(1)
		}

		if chapter != 0 {
			// Single chapter mode
			if err := downloadChapter(title, chapter); err != nil {
				fmt.Println("❌", err)
				os.Exit(1)
			}
			return
		}

		// Range mode
		if from > to {
			fmt.Println("Error: --from must be less than or equal to --to")
			os.Exit(1)
		}

		for i := from; i <= to; i++ {
			fmt.Printf("\n📘 Downloading chapter %d...\n", i)
			if err := downloadChapter(title, i); err != nil {
				fmt.Printf("❌ Failed to download chapter %d: %v\n", i, err)
			}
		}
	},
}


func init(){
	downloadCmd.Flags().Int("from", 0, "Start of chapter range")
	downloadCmd.Flags().Int("to", 0, "End of chapter range")
	downloadCmd.Flags().IntVarP(&chapter, "chapter", "c", 0, "Specific chapter number")
	downloadCmd.Flags().StringVarP(&title, "title", "t", "", "Manga title (required)")
	
	downloadCmd.MarkFlagRequired("title")

	AddSubCommand(downloadCmd)
}

func downloadChapter(title string, chapter int) error {
	mangaID, err := api.GetMangaIDByTitle(title)
	if err != nil {
		return fmt.Errorf("failed to get manga ID: %w", err)
	}

	chapterID, err := api.GetChapterIDByNumber(mangaID, chapter)
	if err != nil {
		return fmt.Errorf("failed to get chapter ID: %w", err)
	}

	pages, err := api.GetChapterPages(chapterID)
	if err != nil {
		return fmt.Errorf("failed to get pages: %w", err)
	}

	basePath := "./downloads"
	if val, err := config.GetConfigOption("path"); err == nil {
		basePath = fmt.Sprintf("%v", val)
	}
	targetPath := filepath.Join(basePath, title, fmt.Sprintf("%v", chapter))
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := downloader.DownloadImages(pages, targetPath); err != nil {
		return fmt.Errorf("failed to download images: %w", err)
	}

	fmt.Printf("✅ Chapter %d of '%s' downloaded to: %s\n", chapter, title, targetPath)
	return nil
}
