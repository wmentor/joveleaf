package joveleaf

import (
	"os"
	"strings"
	"testing"
)

func TestStringOrList_UnmarshalJSON_String(t *testing.T) {
	t.Parallel()

	data := []byte(`"hello world"`)
	var s StringOrList
	if err := s.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.String() != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", s.String())
	}
}

func TestStringOrList_UnmarshalJSON_List(t *testing.T) {
	t.Parallel()

	data := []byte(`["line1\n", "line2\n", "line3"]`)
	var s StringOrList
	if err := s.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "line1\nline2\nline3"
	if s.String() != expected {
		t.Fatalf("expected %q, got %q", expected, s.String())
	}
}

func TestParse_V4_Basic(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"metadata": {
			"language_info": {"name": "python", "version": "3.9"}
		},
		"nbformat": 4,
		"nbformat_minor": 5,
		"cells": [
			{
				"cell_type": "markdown",
				"metadata": {},
				"source": "# Hello"
			},
			{
				"cell_type": "code",
				"execution_count": 1,
				"metadata": {},
				"source": ["print(", "'hello')"],
				"outputs": [
					{
						"output_type": "stream",
						"name": "stdout",
						"text": "hello\n"
					}
				]
			}
		]
	}`)

	nb, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nb.NBFormat != 4 {
		t.Fatalf("expected nbformat 4, got %d", nb.NBFormat)
	}
	if len(nb.Cells) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(nb.Cells))
	}
	if nb.Cells[0].CellType != "markdown" {
		t.Fatalf("expected markdown cell, got %s", nb.Cells[0].CellType)
	}
	if nb.Cells[0].Source.String() != "# Hello" {
		t.Fatalf("expected source %q, got %q", "# Hello", nb.Cells[0].Source.String())
	}
	if nb.Cells[1].CellType != "code" {
		t.Fatalf("expected code cell, got %s", nb.Cells[1].CellType)
	}
	if nb.Cells[1].Source.String() != "print('hello')" {
		t.Fatalf("expected source %q, got %q", "print('hello')", nb.Cells[1].Source.String())
	}
	if len(nb.Cells[1].Outputs) != 1 {
		t.Fatalf("expected 1 output, got %d", len(nb.Cells[1].Outputs))
	}
}

func TestParse_V3_Basic(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"metadata": {},
		"nbformat": 3,
		"nbformat_minor": 0,
		"worksheets": [
			{
				"cells": [
					{
						"cell_type": "heading",
						"level": 2,
						"source": "My Heading"
					},
					{
						"cell_type": "code",
						"language": "python",
						"prompt_number": 1,
						"source": "x = 1",
						"outputs": [
							{
								"output_type": "stream",
								"stream": "stdout",
								"text": "done\n"
							}
						]
					}
				]
			}
		]
	}`)

	nb, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nb.NBFormat != 4 {
		t.Fatalf("expected normalized to nbformat 4, got %d", nb.NBFormat)
	}
	if len(nb.Cells) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(nb.Cells))
	}
	// Heading should be converted to markdown with ##
	if nb.Cells[0].CellType != "markdown" {
		t.Fatalf("expected markdown cell, got %s", nb.Cells[0].CellType)
	}
	if !strings.HasPrefix(nb.Cells[0].Source.String(), "## ") {
		t.Fatalf("expected heading prefix, got %q", nb.Cells[0].Source.String())
	}
	// Code cell output stream should be normalized
	if nb.Cells[1].Outputs[0].Name != "stdout" {
		t.Fatalf("expected stream name 'stdout', got %q", nb.Cells[1].Outputs[0].Name)
	}
}

func TestParse_UnsupportedVersion(t *testing.T) {
	t.Parallel()

	data := []byte(`{"nbformat": 2, "cells": []}`)
	_, err := Parse(data)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
}

func TestToMarkdown_Basic(t *testing.T) {
	t.Parallel()

	nb := &Notebook{
		NBFormat:      4,
		NBFormatMinor: 5,
		Metadata: Metadata{
			LanguageInfo: &LanguageInfo{Name: "python"},
		},
		Cells: []Cell{
			{
				CellType: "markdown",
				Source:   "# Title\n\nSome text",
			},
			{
				CellType:       "code",
				Source:         "print('hello')",
				ExecutionCount: new(1),
				Outputs: []Output{
					{
						OutputType: "stream",
						Name:       "stdout",
						Text:       "hello\n",
					},
				},
			},
		},
	}

	md := ToMarkdown(nb)

	if !strings.Contains(md, "# Title") {
		t.Fatal("expected markdown content in output")
	}
	if !strings.Contains(md, "**In [1]:**") {
		t.Fatal("expected execution count label")
	}
	if !strings.Contains(md, "```python") {
		t.Fatal("expected python code fence")
	}
	if !strings.Contains(md, "print('hello')") {
		t.Fatal("expected code in output")
	}
	if !strings.Contains(md, "**stdout:**") {
		t.Fatal("expected stdout output")
	}
	if !strings.Contains(md, "---") {
		t.Fatal("expected cell separator")
	}
}

func TestToMarkdown_Error(t *testing.T) {
	t.Parallel()

	nb := &Notebook{
		NBFormat: 4,
		Cells: []Cell{
			{
				CellType: "code",
				Source:   "1/0",
				Outputs: []Output{
					{
						OutputType: "error",
						EnName:     "ZeroDivisionError",
						EnValue:    "division by zero",
						Traceback:  []string{"line 1", "  1/0"},
					},
				},
			},
		},
	}

	md := ToMarkdown(nb)

	if !strings.Contains(md, "ZeroDivisionError") {
		t.Fatal("expected error name in output")
	}
	if !strings.Contains(md, "division by zero") {
		t.Fatal("expected error value in output")
	}
}

func TestWriteMimeBundle_SVG(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"image/svg+xml": "<svg>test</svg>",
		"text/plain":    "<graphviz.dot.Digraph at 0x123>",
	}

	var b strings.Builder
	writeMimeBundle(&b, data, nil)

	result := b.String()
	if !strings.Contains(result, "<svg>test</svg>") {
		t.Fatal("expected SVG content")
	}
	if strings.Contains(result, "<graphviz.dot.Digraph") {
		t.Fatal("should not include useless text/plain repr when SVG is present")
	}
}

func TestWriteMimeBundle_HTML(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"text/html":  "<table><tr><td>data</td></tr></table>",
		"text/plain": "<pandas.DataFrame>",
	}

	var b strings.Builder
	writeMimeBundle(&b, data, nil)

	result := b.String()
	if !strings.Contains(result, "<table>") {
		t.Fatal("expected HTML content")
	}
	if strings.Contains(result, "<pandas.DataFrame>") {
		t.Fatal("should not include useless text/plain repr when HTML is present")
	}
}

func TestWriteMimeBundle_TextPlain_NoRaster(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"text/plain": "some actual data",
	}

	var b strings.Builder
	writeMimeBundle(&b, data, nil)

	result := b.String()
	if !strings.Contains(result, "some actual data") {
		t.Fatal("expected text/plain content when no raster images")
	}
}

func TestWriteMimeBundle_TextPlain_WithRaster(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"text/plain": "<Figure size 720x576>",
		"image/png":  "base64data",
	}

	var b strings.Builder
	writeMimeBundle(&b, data, nil)

	result := b.String()
	if strings.Contains(result, "<Figure size") {
		t.Fatal("should not include text/plain repr when raster image is present")
	}
	if strings.Contains(result, "base64data") {
		t.Fatal("should not include raster image data")
	}
	// Should be empty output
	if result != "" {
		t.Fatalf("expected empty output, got %q", result)
	}
}

func TestWriteMimeBundle_JSON(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"application/json": `{"key": "value"}`,
	}

	var b strings.Builder
	writeMimeBundle(&b, data, nil)

	result := b.String()
	if !strings.Contains(result, "```json") {
		t.Fatal("expected json code block")
	}
	if !strings.Contains(result, `{"key": "value"}`) {
		t.Fatal("expected JSON content")
	}
}

func TestWriteMimeBundle_ObjectMimeType(t *testing.T) {
	t.Parallel()

	// Test that object mime types (like widget-view+json) are handled gracefully
	data := map[string]any{
		"application/vnd.jupyter.widget-view+json": map[string]any{
			"model_id": "55e9b3adbda8474e87f39a259b6c1602",
		},
	}

	var b strings.Builder
	writeMimeBundle(&b, data, nil)

	// Should output as JSON
	result := b.String()
	if !strings.Contains(result, "model_id") {
		t.Fatal("expected model_id in output")
	}
}

func TestStripANSI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"\033[31mred\033[0m", "red"},
		{"\033[1;32mbold green\033[0m", "bold green"},
		{"no ansi here", "no ansi here"},
	}

	for _, tt := range tests {
		result := stripANSI(tt.input)
		if result != tt.expected {
			t.Errorf("stripANSI(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDetectLanguage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		nb       *Notebook
		expected string
	}{
		{
			name: "from language_info",
			nb: &Notebook{
				Metadata: Metadata{LanguageInfo: &LanguageInfo{Name: "python"}},
			},
			expected: "python",
		},
		{
			name: "from kernelspec",
			nb: &Notebook{
				Metadata: Metadata{Kernelspec: &Kernelspec{Language: "javascript"}},
			},
			expected: "javascript",
		},
		{
			name:     "empty metadata",
			nb:       &Notebook{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := detectLanguage(tt.nb)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParse_TestdataFile(t *testing.T) {
	t.Parallel()

	nb, err := ReadFile("testdata/test001.ipynb")
	if err != nil {
		t.Fatalf("failed to read testdata/test001.ipynb: %v", err)
	}

	if nb.NBFormat != 4 {
		t.Fatalf("expected nbformat 4, got %d", nb.NBFormat)
	}

	if len(nb.Cells) == 0 {
		t.Fatal("expected at least one cell")
	}

	// Verify language detection works
	lang := detectLanguage(nb)
	if lang != "python" {
		t.Fatalf("expected language 'python', got %q", lang)
	}

	// Verify markdown conversion works
	md := ToMarkdown(nb)
	if md == "" {
		t.Fatal("expected non-empty markdown output")
	}
}

func TestParse_TestdataFile2(t *testing.T) {
	t.Parallel()

	nb, err := ReadFile("testdata/test002.ipynb")
	if err != nil {
		t.Fatalf("failed to read testdata/test002.ipynb: %v", err)
	}

	// test002 is nbformat v3, should be normalized to v4
	if nb.NBFormat != 4 {
		t.Fatalf("expected normalized to nbformat 4, got %d", nb.NBFormat)
	}

	if len(nb.Cells) == 0 {
		t.Fatal("expected at least one cell")
	}

	md := ToMarkdown(nb)
	if md == "" {
		t.Fatal("expected non-empty markdown output")
	}

	// Verify In[N] labels are present
	if !strings.Contains(md, "**In [1]:**") {
		t.Fatal("expected execution count label In [1]")
	}
}

func TestParse_TestdataFile3(t *testing.T) {
	t.Parallel()

	nb, err := ReadFile("testdata/test003.ipynb")
	if err != nil {
		t.Fatalf("failed to read testdata/test003.ipynb: %v", err)
	}

	if nb.NBFormat != 4 {
		t.Fatalf("expected nbformat 4, got %d", nb.NBFormat)
	}

	if len(nb.Cells) == 0 {
		t.Fatal("expected at least one cell")
	}

	md := ToMarkdown(nb)
	if md == "" {
		t.Fatal("expected non-empty markdown output")
	}

	// Verify object mime types (like widget-view+json) are handled
	if !strings.Contains(md, "model_id") {
		t.Fatal("expected model_id from widget-view+json output")
	}
}

func TestReadToMarkdown(t *testing.T) {
	t.Parallel()

	md, err := ReadToMarkdown("testdata/test001.ipynb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if md == "" {
		t.Fatal("expected non-empty markdown output")
	}
	if !strings.Contains(md, "# Title") && !strings.Contains(md, "Chapter") {
		t.Fatal("expected markdown content in output")
	}
}

func TestReadToMarkdown_InvalidPath(t *testing.T) {
	t.Parallel()

	_, err := ReadToMarkdown("nonexistent.ipynb")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestReadToMarkdown_GoldenFiles(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input  string
		golden string
	}{
		{"testdata/test001.ipynb", "testdata/result001.md"},
		{"testdata/test002.ipynb", "testdata/result002.md"},
		{"testdata/test003.ipynb", "testdata/result003.md"},
		{"testdata/test004.ipynb", "testdata/result004.md"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			// Read and convert the notebook
			actual, err := ReadToMarkdown(tc.input)
			if err != nil {
				t.Fatalf("ReadToMarkdown(%q) failed: %v", tc.input, err)
			}

			// Read the golden file
			expectedBytes, err := os.ReadFile(tc.golden)
			if err != nil {
				t.Fatalf("failed to read golden file %q: %v", tc.golden, err)
			}
			expected := string(expectedBytes)

			// Compare
			if actual != expected {
				t.Errorf("output for %q does not match %q", tc.input, tc.golden)
				t.Errorf("expected length: %d, actual length: %d", len(expected), len(actual))

				// Find first difference
				for i := 0; i < len(expected) && i < len(actual); i++ {
					if expected[i] != actual[i] {
						start := max(i-50, 0)
						end := min(i+50, len(expected))
						t.Errorf("first difference at position %d", i)
						t.Errorf("expected around diff: %q", expected[start:end])
						t.Errorf("actual around diff:   %q", actual[start:end])
						break
					}
				}
			}
		})
	}
}

//go:fix inline
func intPtr(v int) *int {
	return new(v)
}
