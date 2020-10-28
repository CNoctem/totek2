package config

import "errors"

type EntityType string

type Option struct {
	Name        string
	Description string
	nArgs       int
}

func NewCommand(options *[]Option, name, description string) {
	NewOption(options, name, description, 0)
}

func NewOption(options *[]Option, name, description string, nArgs int) {
	newOpt := Option{name, description, nArgs}
	*options = append(*options, newOpt)
}

func ParseArgs(defined []Option, args []string) (*map[string][]string, error) {
	if len(defined) == 0 {
		return nil, errors.New("no commands or options defined")
	}

	entMap := make(map[string][]string, 0)
	for i := 0; i < len(args); i++ {
		a := args[i]
		for _, e := range defined {
			if e.Name == a {
				eArgs := make([]string, e.nArgs)
				for j := 0; j < e.nArgs; j++ {
					i++
					eArgs[j] = args[i]
				}
				entMap[e.Name] = eArgs
			}
		}
	}

	return &entMap, nil
}

func GetDescription(options []Option, cmdName string) string {
	for _, option := range options {
		if option.Name == cmdName {
			return option.Description
		}
	}
	return "Undefined command or option: "+ cmdName
}



