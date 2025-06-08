package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"manga-cli/internals/utils"
	"net/http"
	"os"
	"path/filepath"
)

type AtHomeResponse struct {
	BaseURL string  `json:"baseUrl"`
	Chapter Chapter `json:"chapter"`
}

type Chapter struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

func DownloadChapter(title string, chapterID string, chapterNo string, useDataSaver bool) error {
	fmt.Printf(" Downloading chapter %s of \"%s\"...\n", chapterNo, title)

	savePath, err := searchOrCreateFolder(title, chapterNo)
	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}
	fmt.Printf(" Saving to: %s\n", savePath)

	atHomeResp, err := getAtHomeServer(chapterID)
	if err != nil {
		return fmt.Errorf("failed to get at-home server info: %w", err)
	}

	fmt.Printf("Got server: %s\n", atHomeResp.BaseURL)
	fmt.Printf(" Found %d pages to download.\n", len(atHomeResp.Chapter.Data))
	if useDataSaver {
		fmt.Println("Using Data Saver mode")
	}

	err = downloadChapterPages(atHomeResp, savePath, useDataSaver)
	if err != nil {
		return fmt.Errorf("failed to download pages: %w", err)
	}

	fmt.Println(" Chapter download complete.")
	return nil
}


func searchOrCreateFolder(title string, chapterNo string) (string, error) {
	mangaCliDir, err := utils.GetOrCreateMangaCliDir()
	if err != nil {
		return "", err
	}

	savePath := filepath.Join(mangaCliDir, title, chapterNo)
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return savePath, nil
}

func getAtHomeServer(chapterID string) (*AtHomeResponse, error) {
	url := fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", chapterID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var atHomeResp AtHomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&atHomeResp); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return &atHomeResp, nil
}
func downloadChapterPages(atHomeResp *AtHomeResponse, folderPath string, useDataSaver bool) error {
	pages := atHomeResp.Chapter.Data
	if useDataSaver {
		pages = atHomeResp.Chapter.DataSaver
	}

	var failedPages []string
	totalPages := len(pages)

	for i, page := range pages {
		filePath := filepath.Join(folderPath, page)

		progress := fmt.Sprintf("[%d/%d]", i+1, totalPages)

		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("%s Skipped (exists): %s\n", progress, page)
			continue
		}

		url := fmt.Sprintf("%s/data/%s/%s", atHomeResp.BaseURL, atHomeResp.Chapter.Hash, page)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s Failed to GET %s: %v\n", progress, page, err)
			failedPages = append(failedPages, page)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("%s Bad status for %s: %s\n", progress, page, resp.Status)
			failedPages = append(failedPages, page)
			continue
		}

		outFile, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("%s Failed to create file %s: %v\n", progress, page, err)
			failedPages = append(failedPages, page)
			continue
		}

		_, err = io.Copy(outFile, resp.Body)
		outFile.Close()

		if err != nil {
			fmt.Printf("%s  Failed to write file %s: %v\n", progress, page, err)
			failedPages = append(failedPages, page)
			continue
		}

		fmt.Printf("%s Downloaded: %s\n", progress, page)
	}

	if len(failedPages) > 0 {
		return fmt.Errorf("failed to download %d pages: %v", len(failedPages), failedPages)
	}

	return nil
}
