package listUtils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)


func ListTopLevelFolders(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("📁 %s\n", entry.Name())
		}
	}
	return nil
}

func ListFiles(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(root, path)
			fmt.Printf("📄 %s\n", rel)
		}
		return nil
	})
}


func ListLocalDownloads(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}

		depth := len(strings.Split(rel, string(os.PathSeparator))) - 1
		prefix := strings.Repeat("  ", depth)

		if d.IsDir() {
			fmt.Printf("%s📁 %s\n", prefix, d.Name())
		} else {
			fmt.Printf("%s📄 %s\n", prefix, d.Name())
		}
		return nil
	})
}
