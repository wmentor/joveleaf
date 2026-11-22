# joveleaf

A Go library for reading Jupyter Notebook (`.ipynb`) files and converting their content to Markdown.

## Features

- **Read `.ipynb` files** — parse notebook JSON into Go structures
- **nbformat v3 + v4 support** — v3 notebooks are automatically normalized to v4
- **Convert to Markdown** — cells, code blocks, outputs (stdout, stderr, errors, text, HTML, LaTeX, SVG, JSON)
- **Multi-line string handling** — transparently joins string lists (notebook convention)
- **Language detection** — extracts language from metadata for proper code fences
- **ANSI stripping** — cleans terminal color codes from tracebacks

## Installation

```bash
go get github.com/wmentor/joveleaf
```

## Usage

### Command-line tool

Install the CLI:

```bash
go install github.com/wmentor/joveleaf/cmd/jnb2md@latest
```

Convert notebooks:

```bash
# Single file
jnb2md notebook.ipynb > output.md

# Multiple files (separated by ---)
jnb2md file1.ipynb file2.ipynb > combined.md

# Output to file
jnb2md -o output.md notebook.ipynb

# Read from stdin
cat notebook.ipynb | jnb2md

# Custom separator
jnb2md -sep "\n\n" file1.ipynb file2.ipynb
```

Flags:

- `-o string` — output file path (default: stdout)
- `-sep string` — separator between notebooks (default: `\n---\n\n`)
- `-h, -help` — show help

### Library

```go
package main

import (
    "fmt"
    "log"

    "github.com/wmentor/joveleaf"
)

func main() {
    // Quick conversion: read file and convert to markdown in one call
    md, err := joveleaf.ReadToMarkdown("notebook.ipynb")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(md)
    
    // Or use separate steps for more control
    nb, err := joveleaf.ReadFile("notebook.ipynb")
    if err != nil {
        log.Fatal(err)
    }
    
    md = joveleaf.ToMarkdown(nb)
    fmt.Println(md)
}
```

## API

| Function | Description |
|---|---|
| `ReadToMarkdown(path)` | Read notebook from file and convert to Markdown string |
| `ReadFile(path)` | Read notebook from file path |
| `Read(reader)` | Read notebook from `io.Reader` |
| `Parse(data)` | Parse notebook from `[]byte` |
| `ToMarkdown(nb)` | Convert notebook to Markdown string |

## Output mapping

| Notebook element | Markdown output |
|---|---|
| Markdown cell | Markdown text as-is |
| Code cell | Fenced code block with detected language |
| Raw cell | Content as-is |
| `stream` output | `stdout`/`stderr` in code block |
| `execute_result` | `Out[N]:` label + `text/plain` |
| `display_data` | Best available: `text/plain` > `text/html` > `text/latex` > `image/svg+xml` |
| `error` output | Exception name, value, and traceback |
| Images (png, jpeg) | Ignored |

## License

MIT
