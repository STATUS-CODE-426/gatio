package generator

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Options struct {
	DryRun         bool
	ForceOverWrite bool
	GoModuleName   string
}

type ProjectGenerator struct {
	opts       Options
	templateFs fs.FS
	infoLog    io.Writer
	errorLog   io.Writer
}

//go:embed template/*
var defaultFS embed.FS

func NewProjectGenerator(opts Options, userTemplateFs fs.FS) *ProjectGenerator {
	var templateFs fs.FS
	//If the user has not specified a template we use the system template.
	if userTemplateFs == nil {
		templateFs = defaultFS
	} else {
		templateFs = userTemplateFs
	}
	return &ProjectGenerator{
		opts:       opts,
		templateFs: templateFs,
		infoLog:    os.Stdout,
		errorLog:   os.Stderr,
	}
}

func (g *ProjectGenerator) Render(targetDir string) error {
	rootFs, err := fs.Sub(g.templateFs, "template")
	if err != nil {
		return fmt.Errorf("Failed to create sub FS: %w", err)
	}

	return fs.WalkDir(rootFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, path)

		// Creating dir
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Copying file
		sourseContent, err := fs.ReadFile(rootFs, path)
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".tmpl") {
			targetPath = strings.TrimSuffix(targetPath, ".tmpl")
			return g.proccessTemplate(targetPath, string(sourseContent), g.opts)
		}

		return os.WriteFile(targetPath, sourseContent, 0644)

	})

}

func (g *ProjectGenerator) proccessTemplate(targetPath string, sourceConetent string, data any) error {
	tmpl, err := template.New(filepath.Base(targetPath)).Parse(sourceConetent)
	if err != nil {
		return fmt.Errorf("Failed to parse template %s: %w", targetPath, err)
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("Failed to create file %s: %w", targetPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		_ = os.Remove(targetPath)
		return fmt.Errorf("Failed to execute template %s: %w", targetPath, err)
	}

	return nil
}
