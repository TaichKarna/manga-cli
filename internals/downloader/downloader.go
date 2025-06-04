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

// DownloadChapter downloads all pages of a given manga chapter
// title: manga title (used for folder structure)
// chapterID: MangaDex chapter UUID string (not the integer chapter number)
// useDataSaver: if true, download lower quality images (dataSaver)
func DownloadChapter(title string, chapterID string, chapterNo string, useDataSaver bool) error {
	savePath, err := searchOrCreateFolder(title,  chapterNo)
	if err != nil {
		return err
	}

	fmt.Println(savePath)

	atHomeResp, err := getAtHomeServer(chapterID)
	if err != nil {
		return fmt.Errorf("failed to get at-home server info: %w", err)
	}

	fmt.Println(atHomeResp)

	err = downloadChapterPages(atHomeResp, savePath, useDataSaver)
	if err != nil {
		return fmt.Errorf("failed to download pages: %w", err)
	}

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

	for _, page := range pages {
		url := fmt.Sprintf("%s/data/%s/%s", atHomeResp.BaseURL, atHomeResp.Chapter.Hash, page)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error: failed to GET %s: %v\n", url, err)
			failedPages = append(failedPages, page)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: bad status for %s: %s\n", page, resp.Status)
			resp.Body.Close()
			failedPages = append(failedPages, page)
			continue
		}

		filePath := filepath.Join(folderPath, page)
		outFile, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error: failed to create file %s: %v\n", filePath, err)
			resp.Body.Close()
			failedPages = append(failedPages, page)
			continue
		}

		_, err = io.Copy(outFile, resp.Body)
		outFile.Close()
		resp.Body.Close()

		if err != nil {
			fmt.Printf("Error: failed to write file %s: %v\n", filePath, err)
			failedPages = append(failedPages, page)
			continue
		}

	}

	if len(failedPages) > 0 {
		return fmt.Errorf("failed to download %d pages: %v", len(failedPages), failedPages)
	}

	return nil
}
