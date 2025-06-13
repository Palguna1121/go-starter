// cmd/root.go
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Execute() {
	args := os.Args

	if len(args) < 3 || args[1] != "new" {
		fmt.Println("Usage: go-starter new <project-name>")
		return
	}

	projectName := args[2]
	err := copyTemplate(projectName)
	if err != nil {
		fmt.Println("Error creating project:", err)
		return
	}

	fmt.Printf("✅ Project '%s' created successfully!\n", projectName)
}

func copyTemplate(projectName string) error {
	src := filepath.Join(".", "template")
	dst := filepath.Join(".", projectName)

	return copyDir(src, dst)
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}

		input, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := strings.ReplaceAll(string(input), "go-starter-template", filepath.Base(dst))

		return os.WriteFile(targetPath, []byte(content), info.Mode())
	})
}
