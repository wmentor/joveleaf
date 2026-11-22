// Command jnb2md converts Jupyter Notebook (.ipynb) files to Markdown.
//
// Usage:
//
//	jnb2md [flags] [file.ipynb ...]
//
// If no files are specified, input is read from stdin.
// Output is written to stdout unless -o is specified.
//
// Flags:
//
//	-o string   output file path (default: stdout)
//	-sep        separator between notebooks when processing multiple files (default: "\n---\n\n")
//	-h, -help   show help message
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/wmentor/joveleaf"
)

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "jnb2md: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("jnb2md", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		outputPath = fs.String("o", "", "output file path (default: stdout)")
		sep        = fs.String("sep", "\n---\n\n", "separator between notebooks")
		showHelp   = fs.Bool("help", false, "show help message")
		showHelpH  = fs.Bool("h", false, "show help message")
	)

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: jnb2md [flags] [file.ipynb ...]\n\n")
		fmt.Fprintf(stderr, "Convert Jupyter Notebook (.ipynb) files to Markdown.\n\n")
		fmt.Fprintf(stderr, "If no files are specified, input is read from stdin.\n")
		fmt.Fprintf(stderr, "Output is written to stdout unless -o is specified.\n\n")
		fmt.Fprintf(stderr, "Flags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	if *showHelp || *showHelpH {
		fs.Usage()
		return nil
	}

	files := fs.Args()

	// Determine output writer
	var out io.Writer = stdout
	var outFile *os.File
	if *outputPath != "" {
		f, err := os.Create(*outputPath)
		if err != nil {
			return errors.Wrap(err, "create output file")
		}
		defer f.Close()
		out = f
		outFile = f
	}

	// If no files specified, read from stdin
	if len(files) == 0 {
		return convertReader(stdin, out, "stdin")
	}

	// Process each file
	for i, path := range files {
		if i > 0 && *sep != "" {
			if _, err := io.WriteString(out, *sep); err != nil {
				return errors.Wrap(err, "write separator")
			}
		}

		if err := convertFile(path, out); err != nil {
			return errors.Wrapf(err, "convert %s", path)
		}
	}

	// Close output file if we created one
	if outFile != nil {
		if err := outFile.Close(); err != nil {
			return errors.Wrap(err, "close output file")
		}
	}

	return nil
}

// convertFile reads a notebook file and writes its markdown to the writer.
func convertFile(path string, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	defer f.Close()

	return convertReader(f, w, filepath.Base(path))
}

// convertReader reads a notebook from the reader and writes markdown to the writer.
func convertReader(r io.Reader, w io.Writer, name string) error {
	nb, err := joveleaf.Read(r)
	if err != nil {
		return errors.Wrapf(err, "read notebook %s", name)
	}

	md := joveleaf.ToMarkdown(nb)

	if _, err := io.WriteString(w, md); err != nil {
		return errors.Wrap(err, "write markdown")
	}

	// Ensure trailing newline
	if !strings.HasSuffix(md, "\n") {
		if _, err := io.WriteString(w, "\n"); err != nil {
			return errors.Wrap(err, "write trailing newline")
		}
	}

	return nil
}
