package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_Help(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	err := run([]string{"jnb2md", "-help"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatal("expected help output on stderr")
	}
}

func TestRun_Stdin(t *testing.T) {
	t.Parallel()

	input := `{
		"metadata": {"language_info": {"name": "python"}},
		"nbformat": 4,
		"nbformat_minor": 0,
		"cells": [
			{"cell_type": "markdown", "metadata": {}, "source": "# Hello"}
		]
	}`

	var stdout bytes.Buffer
	err := run([]string{"jnb2md"}, strings.NewReader(input), &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout.String(), "# Hello") {
		t.Fatal("expected markdown output")
	}
}

func TestRun_File(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inPath := filepath.Join(dir, "test.ipynb")
	outPath := filepath.Join(dir, "out.md")

	notebook := `{
		"metadata": {},
		"nbformat": 4,
		"nbformat_minor": 0,
		"cells": [
			{"cell_type": "code", "execution_count": 1, "metadata": {}, "source": "x = 1", "outputs": []}
		]
	}`

	if err := os.WriteFile(inPath, []byte(notebook), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	var stderr bytes.Buffer
	err := run([]string{"jnb2md", "-o", outPath, inPath}, nil, &bytes.Buffer{}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(data), "x = 1") {
		t.Fatal("expected code in output")
	}
}

func TestRun_MultipleFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file1 := filepath.Join(dir, "a.ipynb")
	file2 := filepath.Join(dir, "b.ipynb")

	nb := func(src string) string {
		return `{
			"metadata": {},
			"nbformat": 4,
			"nbformat_minor": 0,
			"cells": [
				{"cell_type": "markdown", "metadata": {}, "source": "` + src + `"}
			]
		}`
	}

	if err := os.WriteFile(file1, []byte(nb("AAA")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte(nb("BBB")), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	err := run([]string{"jnb2md", file1, file2}, nil, &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "AAA") || !strings.Contains(out, "BBB") {
		t.Fatal("expected both notebooks in output")
	}
	if !strings.Contains(out, "---") {
		t.Fatal("expected separator between notebooks")
	}
}

func TestRun_InvalidNotebook(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.ipynb")
	if err := os.WriteFile(bad, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stderr bytes.Buffer
	err := run([]string{"jnb2md", bad}, nil, &bytes.Buffer{}, &stderr)
	if err == nil {
		t.Fatal("expected error for invalid notebook")
	}
}
