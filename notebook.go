package joveleaf

// Notebook represents a parsed Jupyter notebook (nbformat v4 after normalization).
type Notebook struct {
	Metadata    Metadata `json:"metadata"`
	NBFormat    int      `json:"nbformat"`
	NBFormatMinor int    `json:"nbformat_minor"`
	Cells       []Cell   `json:"cells"`
}

// Metadata holds notebook-level metadata.
type Metadata struct {
	KernelInfo    *KernelInfo    `json:"kernel_info,omitempty"`
	LanguageInfo  *LanguageInfo  `json:"language_info,omitempty"`
	Kernelspec    *Kernelspec    `json:"kernelspec,omitempty"`
	Authors       []Author       `json:"authors,omitempty"`
}

// KernelInfo describes the kernel.
type KernelInfo struct {
	Name string `json:"name"`
}

// LanguageInfo describes the programming language.
type LanguageInfo struct {
	Name           string      `json:"name"`
	Version        string      `json:"version,omitempty"`
	CodemirrorMode interface{} `json:"codemirror_mode,omitempty"` // string or object
}

// Kernelspec is the kernel specification.
type Kernelspec struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Language    string `json:"language,omitempty"`
}

// Author represents a notebook author.
type Author struct {
	Name string `json:"name"`
}

// Cell represents a single notebook cell.
type Cell struct {
	CellType       string         `json:"cell_type"`
	Source         StringOrList   `json:"source"`
	Metadata       CellMetadata   `json:"metadata,omitempty"`
	// Code cell fields
	ExecutionCount *int           `json:"execution_count,omitempty"`
	Outputs        []Output       `json:"outputs,omitempty"`
	// Raw/Markdown cell field
	Attachments    map[string]map[string]interface{} `json:"attachments,omitempty"`
	// Cell id (since v4.5)
	ID             string         `json:"id,omitempty"`
}

// CellMetadata holds cell-level metadata.
type CellMetadata struct {
	Collapsed    *bool    `json:"collapsed,omitempty"`
	Scrolled     interface{} `json:"scrolled,omitempty"` // bool or "auto"
	Deletable    *bool    `json:"deletable,omitempty"`
	Editable     *bool    `json:"editable,omitempty"`
	Format       string   `json:"format,omitempty"`
	Name         string   `json:"name,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Jupyter      *JupyterMeta  `json:"jupyter,omitempty"`
	Execution    *ExecutionMeta `json:"execution,omitempty"`
}

// JupyterMeta holds jupyter-specific cell metadata.
type JupyterMeta struct {
	SourceHidden *bool `json:"source_hidden,omitempty"`
	OutputsHidden *bool `json:"outputs_hidden,omitempty"`
}

// ExecutionMeta holds execution timestamp metadata.
type ExecutionMeta struct {
	IOPubExecuteInput string `json:"iopub.execute_input,omitempty"`
	IOPubStatusBusy   string `json:"iopub.status.busy,omitempty"`
	ShellExecuteReply string `json:"shell.execute_reply,omitempty"`
	IOPubStatusIdle   string `json:"iopub.status.idle,omitempty"`
}

// Output represents a code cell output.
type Output struct {
	OutputType string `json:"output_type"`
	// stream output
	Name string       `json:"name,omitempty"`
	Text StringOrList `json:"text,omitempty"`
	// display_data / execute_result
	Data           map[string]interface{} `json:"data,omitempty"`
	ExecutionCount *int                   `json:"execution_count,omitempty"`
	// error output
	EnName    string   `json:"ename,omitempty"`
	EnValue   string   `json:"evalue,omitempty"`
	Traceback []string `json:"traceback,omitempty"`
}
