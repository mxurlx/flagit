package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenCmdFiles(filePath string, flags map[string]map[string][]string) error {

	err := os.MkdirAll(filepath.Join(filePath, "cmd"), 0755)
	if err != nil {
		return fmt.Errorf("error: failed to create generated directory: %w", err)
	}
	for key := range flags {
		cmd := key
		if key == "." {
			cmd = "root"
		}
		outputPath := filepath.Join("cmd", fmt.Sprintf("%s.go", cmd))
		if _, err := os.Stat(outputPath); err == nil {
			continue
		}
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}

		fileName := cases.Title(language.English, cases.Compact).String(cmd)

		file.WriteString("package cmd\n\n")
		file.WriteString("import \"fmt\"\n\n")
		file.WriteString(fmt.Sprintf("func %s(flags map[string]any, mandatoryArgs []string) error {\n", fileName))
		file.WriteString(fmt.Sprintf(
			`	fmt.Printf("Command '%s' executed\n")`, cmd))
		file.WriteString("\n\t// Implement your command logic here\n")
		file.WriteString("\t return nil\n")
		file.WriteString("}\n")
		file.Close()
		fmt.Printf("cmd/%s.go generated successfully\n", cmd)
	}

	return nil
}
