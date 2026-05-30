package freeplane

import (
	"bytes"
	"os"
	"strings"
	"testing"

	bf "github.com/russross/blackfriday/v2"
)

// TestRenderExampleMarkdown is a golden test that verifies the renderer
// produces the expected output for the example markdown file.
func TestRenderExampleMarkdown(t *testing.T) {
	// Reset global state before test
	resetRendererState()

	markdown, err := os.ReadFile("../docs/example.md")
	if err != nil {
		t.Fatalf("Failed to read example markdown: %v", err)
	}

	expected, err := os.ReadFile("../docs/result.mm")
	if err != nil {
		t.Fatalf("Failed to read expected result: %v", err)
	}

	markdown = []byte(strings.ReplaceAll(string(markdown), "\\", "\\\\"))
	renderer := &Renderer{}
	output := bf.Run(markdown, bf.WithRenderer(renderer))

	if !bytes.Equal(output, expected) {
		t.Errorf("Output does not match expected result.\nGot:\n%s\nExpected:\n%s", string(output), string(expected))
	}
}

// TestHeadingRendering tests that headings are rendered correctly.
func TestHeadingRendering(t *testing.T) {
	resetRendererState()

	input := "# Heading 1\n\n## Heading 2\n"
	renderer := &Renderer{}
	output := bf.Run([]byte(input), bf.WithRenderer(renderer))
	result := string(output)

	if !strings.Contains(result, "<map version=\"freeplane 1.8.0\">") {
		t.Error("Output should contain map element")
	}
	if !strings.Contains(result, "<h1>Heading 1</h1>") {
		t.Error("Output should contain h1 heading")
	}
	if !strings.Contains(result, "<h2>Heading 2</h2>") {
		t.Error("Output should contain h2 heading")
	}
}

// TestListRendering tests that list items are rendered correctly.
func TestListRendering(t *testing.T) {
	resetRendererState()

	input := "* Item 1\n* Item 2\n"
	renderer := &Renderer{}
	output := bf.Run([]byte(input), bf.WithRenderer(renderer))
	result := string(output)

	if !strings.Contains(result, "TEXT=\"Item 1\"") {
		t.Error("Output should contain Item 1 as node TEXT")
	}
	if !strings.Contains(result, "TEXT=\"Item 2\"") {
		t.Error("Output should contain Item 2 as node TEXT")
	}
}

// TestLinkRendering tests that links are rendered correctly.
func TestLinkRendering(t *testing.T) {
	resetRendererState()

	input := "* [Link Text](http://example.com)\n"
	renderer := &Renderer{}
	output := bf.Run([]byte(input), bf.WithRenderer(renderer))
	result := string(output)

	if !strings.Contains(result, "TEXT=\"Link Text\"") {
		t.Error("Output should contain link text as node TEXT")
	}
	if !strings.Contains(result, "LINK=\"http://example.com\"") {
		t.Error("Output should contain link URL as LINK attribute")
	}
}

// TestCodeBlockRendering tests that code blocks are rendered correctly.
func TestCodeBlockRendering(t *testing.T) {
	resetRendererState()

	input := "```\ncode block\n```\n"
	renderer := &Renderer{}
	output := bf.Run([]byte(input), bf.WithRenderer(renderer))
	result := string(output)

	if !strings.Contains(result, "<pre><code>") {
		t.Error("Output should contain pre/code elements")
	}
	if !strings.Contains(result, "code block") {
		t.Error("Output should contain code block content")
	}
}

// TestTableRendering tests that tables are rendered correctly.
func TestTableRendering(t *testing.T) {
	resetRendererState()

	input := "| a | b |\n| --- | --- |\n| 1 | 2 |\n"
	renderer := &Renderer{}
	output := bf.Run([]byte(input), bf.WithRenderer(renderer), bf.WithExtensions(bf.CommonExtensions))
	result := string(output)

	if !strings.Contains(result, "<table>") {
		t.Error("Output should contain table element")
	}
	if !strings.Contains(result, "<th>a</th>") {
		t.Error("Output should contain table header")
	}
	if !strings.Contains(result, "<td>1</td>") {
		t.Error("Output should contain table cell")
	}
}

// TestBlockQuoteRendering tests that block quotes are rendered correctly.
func TestBlockQuoteRendering(t *testing.T) {
	resetRendererState()

	input := "> quoted text\n"
	renderer := &Renderer{}
	output := bf.Run([]byte(input), bf.WithRenderer(renderer))
	result := string(output)

	if !strings.Contains(result, "quoted text") {
		t.Error("Output should contain quoted text")
	}
	if !strings.Contains(result, "<richcontent TYPE=\"NODE\">") {
		t.Error("Output should contain richcontent element for blockquote")
	}
}

// resetRendererState resets the global renderer state between tests.
func resetRendererState() {
	renderingAsNodeElement = false
	headdingStack = []int{0}
}
