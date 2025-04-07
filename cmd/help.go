package cmd

import (
	"fmt"
	"strconv"
)

func PrintHelp(cmd string, flagsMap map[string]map[string][]string) error {
	if cmd == "." {
		fmt.Println("\nMain Command Help:")
		err := printSubHelp(cmd, flagsMap)
		if err != nil {
			return err
		}
		for subcmd := range flagsMap {
			if subcmd != "." {
				PrintHelp(subcmd, flagsMap)
			}
		}
	} else {
		fmt.Printf("\nHelp for subcommand '%s':\n", cmd)
		err := printSubHelp(cmd, flagsMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func insert(slice []string, value string, index int) []string {
	if index < 0 || index > len(slice) {
		slice = append(slice, value)
		return slice
	}

	slice = append(slice[:index], append([]string{value}, slice[index:]...)...)
	return slice
}

func printSubHelp(cmd string, flagsMap map[string]map[string][]string) error {
	var (
		mandatoryArgs []string
		flags         []string
	)
	for flagName, info := range flagsMap[cmd] {
		shortName := info[0]
		defVal := info[1]
		desc := info[2]

		if defVal == "" {
			defVal = "None"
		}

		if defVal != "<mandatory>" {
			flags = append(flags, fmt.Sprintf("\t--%s (-%s): %s (default: %s)\n", flagName, shortName, desc, defVal))
		} else {
			pos, err := strconv.Atoi(shortName)
			if err != nil {
				return err
			}
			mandatoryArgs = insert(mandatoryArgs, fmt.Sprintf("\t%d) %s (%s) (mandatory argument)\n", pos, flagName, desc), pos-1)
		}
	}
	if len(mandatoryArgs) > 0 {
		fmt.Println()
	}
	for _, row := range mandatoryArgs {
		fmt.Print(row)
	}
	fmt.Println()
	for _, row := range flags {
		fmt.Print(row)
	}

	return nil
}
