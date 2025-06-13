package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var templateFS embed.FS

func Execute() {
	args := os.Args

	if len(args) < 3 || args[1] != "new" {
		fmt.Println("Usage: go-starter new <project-name>")
		return
	}

	projectName := args[2]

	if projectName == "" {
		fmt.Println("Error: Project name cannot be empty")
		return
	}

	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' already exists\n", projectName)
		return
	}

	err := copyTemplate(projectName)
	if err != nil {
		fmt.Println("Error creating project:", err)
		return
	}

	fmt.Printf("✅ Project '%s' created successfully!\n", projectName)
}

func copyTemplate(projectName string) error {
	templateDir := "template"

	return fs.WalkDir(templateFS, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == templateDir {
			return nil
		}

		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(projectName, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		content, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		modifiedContent := strings.ReplaceAll(string(content), "go-starter-template", projectName)

		return os.WriteFile(targetPath, []byte(modifiedContent), 0644)
	})
}
