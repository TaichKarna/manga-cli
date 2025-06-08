package utils

import (
	"fmt"
	"manga-cli/internals/config"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func GetOrCreateMangaCliDir() (string, error) {
	configPath, err := config.GetConfigOption("path")
	var basePath string
	print(configPath)

	if err != nil || configPath == nil || configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting home directory: %w", err)
		}
		basePath = filepath.Join(homeDir, "Pictures", "manga-cli")
	} else {
		basePath = fmt.Sprintf("%v", configPath)
	}

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return "", fmt.Errorf("failed to create manga-cli directory: %w", err)
		}
	}

	return basePath, nil
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