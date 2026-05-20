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

func TestExecuteCmdAllowsArgumentlessSubcommandWithoutHandler(t *testing.T) {
	if err := flagit.ExecuteCmd("init", map[string]any{}, nil, map[string]func(map[string]any, []string) error{}); err != nil {
		t.Fatalf("ExecuteCmd returned an error: %v", err)
	}
}

func TestExecuteCmdRejectsUnknownSubcommandWithArguments(t *testing.T) {
	err := flagit.ExecuteCmd("create", map[string]any{}, []string{"name"}, map[string]func(map[string]any, []string) error{})
	if err == nil {
		t.Fatal("ExecuteCmd returned nil error, want unknown subcommand error")
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

func TestGenFilesDoesNotGenerateOutsideSourceWorkspaceByDefault(t *testing.T) {
	t.Setenv("FLAGIT_GENERATE", "")
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
}

func TestGenFilesGeneratesInSourceWorkspaceByDefault(t *testing.T) {
	t.Setenv("FLAGIT_GENERATE", "")
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/app\n"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "common"), 0755); err != nil {
		t.Fatalf("failed to create common dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "common", "flags.go"), []byte("package common\n"), 0644); err != nil {
		t.Fatalf("failed to create flags.go: %v", err)
	}

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
