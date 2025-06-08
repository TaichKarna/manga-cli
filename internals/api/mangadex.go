package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type AtHomeResponse struct {
	BaseURL string  `json:"baseUrl"`
	Chapter Chapter `json:"chapter"`
}

type Chapter struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

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


func GetChapterIDsByRange(title string, from, to int) (map[int]string, error) {
	mangaResult, err := GetMangaIDByTitle(title)
	if err != nil {
		return nil, fmt.Errorf("failed to get manga ID: %w", err)
	}

	query := url.Values{}
	query.Set("manga", mangaResult.Data[0].ID)
	query.Set("limit", "100")
	query.Add("translatedLanguage[]", "en")

	for i := from; i <= to; i++ {
		query.Add("chapter[]", strconv.Itoa(i))
	}

	fullURL := fmt.Sprintf("https://api.mangadex.org/chapter?%s", query.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var result ChapterSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	chapterMap := make(map[int]string)
	for _, ch := range result.Data {
		num, err := strconv.Atoi(ch.Attributes.Chapter)
		if err != nil {
			continue 
		}
		chapterMap[num] = ch.ID
	}

	return chapterMap, nil
}

func GetChapterIDByNumber(title string, chapterNumber int) (ChapterData, error) {
	mangaResult, err := GetMangaIDByTitle(title)
	if err != nil {
		return ChapterData{}, fmt.Errorf("failed to get manga ID: %w", err)
	}

	endpoint := "https://api.mangadex.org/chapter"
	query := url.Values{}
	query.Set("manga", mangaResult.Data[0].ID)
	query.Set("chapter", strconv.Itoa(chapterNumber))
	query.Add("translatedLanguage[]", "en")
	query.Set("limit", "1")

	fullURL := fmt.Sprintf("%s?%s", endpoint, query.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return ChapterData{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ChapterData{}, fmt.Errorf("unexpected response: %s", resp.Status)
	}

	var data ChapterSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ChapterData{}, fmt.Errorf("JSON decode error: %w", err)
	}

	if len(data.Data) == 0 {
		return ChapterData{}, fmt.Errorf("chapter %d not found for '%s'", chapterNumber, title)
	}

	return data.Data[0], nil
}