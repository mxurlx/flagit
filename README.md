# Flagit - Flexible Go CLI Flag Parsing

## Overview

Flagit is a Go library designed to provide a more simple and manageable way to parse command-line flags for your CLI applications. It addresses limitations in the standard `flags` package by allowing developers to declare all flags in one central location, support both short and long flag names, handle mandatory arguments, and manage subcommands with ease.

## Features

*   **Centralized Flag Declaration:** Define all flags in a single map for clear organization.
*   **Short & Long Flags:** Supports both short (e.g., `-s`) and full (e.g., `--shell`) flag names, as well as combined short flags (e.g. `-rf`)
*   **Mandatory Arguments:** Easily define arguments that are required without needing a flag.
*   **Subcommand Support:**  Handles commands like `add`, `remove`, etc.
*   **Module-Aware Embedding:** Embeds the library directly into your executable for easy distribution – no external dependencies

## Installation

```bash
go get github.com/mxurlx/flagit
```

## Usage

Here's a basic example of how to integrate Flagit into your Go CLI application:

```go
package main

import (
	"flagit"
	"fmt"
	"gotestforflagparser/cmd"
	"gotestforflagparser/common"
	"os"
)

func main() {
	// 1) Initialize the Flags map.
	// Generates and populates the `Flags` map file with the template.
	// Must be done once or when the `Flags` file has been deleted.
	// Remove it if `Flags` map file is present in your project
	err := flagit.InitFlagsMap()
	if err != nil {
		fmt.Println(err)
	    return
	}

	// 2) Generate command files based on the Flags map.
	err = flagit.GenFiles(common.Flags)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3) Parse the command-line flags.
	subcmd, flags, mandatoryArgs, err := flagit.ParseFlags(os.Args, common.Flags)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Execute the appropriate command with parsed flags and arguments.
	err = flagit.ExecuteCmd(subcmd, flags, mandatoryArgs, cmd.CmdFuncs)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

**Explanation:**

1.  **`flagit.InitFlagsMap()`**: This function generates and populates the `Flags` map file in the `/common` directory of your project with a template based on your defined flag structure. It needs to be called once, or whenever you delete the `Flags` map file, after which this function must be removed from your code.
2.  **`flagit.GenFiles(common.Flags)`**: Generates command files based on the declared flags in `common.Flags`. This function creates the necessary structures for parsing.
3.  **`flagit.ParseFlags(os.Args, common.Flags)`**: Parses the command-line arguments (`os.Args`) using the defined flags in `common.Flags`, returning the subcommand, a map of parsed flags, and any mandatory arguments.
4.  **`flagit.ExecuteCmd(subcmd, flags, mandatoryArgs, cmd.CmdFuncs)`**: Executes the appropriate command (determined by `subcmd`) with the parsed flags and arguments.

## Flag Declaration

Flags are declared in a map within your project (e.g., `common/flags.go`):

```go
package common

var Flags = map[string]map[string][]string{
    ".": {  // Flags applicable when no subcommand is specified.
        "debug": {"D", "false", "Show debugging process"}
    }
    "add": {
        "username": {"u", "<mandatory>", "Username to create"},
        "homedir": {"d", "", "Home directory of the user"},
        "shell":   {"s", "", "Login shell of the user"},
        // ... more flags ...
    },
}
```

*   The outer key is the **subcommand** (e.g., `"add"`). Use `"."` to indicate flags that apply when no subcommand is provided.
*   Each inner map represents the flags for that subcommand.
*   Each flag entry has a slice of strings: `{"short_flag", "default_value| <mandatory>", "description"}`.
