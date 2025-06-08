package readerUtil

import (
	"bufio"
	"fmt"
	"manga-cli/internals/config"
	"manga-cli/internals/utils"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var viewerCmd string = "viu"

func StartReader(path string, width int, height int) error {

	viewerVal, err := config.GetConfigOption("viewer")
	if err == nil && viewerVal != nil {
		viewerCmd = fmt.Sprintf("%v", viewerVal)
	}

	if !isCommandAvailable(viewerCmd) {
		v, err := promptViewer()
		if err != nil || v == "" {
			return fmt.Errorf("no valid image viewer configured, aborting")
		}
		viewerCmd = v
		_ = config.SetConfigOption("viewer", viewerCmd)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var images []string
	for _, f := range files {
		if !f.IsDir() && isImageFile(f.Name()) {
			images = append(images, filepath.Join(path, f.Name()))
		}
	}

	sort.Strings(images)

	if len(images) == 0 {
		return fmt.Errorf("no images found in folder: %s ", path)
	}

	reader := bufio.NewReader(os.Stdin)

	i := 0
	for {
		utils.ClearTerminal()

		fmt.Printf("Page %d / %d \n", i+1, len(images))
		if err := renderImage(images[i], width, height); err != nil {
			fmt.Println("Error rendering panel:", err)
		}

		fmt.Print("[n] next  [p] prev  [q] quit: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "n":
			if i < len(images)-1 {
				i++
			}
		case "p":
			if i > 0 {
				i--
			}
		case "q":
			return nil
		}
	}
}

func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
}

func renderImage(path string, width int, height int) error {
	cmd := exec.Command(viewerCmd, "-w", strconv.Itoa(width), "-h", strconv.Itoa(height), path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func promptViewer() (string, error) {
	fmt.Print("No image viewer found. Please enter your preferred image viewer command (e.g., viu, chafa): ")
	reader := bufio.NewReader(os.Stdin)
	viewer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	viewer = strings.TrimSpace(viewer)
	return viewer, nil
}
