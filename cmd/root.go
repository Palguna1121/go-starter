package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Execute() {
	if len(os.Args) != 3 || os.Args[1] != "new" {
		log.Fatal("❌ Usage: go-starter new <project-name>")
	}

	projectName := os.Args[2]
	src := filepath.Join(getCurrentDir(), "..", "..", "template")
	dst := filepath.Join(getCurrentDir(), projectName)

	err := copyTemplate(src, dst, "response-std", projectName)
	if err != nil {
		log.Fatalf("❌ Failed to create project: %v", err)
	}

	fmt.Println("✅ Project", projectName, "created successfully!")
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("❌ Cannot get current dir:", err)
	}
	return dir
}

func copyTemplate(src, dst, oldName, newName string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		relPath = strings.ReplaceAll(relPath, oldName, newName)
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}

		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace placeholders
		text := string(content)
		text = strings.ReplaceAll(text, oldName, newName)
		text = strings.ReplaceAll(text, "{{ProjectName}}", newName)

		return os.WriteFile(targetPath, []byte(text), info.Mode())
	})
}
