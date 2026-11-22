package joveleaf

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// notebookV3Raw represents the raw structure of a nbformat v3 notebook.
type notebookV3Raw struct {
	Metadata    json.RawMessage `json:"metadata"`
	NBFormat    int             `json:"nbformat"`
	NBFormatMinor int           `json:"nbformat_minor"`
	Worksheets  []worksheetV3   `json:"worksheets"`
}

type worksheetV3 struct {
	Cells []cellV3 `json:"cells"`
}

type cellV3 struct {
	CellType       string          `json:"cell_type"`
	Source         StringOrList    `json:"source"`
	Input          StringOrList    `json:"input"` // v3 uses "input" instead of "source"
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	// Code cell fields
	Language       string          `json:"language,omitempty"`
	Collapsed      *bool           `json:"collapsed,omitempty"`
	Outputs        []outputV3      `json:"outputs,omitempty"`
	PromptNumber   *int            `json:"prompt_number,omitempty"`
	// Heading cell fields
	Level          int             `json:"level,omitempty"`
	// Markdown cell fields (same as source)
}

type outputV3 struct {
	OutputType     string                 `json:"output_type"`
	// stream
	Stream         string                 `json:"stream,omitempty"`
	Text           StringOrList           `json:"text,omitempty"`
	// display_data / pyout
	MimeType       string                 `json:"mime_type,omitempty"`
	Data           map[string]StringOrList `json:"data,omitempty"`
	PromptNumber   *int                   `json:"prompt_number,omitempty"`
	// png, jpeg, etc. (v3 stores images at top level)
	PNG            string                 `json:"png,omitempty"`
	JPEG           string                 `json:"jpeg,omitempty"`
	HTML           StringOrList           `json:"html,omitempty"`
	Latex          StringOrList           `json:"latex,omitempty"`
	SVG            StringOrList           `json:"svg,omitempty"`
	// error / pyerr
	EnName         string                 `json:"ename,omitempty"`
	EnValue        string                 `json:"evalue,omitempty"`
	Traceback      []string               `json:"traceback,omitempty"`
}

// parseV3 parses a nbformat v3 notebook and normalizes it to v4.
func parseV3(data []byte) (*Notebook, error) {
	var raw notebookV3Raw
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, errors.Wrap(ErrInvalidNotebook, err.Error())
	}

	nb := &Notebook{
		NBFormat:      4,
		NBFormatMinor: 0,
	}

	// Parse metadata (best effort)
	if len(raw.Metadata) > 0 {
		_ = json.Unmarshal(raw.Metadata, &nb.Metadata)
	}

	// Convert worksheets to cells
	for _, ws := range raw.Worksheets {
		for _, cellV3 := range ws.Cells {
			cell, err := convertCellV3toV4(cellV3)
			if err != nil {
				return nil, errors.Wrap(err, "convert v3 cell")
			}
			nb.Cells = append(nb.Cells, *cell)
		}
	}

	return nb, nil
}

// convertCellV3toV4 converts a v3 cell to v4 format.
func convertCellV3toV4(c cellV3) (*Cell, error) {
	cell := &Cell{}

	// In v3, code cells use "input" instead of "source"
	if c.Source != "" {
		cell.Source = c.Source
	} else {
		cell.Source = c.Input
	}

	switch c.CellType {
	case "code":
		cell.CellType = "code"
		cell.ExecutionCount = c.PromptNumber
		if c.Collapsed != nil {
			cell.Metadata.Collapsed = c.Collapsed
		}
		for _, outV3 := range c.Outputs {
			out := convertOutputV3toV4(outV3)
			cell.Outputs = append(cell.Outputs, *out)
		}

	case "markdown":
		cell.CellType = "markdown"

	case "heading":
		// In v3, headings are separate cells with a level.
		// In v4, they are markdown cells with # headings.
		cell.CellType = "markdown"
		prefix := ""
		for i := 0; i < c.Level; i++ {
			prefix += "#"
		}
		cell.Source = StringOrList(fmt.Sprintf("%s %s", prefix, string(cell.Source)))

	case "raw":
		cell.CellType = "raw"

	default:
		// Unknown cell type, treat as raw
		cell.CellType = "raw"
	}

	return cell, nil
}

// convertOutputV3toV4 converts a v3 output to v4 format.
func convertOutputV3toV4(o outputV3) *Output {
	out := &Output{}

	switch o.OutputType {
	case "stream":
		out.OutputType = "stream"
		out.Name = o.Stream
		out.Text = o.Text

	case "pyout":
		// v3 pyout → v4 execute_result
		out.OutputType = "execute_result"
		out.ExecutionCount = o.PromptNumber
		out.Data = buildDataBundleV3(o)

	case "display_data":
		out.OutputType = "display_data"
		out.Data = buildDataBundleV3(o)

	case "pyerr":
		// v3 pyerr → v4 error
		out.OutputType = "error"
		out.EnName = o.EnName
		out.EnValue = o.EnValue
		out.Traceback = o.Traceback

	default:
		// Unknown output type, try to preserve as stream
		out.OutputType = "stream"
		out.Name = "stdout"
		out.Text = o.Text
	}

	return out
}

// buildDataBundleV3 builds a v4 mime-bundle from v3 output fields.
// In v3, mime types can be stored either in the "data" map or at the top level
// of the output object (e.g., "text", "html", "png" as direct keys).
func buildDataBundleV3(o outputV3) map[string]interface{} {
	data := make(map[string]interface{})

	// First, copy from the "data" map if present
	if o.Data != nil {
		for k, v := range o.Data {
			data[k] = v
		}
	}

	// v3 stores some mime types at top level — add them if not already in data
	if o.Text != "" {
		if _, exists := data["text/plain"]; !exists {
			data["text/plain"] = o.Text
		}
	}
	if o.HTML != "" {
		if _, exists := data["text/html"]; !exists {
			data["text/html"] = o.HTML
		}
	}
	if o.Latex != "" {
		if _, exists := data["text/latex"]; !exists {
			data["text/latex"] = o.Latex
		}
	}
	if o.SVG != "" {
		if _, exists := data["image/svg+xml"]; !exists {
			data["image/svg+xml"] = o.SVG
		}
	}
	if o.PNG != "" {
		if _, exists := data["image/png"]; !exists {
			data["image/png"] = string(o.PNG)
		}
	}
	if o.JPEG != "" {
		if _, exists := data["image/jpeg"]; !exists {
			data["image/jpeg"] = string(o.JPEG)
		}
	}

	return data
}
