package joveleaf

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ToMarkdown converts a parsed Notebook to a Markdown string.
func ToMarkdown(nb *Notebook) string {
	var b strings.Builder

	// Detect language from metadata
	lang := detectLanguage(nb)

	for i, cell := range nb.Cells {
		if i > 0 {
			b.WriteString("\n---\n\n")
		}
		writeCell(&b, &cell, lang)
	}

	return b.String()
}

// detectLanguage extracts the programming language from notebook metadata.
func detectLanguage(nb *Notebook) string {
	if nb.Metadata.LanguageInfo != nil && nb.Metadata.LanguageInfo.Name != "" {
		return normalizeLang(nb.Metadata.LanguageInfo.Name)
	}
	if nb.Metadata.Kernelspec != nil && nb.Metadata.Kernelspec.Language != "" {
		return normalizeLang(nb.Metadata.Kernelspec.Language)
	}
	return ""
}

// normalizeLang normalizes language names to common markdown code fence identifiers.
func normalizeLang(lang string) string {
	switch strings.ToLower(lang) {
	case "python", "python3":
		return "python"
	case "r":
		return "r"
	case "julia":
		return "julia"
	case "javascript", "node":
		return "javascript"
	case "typescript":
		return "typescript"
	case "scala":
		return "scala"
	case "java":
		return "java"
	case "c", "c++", "cpp":
		return "cpp"
	case "go":
		return "go"
	case "ruby":
		return "ruby"
	case "bash", "shell":
		return "bash"
	default:
		return strings.ToLower(lang)
	}
}

// writeCell writes a single cell as markdown.
func writeCell(b *strings.Builder, cell *Cell, lang string) {
	switch cell.CellType {
	case "markdown":
		writeMarkdownCell(b, cell)
	case "code":
		writeCodeCell(b, cell, lang)
	case "raw":
		writeRawCell(b, cell)
	default:
		// Unknown cell type, write source as-is
		b.WriteString(cell.Source.String())
		b.WriteString("\n")
	}
}

// writeMarkdownCell writes a markdown cell.
func writeMarkdownCell(b *strings.Builder, cell *Cell) {
	b.WriteString(cell.Source.String())
	b.WriteString("\n")
}

// writeCodeCell writes a code cell with its outputs.
func writeCodeCell(b *strings.Builder, cell *Cell, lang string) {
	// Write "In [N]:" label if execution count is present
	if cell.ExecutionCount != nil {
		b.WriteString(fmt.Sprintf("**In [%d]:**\n\n", *cell.ExecutionCount))
	}

	// Write code block
	b.WriteString("```")
	b.WriteString(lang)
	b.WriteString("\n")
	b.WriteString(cell.Source.String())
	b.WriteString("\n```\n")

	// Write outputs
	for _, output := range cell.Outputs {
		writeOutput(b, &output, cell.ExecutionCount)
	}
}

// writeRawCell writes a raw cell.
func writeRawCell(b *strings.Builder, cell *Cell) {
	// Raw cells are included as-is
	b.WriteString(cell.Source.String())
	b.WriteString("\n")
}

// writeOutput writes a single output as markdown.
func writeOutput(b *strings.Builder, output *Output, execCount *int) {
	switch output.OutputType {
	case "stream":
		writeStreamOutput(b, output)
	case "display_data":
		writeDisplayDataOutput(b, output)
	case "execute_result":
		writeExecuteResultOutput(b, output, execCount)
	case "error":
		writeErrorOutput(b, output)
	}
}

// writeStreamOutput writes stream output (stdout/stderr).
func writeStreamOutput(b *strings.Builder, output *Output) {
	text := output.Text.String()
	if text == "" {
		return
	}

	b.WriteString("\n**")
	if output.Name == "stderr" {
		b.WriteString("stderr")
	} else {
		b.WriteString("stdout")
	}
	b.WriteString(":**\n\n```\n")
	b.WriteString(text)
	if !strings.HasSuffix(text, "\n") {
		b.WriteString("\n")
	}
	b.WriteString("```\n")
}

// writeDisplayDataOutput writes display_data output.
func writeDisplayDataOutput(b *strings.Builder, output *Output) {
	writeMimeBundle(b, output.Data, nil)
}

// writeExecuteResultOutput writes execute_result output.
func writeExecuteResultOutput(b *strings.Builder, output *Output, execCount *int) {
	count := output.ExecutionCount
	if count == nil {
		count = execCount
	}

	if count != nil {
		b.WriteString(fmt.Sprintf("\n**Out[%d]:**\n\n", *count))
	} else {
		b.WriteString("\n**Out:**\n\n")
	}

	writeMimeBundle(b, output.Data, count)
}

// writeErrorOutput writes error output.
func writeErrorOutput(b *strings.Builder, output *Output) {
	b.WriteString("\n**Error:** ")
	b.WriteString(output.EnName)
	if output.EnValue != "" {
		b.WriteString(": ")
		b.WriteString(output.EnValue)
	}
	b.WriteString("\n\n")

	if len(output.Traceback) > 0 {
		b.WriteString("```\n")
		for _, line := range output.Traceback {
			// Strip ANSI color codes from traceback
			cleaned := stripANSI(line)
			b.WriteString(cleaned)
			if !strings.HasSuffix(cleaned, "\n") {
				b.WriteString("\n")
			}
		}
		b.WriteString("```\n")
	}
}

// writeMimeBundle writes a mime-bundle, selecting the best available format.
// Priority: image/svg+xml > text/html > text/latex > text/plain (only if no raster images)
// Raster images (png, jpeg) are ignored as per requirements.
// text/plain is skipped when raster images are present, because in that case
// text/plain is typically just a Python repr like "<Figure size 720x576>".
func writeMimeBundle(b *strings.Builder, data map[string]interface{}, execCount *int) {
	if len(data) == 0 {
		return
	}

	hasRaster := false
	if _, ok := data["image/png"]; ok {
		hasRaster = true
	}
	if _, ok := data["image/jpeg"]; ok {
		hasRaster = true
	}

	// Try rich mime types first
	type mimeHandler struct {
		mime  string
		lang  string
	}
	handlers := []mimeHandler{
		{"image/svg+xml", "svg"},
		{"text/html", "html"},
		{"text/latex", "latex"},
	}

	for _, h := range handlers {
		if content, ok := data[h.mime]; ok {
			text := extractString(content)
			if text == "" {
				continue
			}
			b.WriteString("```" + h.lang + "\n")
			b.WriteString(text)
			if !strings.HasSuffix(text, "\n") {
				b.WriteString("\n")
			}
			b.WriteString("```\n")
			return
		}
	}

	// text/plain: only show if there are no raster images
	// (otherwise it is just a useless Python repr)
	if !hasRaster {
		if content, ok := data["text/plain"]; ok {
			text := extractString(content)
			if text != "" {
				b.WriteString("```\n")
				b.WriteString(text)
				if !strings.HasSuffix(text, "\n") {
					b.WriteString("\n")
				}
				b.WriteString("```\n")
				return
			}
		}
	}

	// Fallback: application/json or any mime type ending with +json
	if content, ok := data["application/json"]; ok {
		text := extractString(content)
		if text != "" {
			b.WriteString("```json\n")
			b.WriteString(text)
			if !strings.HasSuffix(text, "\n") {
				b.WriteString("\n")
			}
			b.WriteString("```\n")
			return
		}
	}

	// Check for any mime type ending with +json
	for mime, content := range data {
		if strings.HasSuffix(mime, "+json") {
			text := extractString(content)
			if text != "" {
				b.WriteString("```json\n")
				b.WriteString(text)
				if !strings.HasSuffix(text, "\n") {
					b.WriteString("\n")
				}
				b.WriteString("```\n")
				return
			}
		}
	}
}

// extractString extracts a string from an interface{} value.
// Handles string, []string (joined), []interface{} (joined), and other types.
func extractString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []string:
		return strings.Join(val, "")
	case []interface{}:
		parts := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, "")
	default:
		// For other types (objects, numbers, etc.), convert to JSON string
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(jsonBytes)
	}
}

// stripANSI removes ANSI escape codes from a string.
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false

	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}

	return result.String()
}
