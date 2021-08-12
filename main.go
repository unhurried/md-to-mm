package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"md-to-mm/freeplane"

	"github.com/russross/blackfriday/v2"
)

func main() {
	// Check whether an input filename is specified.
	if len(os.Args) < 2 {
		fmt.Println("Please specify a markdown file.")
		os.Exit(1)
	}
	file := os.Args[1]

	markdown, error := ioutil.ReadFile(file)
	if error != nil {
		fmt.Println("Please specify a markdown file.")
		os.Exit(1)
	}

	// Construct a path to save the output file.
	targetDir := filepath.Dir(file)
	targetPath := filepath.Join(targetDir, "result.mm")

	// Render a mind map file and store it to the destination path.
	markdown = []byte(strings.ReplaceAll(string(markdown), "\\", "\\\\"))
	renderer := &freeplane.Renderer{}
	output := blackfriday.Run(markdown, blackfriday.WithRenderer(renderer))
	ioutil.WriteFile(targetPath, output, 0644)
}
