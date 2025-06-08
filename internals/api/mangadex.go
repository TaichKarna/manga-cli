package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

const baseURL = "https://api.mangadex.org"

type MangaSearchResult struct {
	Data []MangaData `json:"data"`
}

type MangaData struct {
		ID         string `json:"id"`
		Attributes struct {
			Title map[string]string `json:"title"`
		} `json:"attributes"`
} 


type ChapterSearchResult struct {
	Data []ChapterData `json:"data"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type ChapterData struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Volume             string `json:"volume"`
		Chapter            string `json:"chapter"`
		Title              string `json:"title"`
		TranslatedLanguage string `json:"translatedLanguage"`
	} `json:"attributes"`
	
}

func GetMangaIDByTitle(title string) (MangaSearchResult, error) {
	endpoint := fmt.Sprintf("%s/manga", baseURL)

	params := url.Values{}
	params.Set("title", title)

	resp, err := http.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return MangaSearchResult{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MangaSearchResult{}, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MangaSearchResult{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var result MangaSearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return MangaSearchResult{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(result.Data) == 0 {
		return MangaSearchResult{}, errors.New("no manga found with the given title")
	}

	return result, nil
}


func GetChapterList(mangaId string, limit int, offset int) (*ChapterSearchResult, error) {
	endpoint := fmt.Sprintf("%s/manga/%s/feed", baseURL, mangaId)

	params := url.Values{}
	params.Add("translatedLanguage[]", "en")
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))
	params.Add("order[chapter]", "asc")
	params.Add("order[volume]", "asc")

	finalURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	resp, err := http.Get(finalURL)
	if err != nil {
		return &ChapterSearchResult{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &ChapterSearchResult{}, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ChapterSearchResult{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var result ChapterSearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &ChapterSearchResult{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &result, nil
}

func FetchAllChapters(mangaID string) ([]*ChapterData, error) {
	var all []*ChapterData
	limit := 100
	offset := 0

	for {
		list, err := GetChapterList(mangaID, limit, offset)
		if err != nil {
			return nil, err
		}

		for i := range list.Data {
			all = append(all, &list.Data[i])
		}


		if len(list.Data) < limit {
			break
		}
		offset += limit
	}
	return all, nil
}

// Structs to unmarshal at-home server response JSON
type AtHomeResponse struct {
	BaseURL string  `json:"baseUrl"`
	Chapter Chapter `json:"chapter"`
}

type Chapter struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

// Fetch the at-home server info for the chapter
func GetAtHomeServer(chapterID string) (*AtHomeResponse, error) {
	url := fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", chapterID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get at-home server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	var atHomeResp AtHomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&atHomeResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &atHomeResp, nil
}

func DownloadChapterPages(atHomeResp *AtHomeResponse, folderPath string, useDataSaver bool) error {
	pages := atHomeResp.Chapter.Data
	if useDataSaver {
		pages = atHomeResp.Chapter.DataSaver
	}

	for _, page := range pages {
		url := fmt.Sprintf("%s/data/%s/%s", atHomeResp.BaseURL, atHomeResp.Chapter.Hash, page)

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download page %s: %w", page, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad response for page %s: %s", page, resp.Status)
		}

		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create folder %s: %w", folderPath, err)
		}

		filePath := filepath.Join(folderPath, page)
		outFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filePath, err)
		}

		_, err = io.Copy(outFile, resp.Body)
		outFile.Close()
		if err != nil {
			return fmt.Errorf("failed to write page %s: %w", page, err)
		}

		fmt.Printf("Downloaded %s\n", page)
	}

	fmt.Printf("Downloaded %d pages.\n", len(pages))
	return nil
}

