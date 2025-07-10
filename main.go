package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mholt/archiver/v3"
	"github.com/spf13/cobra"
)

const (
	templateRepo = "https://raw.githubusercontent.com/Palguna1121/go-starter/main/template.zip"
)

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "go-starter",
	Short: "Go Starter CLI",
	Long:  "A CLI to scaffold new Go projects from GitHub template.",
}

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new Go project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		projectName := args[0]

		fmt.Println("🚀 Welcome to Go Starter!")
		fmt.Printf("🛠  Creating project: %s\n\n", projectName)

		// Validate and create project directory first
		if err := os.Mkdir(projectName, 0755); err != nil {
			fmt.Printf("❌ Failed to create project directory: %v\n", err)
			os.Exit(1)
		}

		// Step 1: Download template
		fmt.Println("🔽 Downloading template...")
		err := downloadTemplate()
		if err != nil {
			fmt.Printf("❌ Download failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Template downloaded")

		// Step 2: Extract template
		fmt.Println("\n📦 Extracting template...")
		extractedDir, err := extractTemplate()
		if err != nil {
			fmt.Printf("❌ Extraction failed: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(extractedDir)

		// Find template directory
		templateBaseDir := findTemplateDir(extractedDir)
		if templateBaseDir == "" {
			fmt.Println("❌ Template folder not found in the extracted files")
			os.Exit(1)
		}
		fmt.Printf("✅ Template found at: %s\n", templateBaseDir)

		// Step 3: Process template files
		fmt.Println("\n🔄 Processing files...")
		err = processFiles(templateBaseDir, projectName)
		if err != nil {
			fmt.Printf("❌ Processing failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Files processed")

		// Step 4: Initialize go module
		fmt.Println("\n📦 Initializing Go module...")
		cmdInit := exec.Command("go", "mod", "init", projectName)
		cmdInit.Dir = projectName
		cmdInit.Stdout = os.Stdout
		cmdInit.Stderr = os.Stderr
		if err := cmdInit.Run(); err != nil {
			fmt.Printf("❌ Failed to initialize go module: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Go module initialized")

		// Step 5: Clean up
		fmt.Println("\n🧹 Cleaning up...")
		os.Remove("template.zip")

		// Final output
		fmt.Printf("\n🎉 Project created successfully in %.2f seconds!\n", time.Since(startTime).Seconds())
		fmt.Println("🚀 You're ready to go!")

		fmt.Println("\n👉 Next steps:")
		fmt.Printf("   1. Change into the project directory:\n      cd %s\n", projectName)
		fmt.Println("   2. Install dependencies and set up environment:")
		fmt.Println("      go mod tidy")
		fmt.Println("      cp .env.example .env")
		fmt.Println("   3. Run the app easily with Makefile:")
		fmt.Println("      make install   # this command will do `go mod tidy` and `cp .env.example .env`")
		fmt.Println("      make run       # this command will Starts the application")

		fmt.Println("\n💡 Tip: You can edit .env for your local config.")
		fmt.Println("\nHappy coding! 💻")
	},
}

func downloadTemplate() error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(templateRepo)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %s", resp.Status)
	}

	out, err := os.Create("template.zip")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}

	return nil
}

func extractTemplate() (string, error) {
	tempDir, err := os.MkdirTemp("", "go-starter-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	err = archiver.Unarchive("template.zip", tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to unzip: %w", err)
	}

	return tempDir, nil
}

func findTemplateDir(baseDir string) string {
	// Cari folder template di beberapa lokasi yang mungkin
	possiblePaths := []string{
		filepath.Join(baseDir, "go-starter-main", "template"),
		filepath.Join(baseDir, "template"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func processFiles(sourceDir, projectName string) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(projectName, relPath)

		// Skip root directory
		if path == sourceDir {
			return nil
		}

		// Create parent directory first if needed
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
		}

		// Skip directories (already created by MkdirAll)
		if info.IsDir() {
			fmt.Printf("📁 %s\n", relPath)
			return nil
		}

		// Process files
		fmt.Printf("📝 %s\n", relPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Process text files
		newContent := strings.ReplaceAll(string(content), "response-std", projectName)

		return os.WriteFile(targetPath, []byte(newContent), 0644)
	})
}

func Execute() {
	rootCmd.AddCommand(newCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}
