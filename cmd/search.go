package cmd

import (
	"bufio"
	"fmt"
	"manga-cli/internals/api"
	"manga-cli/internals/downloader"
	readerUtil "manga-cli/internals/reader"
	"manga-cli/internals/utils"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)


var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "List matching manga titles",
	Run: func(cmd *cobra.Command, args []string) {
		if title == "" {
			fmt.Println("Usage: manga-cli search --title 'One Piece'")
			os.Exit(1)
		}
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")

		fmt.Println("Searching for manga:", title)

		resp, err := api.GetMangaIDByTitle(title)

		if err != nil {
			fmt.Println("Error searching manga", err)
			os.Exit(1)
		}
		
		if len(resp.Data) == 0 {
			fmt.Println("No manga found with that title.")
			return
		}

		fmt.Printf("\nFound %d manga(s):\n\n", len(resp.Data))
		for i, manga := range resp.Data {
			fmt.Printf("%d. %s\n", i+1, manga.Attributes.Title["en"])
			fmt.Println()
		}

		selectedManga := selectManga(resp.Data)
		fmt.Println()
		fmt.Print(selectedManga.Attributes.Title["en"])
		selectedChapter := ShowChaptersList(selectedManga.ID)

		err = downloader.DownloadChapter(selectedManga.Attributes.Title["en"], selectedChapter.ID, selectedChapter.Attributes.Chapter, false)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		chapterStr := selectedChapter.Attributes.Chapter

		
		chapterNo, err := strconv.Atoi(chapterStr)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		
				
		path, err := utils.GetPathByTitleAndChapter(selectedManga.Attributes.Title["en"], int(chapterNo))
		
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
	searchCmd.Flags().StringVarP(&title, "title", "t", "", "Manga title (required)")

	searchCmd.MarkFlagRequired("title")

	AddSubCommand(searchCmd)
}


func selectManga(mangaList []api.MangaData) *api.MangaData{
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Select a manga (1-%d) or 'q' to quit: ", len(mangaList))
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		
		if input == "q" || input == "quit" {
			return nil
		}

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter a valid number or 'q' to quit.")
			continue
		}

		if choice < 1 || choice > len(mangaList) {
			fmt.Printf("Please enter a number between 1 and %d.\n", len(mangaList))
			continue
		}

		return &mangaList[choice-1]
	}
}

func ShowChaptersList(selectedManga string) *api.ChapterData {
	reader := bufio.NewReader(os.Stdin)
	limit := 10
	offset := 0
	var allChapters []*api.ChapterData 

	for {
		utils.ClearTerminal()

		listChapters, err := api.GetChapterList(selectedManga, limit, offset)
		if err != nil {
			fmt.Println("Error occurred fetching chapters:", err)
			return nil
		}
		fmt.Printf("Page %d\n", (offset/limit)+1)
		renderChapterList(listChapters.Data)

		fmt.Print("[n] next  [p] prev  [q] quit  [number] select chapter: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "n":
			offset += limit
		case "p":
			if offset >= limit {
				offset -= limit
			}
		case "q":
			return nil
		default:
			num, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input. Press enter to continue.")
				reader.ReadString('\n')
				continue
			}

			if num >= 1 && num <= len(listChapters.Data) {
				return &listChapters.Data[num-1]
			}

			if allChapters == nil {
				allChapters, err = api.FetchAllChapters(selectedManga)
				if err != nil {
					fmt.Println("Error fetching all chapters for search:", err)
					reader.ReadString('\n')
					continue
				}
			}

			for _, ch := range allChapters {
				if ch.Attributes.Chapter == fmt.Sprint(num) {
					return ch
				}
			}

			fmt.Println("Chapter not found. Press enter to continue.")
			reader.ReadString('\n')
		}
	}
}




func renderChapterList(chapterList []api.ChapterData) {
	fmt.Printf("\nChapters (%d):\n\n", len(chapterList))
	for i, chapter := range chapterList {
		title := chapter.Attributes.Title
		if title == "" {
			title = "Untitled Chapter"
		}
		fmt.Printf("%2d. %-40s (Chapter #%s)\n", i+1, title, chapter.Attributes.Chapter)
	}
	fmt.Println()
}


