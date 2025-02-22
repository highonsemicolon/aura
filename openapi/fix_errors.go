package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var errorPattern = regexp.MustCompile(`(errors\.New|fmt\.Errorf)\(\s*"([^"]+)"`)

func processFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	updatedContent := errorPattern.ReplaceAllStringFunc(string(content), func(match string) string {
		matches := errorPattern.FindStringSubmatch(match)
		if len(matches) == 3 {
			errorFunc, errorMsg := matches[1], matches[2]
			if len(errorMsg) > 0 && strings.ToUpper(errorMsg[:1]) == errorMsg[:1] {
				errorMsg = strings.ToLower(errorMsg[:1]) + errorMsg[1:]
				return fmt.Sprintf(`%s("%s"`, errorFunc, errorMsg)
			}
		}
		return match
	})

	if string(content) != updatedContent {
		return os.WriteFile(path, []byte(updatedContent), 0600)
	}
	return nil
}

func processDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			log.Println("Processing", path)
			return processFile(path)
		}
		return nil
	})
}

func main() {
	dir := ".src/api"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	if err := processDirectory(dir); err != nil {
		fmt.Println("Error processing files:", err)
		os.Exit(1)
	}

	fmt.Println("All error messages have been fixed!")
}
