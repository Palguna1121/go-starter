package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed template/**/*
var TemplateFS embed.FS

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "go-starter",
	Short: "Go Starter CLI",
	Long:  "A CLI to scaffold new Go projects from a template.",
}

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new Go project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		fmt.Println("Welcome to Go Starter!")
		fmt.Printf("Creating new project '%s'...\n", projectName)

		templateRoot := "template"
		err := fs.WalkDir(TemplateFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing path %s: %w", path, err)
			}

			relPath, err := filepath.Rel(templateRoot, path)
			if err != nil {
				return fmt.Errorf("error getting relative path: %w", err)
			}

			targetPath := filepath.Join(projectName, relPath)

			if d.IsDir() {
				return os.MkdirAll(targetPath, 0755)
			}

			data, err := TemplateFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}

			// Skip binary files
			if isBinary(data) {
				return os.WriteFile(targetPath, data, 0644)
			}

			content := string(data)

			// Replace all occurrences of template placeholders
			replacements := map[string]string{
				"response-std": projectName,
				// Tambahkan placeholder lain jika diperlukan
			}

			for old, new := range replacements {
				content = strings.ReplaceAll(content, old, new)
			}

			// Untuk file go.mod, kita perlu replace module path secara khusus
			// if strings.HasSuffix(path, "go.mod") {
			// 	content = fmt.Sprintf("module %s\n\n%s", projectName,
			// 		strings.SplitN(content, "\n", 2)[1])
			// }

			return os.WriteFile(targetPath, []byte(content), 0644)
		})

		if err != nil {
			fmt.Printf("❌ Error creating project: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Println("✅ Project created successfully!")
		}
		fmt.Println("Thank you for using Go Starter!")
		fmt.Println("Happy coding!!! 🚀")
	},
}

func isBinary(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}

func Execute() {
	rootCmd.AddCommand(newCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}
