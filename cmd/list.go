package cmd

import (
	"fmt"
	"manga-cli/internals/config"
	"manga-cli/internals/listUtils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List downloaded manga, chapters, or images",
	Run: func(cmd *cobra.Command, args []string) {
		pathFlag, _ := cmd.Flags().GetString("path")
		title, _ := cmd.Flags().GetString("title")
		chapter, _ := cmd.Flags().GetString("chapter")

		var basePath string
		if pathFlag != "" {
			basePath = pathFlag
		} else {
			val, err := config.GetConfigOption("path")
			if err != nil {
				basePath = "./downloads"
			} else {
				basePath = fmt.Sprintf("%v", val)
			}
		}

		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			fmt.Printf("Download path '%s' does not exist or is empty.\n", basePath)
			return
		}

		if title == "" {
			fmt.Println("Available manga:")
			if err := listUtils.ListTopLevelFolders(basePath); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}

		mangaPath := filepath.Join(basePath, title)
		if _, err := os.Stat(mangaPath); os.IsNotExist(err) {
			fmt.Printf("Manga title '%s' not found in downloads.\n", title)
			return
		}

		if chapter == "" {
			fmt.Printf("Chapters for manga '%s':\n", title)
			if err := listUtils.ListTopLevelFolders(mangaPath); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}

		chapterPath := filepath.Join(mangaPath, chapter)
		if _, err := os.Stat(chapterPath); os.IsNotExist(err) {
			fmt.Printf("Chapter '%s' not found under manga '%s'.\n", chapter, title)
			return
		}

		fmt.Printf("Images in chapter '%s' of manga '%s':\n", chapter, title)
		if err := listUtils.ListFiles(chapterPath); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	listCmd.Flags().String("path", "", "Override the download path")
	listCmd.Flags().String("title", "", "Manga title to list chapters")
	listCmd.Flags().String("chapter", "", "Chapter number/title to list images")
	AddSubCommand(listCmd)
}
