package commandparser

import "errors"

type Command struct {
	Exe         string
	Args        map[string]string
	Arrayed     map[string][]string
	Booleans    map[string]bool
	Positionals []string
}

func (c Command) String() string {
	res := c.Exe

	for key, value := range c.Args {
		res += " --" + key + " '" + value + "'"
	}

	for key, arrayedValue := range c.Arrayed {
		for _, value := range arrayedValue {
			res += " --" + key + " '" + value + "'"
		}
	}

	for key := range c.Booleans {
		res += " --" + key
	}

	for _, positional := range c.Positionals {
		res += " " + positional
	}

	return res
}

type CommandParser struct {
	keys     map[string]struct{}
	booleans []string
	arrayed  []string
	aliases  map[string]string
}

func New() *CommandParser {
	return &CommandParser{
		map[string]struct{}{},
		[]string{},
		[]string{},
		map[string]string{},
	}
}

type Option interface {
	GetKey() string
	GetAliases() []string
	IsBoolean() bool
	IsArrayed() bool
}

func (parser *CommandParser) AddOption(options ...Option) {
	for _, option := range options {
		parser.keys[option.GetKey()] = struct{}{}
		parser.addAliases(option.GetKey(), option.GetAliases()...)
		if option.IsBoolean() {
			parser.addBoolean(option.GetKey())
		}
		if option.IsArrayed() {
			parser.addArrayed(option.GetKey())
		}
	}
}

func (parser *CommandParser) addBoolean(key ...string) {
	parser.booleans = append(parser.booleans, key...)
}

func (parser *CommandParser) addArrayed(key ...string) {
	parser.arrayed = append(parser.arrayed, key...)
}

func (parser *CommandParser) addAliases(from string, to ...string) {
	for _, t := range to {
		parser.aliases[t] = from
	}
}

func (parser *CommandParser) isBoolean(key string) bool {
	for _, boolKey := range parser.booleans {
		alias, _ := parser.aliases[key]
		if boolKey == key || boolKey == alias {
			return true
		}
	}

	return false
}

func (parser *CommandParser) isArrayed(key string) bool {
	for _, arrayKey := range parser.arrayed {
		if arrayKey == key {
			return true
		}
	}

	return false
}

func (parser *CommandParser) Parse(argv []string) (*Command, error) {
	var (
		positionals []string
		arguments   = map[string]string{}
		arrayed     = map[string][]string{}
		booleans    = map[string]bool{}
	)

	for i := 1; i < len(argv); i++ {
		if argv[i][0] == '-' && len(argv[i]) > 1 {
			key := ""
			if argv[i][1] == '-' { // --long_key_param
				key = argv[i][2:]
			} else { // -short
				key = argv[i][1:]
			}

			alias, hasAlias := parser.aliases[key]
			if hasAlias { // convert aliased key to original key
				key = alias
			}

			if _, found := parser.keys[key]; !found {
				return nil, errors.New("Invalid key '" + key + "'")
			}

			if parser.isBoolean(key) && parser.isArrayed(key) {
				arrayed[key] = append(arrayed[key], "1")
			} else if parser.isArrayed(key) {
				arrayed[key] = append(arrayed[key], argv[i+1])
				i++
			} else if parser.isBoolean(key) {
				booleans[key] = true
			} else {
				arguments[key] = argv[i+1]
				i++
			}
		} else {
			positionals = append(positionals, argv[i])
		}
	}

	return &Command{
		Exe:         argv[0],
		Args:        arguments,
		Arrayed:     arrayed,
		Booleans:    booleans,
		Positionals: positionals,
	}, nil
}
