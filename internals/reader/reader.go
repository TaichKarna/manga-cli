package readerUtil

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"manga-cli/internals/config"
	"manga-cli/internals/utils"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var viewerCmd string = "viu"
var viewerFlags = []string{"-w", "-h"}

func StartReader(path string, width int, height int) error {
	viewerVal, err := config.GetConfigOption("viewer")
	if err == nil && viewerVal != nil {
		viewerCmd = fmt.Sprintf("%v", viewerVal)
	}

	// Optional: load viewer flags from config if you want
	// e.g. viewerFlagsRaw, _ := config.GetConfigOption("viewerFlags")

	if !isCommandAvailable(viewerCmd) {
		v, err := promptViewer()
		if err != nil || v == "" {
			return errors.New("no valid image viewer configured, aborting")
		}
		viewerCmd = v
		_ = config.SetConfigOption("viewer", viewerCmd)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory '%s': %w", path, err)
	}

	var images []string
	for _, f := range files {
		if !f.IsDir() && isImageFile(f.Name()) {
			images = append(images, filepath.Join(path, f.Name()))
		}
	}

	sort.Strings(images)

	if len(images) == 0 {
		return fmt.Errorf("no images found in folder: %s", path)
	}

	reader := bufio.NewReader(os.Stdin)
	i := 0

	for {
		utils.ClearTerminal()
		fmt.Printf("Page %d / %d \n", i+1, len(images))
		fmt.Println("Commands: [n]ext, [p]rev, [q]uit, [number] jump to page")

		err := renderImageWithTimeout(images[i], width, height, 5*time.Second)
		if err != nil {
			fmt.Printf("Error rendering image: %v\n", err)
		}

		fmt.Print("Enter command: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "n", "":
			if i < len(images)-1 {
				i++
			}
		case "p":
			if i > 0 {
				i--
			}
		case "q":
			return nil
		default:
			pageNum, err := strconv.Atoi(input)
			if err != nil || pageNum < 1 || pageNum > len(images) {
				fmt.Println("Invalid input, enter 'n', 'p', 'q' or a valid page number. Press enter to continue.")
				reader.ReadString('\n')
				continue
			}
			i = pageNum - 1
		}
	}
}

func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
}

func renderImageWithTimeout(path string, width int, height int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := []string{}
	if width > 0 {
		args = append(args, "-w", strconv.Itoa(width))
	}
	if height > 0 {
		args = append(args, "-h", strconv.Itoa(height))
	}
	args = append(args, path)

	cmd := exec.CommandContext(ctx, viewerCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("image viewer command timed out after %s", timeout)
	}

	return err
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
	return strings.TrimSpace(viewer), nil
}
