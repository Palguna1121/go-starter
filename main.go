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

//go:embed template/*
var templateFS embed.FS

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

		// Verifikasi template ada
		if _, err := templateFS.ReadDir("template"); err != nil {
			fmt.Println("❌ Error: Failed to read template files")
			fmt.Println("Make sure the 'template' folder exists and contains files")
			fmt.Println("Original error:", err)
			os.Exit(1)
		}

		err := fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("access error: %w", err)
			}

			relPath := strings.TrimPrefix(path, "template/")
			targetPath := filepath.Join(projectName, relPath)

			if d.IsDir() {
				return os.MkdirAll(targetPath, 0755)
			}

			data, err := templateFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}

			content := string(data)
			content = strings.ReplaceAll(content, "response-std", projectName)

			return os.WriteFile(targetPath, []byte(content), 0644)
		})

		if err != nil {
			fmt.Printf("❌ Error creating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Project created successfully!")
		fmt.Println("Thank you for using Go Starter!")
		fmt.Println("Happy coding!!! 🚀")
	},
}

func Execute() {
	rootCmd.AddCommand(newCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}
