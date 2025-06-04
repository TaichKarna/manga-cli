package readerUtil

import (
	"bufio"
	"fmt"
	"manga-cli/internals/utils"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func StartReader(path string) error{

	files, err := os.ReadDir(path)

	if err != nil{
		return err
	}

	var images []string

	for _, f := range files {
		if !f.IsDir() && isImageFile(f.Name()){
			images  = append(images, filepath.Join(path, f.Name()))
		}
	}

	sort.Strings(images)

	if len(images) == 0{
		return  fmt.Errorf("no images found in folder: %s ", path)
	}

	reader := bufio.NewReader(os.Stdin)

	i := 0

	for{
		utils.ClearTerminal()

		fmt.Printf("Page %d / %d \n ", i + 1, len(images))
		if err := renderImage(images[i]); err != nil{
			fmt.Println("Error rendering panel")
		}

		fmt.Print("[n] next  [p] prev  [q] quit: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input{
		case "n":
			if i < len(images) - 1 {
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


func renderImage(path string) error {
	cmd := exec.Command("viu", "-w", "80", "-h", "40", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

