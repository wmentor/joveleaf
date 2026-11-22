package joveleaf

import "github.com/pkg/errors"

var (
	// ErrInvalidNotebook is returned when the notebook JSON is malformed or missing required fields.
	ErrInvalidNotebook = errors.New("invalid notebook")

	// ErrUnsupportedVersion is returned when the notebook format version is not supported.
	ErrUnsupportedVersion = errors.New("unsupported notebook version")

	// ErrEmptySource is returned when a required cell source is empty.
	ErrEmptySource = errors.New("empty cell source")
)
