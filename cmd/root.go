package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-starter",
	Short: "Go Starter CLI",
}

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new Go starter project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		templatePath, err := getTemplatePath()
		if err != nil {
			return err
		}

		outputPath, err := filepath.Abs(projectName)
		if err != nil {
			return err
		}

		fmt.Printf("Creating new project '%s'...\n", projectName)

		// Ganti "response-std" jadi nama baru
		err = copyDir(templatePath, outputPath, map[string]string{
			"response-std":    projectName,
			"__MODULE_NAME__": projectName,
		})
		if err != nil {
			return err
		}

		fmt.Println("✅ Project created successfully at", outputPath)
		return nil
	},
}

func Execute() {
	rootCmd.AddCommand(newCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}

func getTemplatePath() (string, error) {
	// Coba cari berdasarkan relative path dari source
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	localTemplate := filepath.Join(cwd, "template")
	if _, err := os.Stat(localTemplate); err == nil {
		return localTemplate, nil
	}

	// Kalau tidak ditemukan, fallback ke path relatif dari executable
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	execDir := filepath.Dir(execPath)
	parentDir := filepath.Dir(execDir)
	embeddedTemplate := filepath.Join(parentDir, "template")

	if _, err := os.Stat(embeddedTemplate); err == nil {
		return embeddedTemplate, nil
	}

	return "", fmt.Errorf("template folder not found")
}

func copyDir(src string, dst string, replacements map[string]string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
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

		content := string(input)
		for old, new := range replacements {
			content = strings.ReplaceAll(content, old, new)
		}

		return os.WriteFile(targetPath, []byte(content), info.Mode())
	})
}
