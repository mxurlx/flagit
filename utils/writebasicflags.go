package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitFlagsMap() error {
	filePath, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(filePath, "common"), 0755)
	if err != nil {
		return fmt.Errorf("error: failed to create flags directory: %w", err)
	}
	outputPath := filepath.Join("common", "flags.go")

	file, err := os.Create(outputPath)

	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(`package common

var Flags = map[string]map[string][]string{
	".": {
		"help":    {"h", "false", "Show help message"},
		"version": {"v", "false", "Show version"},
	},
	"subcommand": {
		"username": {"1", "<mandatory>", "Username to create"},
		"another":  {"2", "<mandatory>", "Another username to create"},
		"homedir":  {"d", "", "Home directory of the user"},
		"shell":    {"s", "", "Login shell of the user"},
		"group":    {"g", "", "Primary group of the user"},
		"groups":   {"G", "", "Supplementary groups of the user (comma-separated)"},
		"debug":    {"D", "false", "Show debugging process"},
		"number":   {"n", "13134", "Number"},
	},
}
`)

	return nil
}
