package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mxurlx/flagit"
)

func TestParseFlagsAllowsArgumentlessSubcommand(t *testing.T) {
	flagsMap := map[string]map[string][]string{
		".": {
			"help": {"h", "false", "Show help message"},
		},
		"init": {},
	}

	subcmd, flags, mandatoryArgs, err := flagit.ParseFlags([]string{"app", "init"}, flagsMap)
	if err != nil {
		t.Fatalf("ParseFlags returned an error: %v", err)
	}
	if subcmd != "init" {
		t.Fatalf("subcmd = %q, want %q", subcmd, "init")
	}
	if len(flags) != 0 {
		t.Fatalf("flags = %v, want no flags", flags)
	}
	if len(mandatoryArgs) != 0 {
		t.Fatalf("mandatoryArgs = %v, want no mandatory args", mandatoryArgs)
	}
}

func TestParseFlagsAllowsEmptyArguments(t *testing.T) {
	subcmd, flags, mandatoryArgs, err := flagit.ParseFlags(nil, map[string]map[string][]string{".": {}})
	if err != nil {
		t.Fatalf("ParseFlags returned an error: %v", err)
	}
	if subcmd != "." {
		t.Fatalf("subcmd = %q, want root command", subcmd)
	}
	if help, ok := flags["help"].(bool); !ok || !help {
		t.Fatalf("help flag = %v, want true", flags["help"])
	}
	if len(mandatoryArgs) != 0 {
		t.Fatalf("mandatoryArgs = %v, want no mandatory args", mandatoryArgs)
	}
}

func TestParseFlagsReturnsErrorForUnknownFlagOnArgumentlessSubcommand(t *testing.T) {
	flagsMap := map[string]map[string][]string{
		".":    {},
		"init": {},
	}

	_, _, _, err := flagit.ParseFlags([]string{"app", "init", "--unknown"}, flagsMap)
	if err == nil {
		t.Fatal("ParseFlags returned nil error, want unknown flag error")
	}
}

func TestGenFilesDoesNotGenerateWhenDisabled(t *testing.T) {
	t.Setenv("FLAGIT_GENERATE", "0")
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	flagsMap := map[string]map[string][]string{
		".": {
			"help": {"h", "false", "Show help message"},
		},
		"init": {},
	}

	if err := flagit.GenFiles(flagsMap); err != nil {
		t.Fatalf("GenFiles returned an error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "flags_hash.txt")); !os.IsNotExist(err) {
		t.Fatalf("flags_hash.txt stat error = %v, want not exist", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "cmd")); !os.IsNotExist(err) {
		t.Fatalf("cmd stat error = %v, want not exist", err)
	}

	if err := flagit.InitFlagsMap(); err != nil {
		t.Fatalf("InitFlagsMap returned an error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "common")); !os.IsNotExist(err) {
		t.Fatalf("common stat error = %v, want not exist", err)
	}
}

func TestGenFilesCanBeForcedForDevelopmentGeneration(t *testing.T) {
	t.Setenv("FLAGIT_GENERATE", "1")
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	flagsMap := map[string]map[string][]string{
		".": {
			"help": {"h", "false", "Show help message"},
		},
		"init": {},
	}

	if err := flagit.GenFiles(flagsMap); err != nil {
		t.Fatalf("GenFiles returned an error: %v", err)
	}

	for _, path := range []string{
		filepath.Join(tmpDir, "flags_hash.txt"),
		filepath.Join(tmpDir, "cmd", "cmdfuncs.go"),
		filepath.Join(tmpDir, "cmd", "init.go"),
		filepath.Join(tmpDir, "cmd", "root.go"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected generated file %s: %v", path, err)
		}
	}
}
