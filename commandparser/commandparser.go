package commandparser

type Command struct {
	Exe         string
	Args        map[string]string
	Positionals []string
}

type CommandParser struct {
	booleans []string
	aliases  map[string]string
}

func New() CommandParser {
	return CommandParser{
		[]string{},
		map[string]string{},
	}
}

func (parser *CommandParser) AddBoolean(key ...string) {
	parser.booleans = append(parser.booleans, key...)
}

func (parser *CommandParser) AddAliases(from string, to ...string) {
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

func (parser CommandParser) Parse(argv []string) Command {
	arguments := map[string]string{}
	positionals := []string{}

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

			if parser.isBoolean(key) {
				arguments[key] = "1"
			} else {
				arguments[key] = argv[i+1]
				i++
			}
		} else {
			positionals = append(positionals, argv[i])
		}
	}

	return Command{
		Exe:         argv[0],
		Args:        arguments,
		Positionals: positionals,
	}
}
