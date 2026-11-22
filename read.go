package joveleaf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// ReadFile reads a Jupyter notebook from the given file path.
// It supports both nbformat v3 and v4, normalizing v3 to v4 internally.
func ReadFile(path string) (*Notebook, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open notebook file")
	}
	defer f.Close()

	return Read(f)
}

// Read reads a Jupyter notebook from the given reader.
// It supports both nbformat v3 and v4, normalizing v3 to v4 internally.
func Read(r io.Reader) (*Notebook, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read notebook data")
	}

	return Parse(data)
}

// Parse parses a Jupyter notebook from raw JSON bytes.
// It supports both nbformat v3 and v4, normalizing v3 to v4 internally.
func Parse(data []byte) (*Notebook, error) {
	// First, detect the version
	var version struct {
		NBFormat int `json:"nbformat"`
	}
	if err := json.Unmarshal(data, &version); err != nil {
		return nil, errors.Wrap(ErrInvalidNotebook, err.Error())
	}

	switch version.NBFormat {
	case 4:
		return parseV4(data)
	case 3:
		return parseV3(data)
	default:
		return nil, errors.Wrap(ErrUnsupportedVersion,
			fmt.Sprintf("nbformat %d is not supported", version.NBFormat))
	}
}

// parseV4 parses a nbformat v4 notebook.
func parseV4(data []byte) (*Notebook, error) {
	var nb Notebook
	if err := json.Unmarshal(data, &nb); err != nil {
		return nil, errors.Wrap(ErrInvalidNotebook, err.Error())
	}
	if nb.NBFormat != 4 {
		return nil, errors.Wrap(ErrInvalidNotebook,
			fmt.Sprintf("expected nbformat 4, got %d", nb.NBFormat))
	}
	return &nb, nil
}

// ReadToMarkdown reads a Jupyter notebook from the given file path
// and converts it to a Markdown string.
// It is a convenience wrapper around ReadFile and ToMarkdown.
func ReadToMarkdown(path string) (string, error) {
	nb, err := ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "read notebook")
	}
	return ToMarkdown(nb), nil
}
