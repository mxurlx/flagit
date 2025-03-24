package flagit

import (
	"flagit/cmd"
	"flagit/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func InitFlagsMap() error {
	return utils.InitFlagsMap()
}

func PrintHelp(command string, flagsMap map[string]map[string][]string) {
	cmd.PrintHelp(command, flagsMap)
}

func GenFiles(flags map[string]map[string][]string) error {
	filePath, err := os.Getwd()
	if err != nil {
		return err
	}

	return utils.GenFiles(filePath, flags)
}

func ExecuteCmd(subcmd string, flags map[string]any, mandatoryArgs []string, cmdFuncs map[string]func(flags map[string]any, mandatoryArgs []string) error) error {
	command, ok := cmdFuncs[subcmd]
	if !ok {
		return fmt.Errorf("unknown subcommand '%s'", subcmd)
	}

	err := command(flags, mandatoryArgs)

	return err
}

func ParseFlags(arguments []string, flagsMap map[string]map[string][]string) (string, map[string]any, []string, error) {
	args := arguments[1:]
	subcmd := "."
	flags := make(map[string]any)
	mandatoryArgs := []string{}
	mandatoryArgsCheck := 0
	i := 1

	if len(args) > 0 {
		subcmd = args[0]
	} else {
		flags["help"] = true
		return subcmd, flags, mandatoryArgs, nil
	}

	_, ok := flagsMap[subcmd]
	if !ok || strings.HasPrefix(subcmd, "-") {
		subcmd = "."
		i = 0
	}

	for flagName, info := range flagsMap[subcmd] {
		populateFlags(flags, subcmd, flagName, info[1], flagsMap, true)
		if info[1] == "<mandatory>" {
			mandatoryArgsCheck++
		}
	}

	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flagName := ""
			value := ""
			flagType := ""

			if strings.HasPrefix(arg, "--") {
				flagName = arg[2:]
			} else {
				shorts := splitShortFlags(args[i])
				for _, flag := range shorts {
					flagName = flag
					tmp, err := findFullbyShort(flagName, subcmd, flagsMap)
					if err != nil {
						fmt.Printf("error: short flag %s not found\n", flagName)
						i++
						continue
					}
					flagName = tmp
					populateFlags(flags, subcmd, flagName, value, flagsMap, false)
				}
			}
			flagType = getDataType(flagsMap[subcmd][flagName][1])
			i++
			if i < len(args) && !strings.HasPrefix(args[i], "-") && flagType != "bool" {
				value = args[i]
				i++
			}
			populateFlags(flags, subcmd, flagName, value, flagsMap, false)
		} else {
			mandatoryArgs = append(mandatoryArgs, arg)
			i++
		}
	}

	if mandatoryArgsCheck != len(mandatoryArgs) {
		return subcmd, flags, mandatoryArgs, fmt.Errorf("error: mandatory arguments do not match")
	}

	return subcmd, flags, mandatoryArgs, nil
}

func populateFlags(flags map[string]any, subcmd, flagName, value string, flagsMap map[string]map[string][]string, init bool) {
	flagInfo, ok := flagsMap[subcmd][flagName]
	if ok {
		dataType := getDataType(flagInfo[1])
		switch dataType {
		case "bool":
			tmp, _ := strconv.ParseBool(flagInfo[1])
			if init {
				flags[flagName] = tmp
			} else {
				flags[flagName] = !tmp
			}
		case "int":
			intVal, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("error: expected int value for %s, got '%s'\n", flagName, value)
				break
			}
			flags[flagName] = intVal
		default:
			flags[flagName] = value
		}
	}
}

func findFullbyShort(shortFlag string, subcmd string, flagsMap map[string]map[string][]string) (string, error) {
	flagInfo := flagsMap[subcmd]
	for fullName, info := range flagInfo {
		if info[0] == shortFlag {
			return fullName, nil
		}
	}

	return "", fmt.Errorf("error: no full flag name found for short flag '%s'", shortFlag)
}

func isInt(val string) bool {
	_, err := strconv.Atoi(val)
	return err == nil
}

func getDataType(defVal string) string {
	if defVal == "true" || defVal == "false" {
		return "bool"
	}
	if isInt(defVal) {
		return "int"
	}
	return "string"
}

func splitShortFlags(str string) []string {
	var args []string

	for i := 1; i < len(str); i++ {
		args = append(args, string(str[i]))
	}

	return args
}
