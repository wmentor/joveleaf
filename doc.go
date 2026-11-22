// Package joveleaf reads Jupyter Notebook (.ipynb) files and converts their
// content to Markdown.
//
// It supports both nbformat v3 and v4, automatically normalizing v3 notebooks
// to the v4 structure internally.
//
// Basic usage:
//
//	nb, err := joveleaf.ReadFile("notebook.ipynb")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	md := joveleaf.ToMarkdown(nb)
//	fmt.Println(md)
package joveleaf
