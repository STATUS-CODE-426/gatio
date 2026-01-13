package generator

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//NOTE: I tested this code and it works fine but i think it needs refactoring later

//go:embed template/**/*
var templateEmbed embed.FS
var templateFs fs.FS = templateEmbed

func RenderProject(targetDir string, templateFs fs.FS, data any) error {
	return fs.WalkDir(templateFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetDir, path)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		content, err := fs.ReadFile(templateFs, path)
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".tmp") {
			targetPath = strings.TrimSuffix(targetPath, ".tmp")

			tmpl, err := template.New(path).Parse(string(content))
			if err != nil {
				return err
			}

			f, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer f.Close()

			return tmpl.Execute(f, data)
		}

		return os.WriteFile(targetPath, content, 0644)
	})
}
