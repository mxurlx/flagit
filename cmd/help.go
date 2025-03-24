package cmd

import (
	"fmt"
)

func PrintHelp(cmd string, flagsMap map[string]map[string][]string) {
	if cmd == "." {
		fmt.Println("\nMain Command Help:")
		printSubHelp(cmd, flagsMap)
		for subcmd := range flagsMap {
			if subcmd != "." {
				PrintHelp(subcmd, flagsMap)
			}
		}
	} else {
		fmt.Printf("\nHelp for subcommand '%s':\n", cmd)
		printSubHelp(cmd, flagsMap)
	}
}

func printSubHelp(cmd string, flagsMap map[string]map[string][]string) {
	for flagName, info := range flagsMap[cmd] {
		shortName := info[0]
		defVal := info[1]
		desc := info[2]

		if defVal == "" {
			defVal = "None"
		}

		fmt.Printf("\t--%s (-%s): %s (default: %s)\n", flagName, shortName, desc, defVal)
	}
}
