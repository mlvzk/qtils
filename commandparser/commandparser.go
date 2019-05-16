package commandparser

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
	options map[string]Option
	aliases map[string]string
}

func New() *CommandParser {
	return &CommandParser{
		map[string]Option{},
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
		parser.options[option.GetKey()] = option
		for _, t := range option.GetAliases() {
			parser.aliases[t] = option.GetKey()
		}
	}
}

func (parser *CommandParser) Parse(argv []string) (*Command, error) {
	if len(argv) == 0 {
		return nil, EmptyArgvError
	}

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

			option, optionFound := parser.options[key]
			if !optionFound {
				return nil, NewInvalidKeyError(key)
			}

			if option.IsBoolean() && option.IsArrayed() {
				arrayed[key] = append(arrayed[key], "1")
			} else if option.IsArrayed() {
				if len(argv) <= i+1 {
					return nil, NewMissingValueError(key)
				}
				arrayed[key] = append(arrayed[key], argv[i+1])
				i++
			} else if option.IsBoolean() {
				booleans[key] = true
			} else {
				if len(argv) <= i+1 {
					return nil, NewMissingValueError(key)
				}
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
