package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func GetOrCreateMangaCliDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	picturesDir := filepath.Join(homeDir, "Pictures")
	mangaCliDir := filepath.Join(picturesDir, "manga-cli")

	if _, err := os.Stat(mangaCliDir); os.IsNotExist(err) {
		err = os.MkdirAll(mangaCliDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create manga-cli directory: %w", err)
		}
	}

	return mangaCliDir, nil
}


func GetPathByTitleAndChapter(title string, chapter int) (string, error) {
	mangaCliDir, err := GetOrCreateMangaCliDir()
	if err != nil {
		return "", err
	}

	chapterStr := strconv.Itoa(chapter)
	path := filepath.Join(mangaCliDir, title, chapterStr)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("manga doesn't exist; download it using the 'download' command: %w", err)
	}

	return path, nil
}


func ClearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}