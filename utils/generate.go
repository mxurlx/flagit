package utils

import (
	"fmt"
	"os"
)

const flagsHashFile = "flags_hash.txt"

func GenFiles(filePath string, flags map[string]map[string][]string) error {
	lastFlagsHash, err := readHashFromFile()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error: failed to read last hash from file: %w", err)
	}

	currentFlagsHash := calcFlagsHash(flags)

	if currentFlagsHash == lastFlagsHash {
		return nil
	}
	fmt.Println("New hash detected. Regenerating files")
	err = GenCmdMap(filePath, flags)
	if err != nil {
		return fmt.Errorf("error: couldn't generate Cmd map: %w", err)
	}
	err = GenCmdFiles(filePath, flags)
	if err != nil {
		return fmt.Errorf("error: couldn't generate Cmd files: %w", err)
	}

	err = writeHashToFile(currentFlagsHash)
	if err != nil {
		return fmt.Errorf("error: failed to write current hash to file: %w", err)
	}

	return nil
}

func calcFlagsHash(flags map[string]map[string][]string) string {
	return fmt.Sprintf("%v", flags)
}

func readHashFromFile() (string, error) {
	data, err := os.ReadFile(flagsHashFile)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("error: failed to read hash from file: %w", err)
	}
	return string(data), nil
}

func writeHashToFile(hash string) error {
	return os.WriteFile(flagsHashFile, []byte(hash), 0644)
}
