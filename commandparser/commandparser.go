package commandparser

type Command struct {
	Exe         string
	Args        map[string]string
	Positionals []string
}

type CommandParser struct {
	booleans []string
}

func (parser *CommandParser) AddBoolean(key ...string) *CommandParser {
	parser.booleans = append(parser.booleans, key...)
	return parser
}

func (parser *CommandParser) isBoolean(key string) bool {
	for _, boolKey := range parser.booleans {
		if boolKey == key {
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

			if parser.isBoolean(key) { // don't progress if key is boolean
				arguments[key] = "1"
			} else {
				arguments[key] = argv[i+1]
				i++
			}
		} else { // positional
			positionals = append(positionals, argv[i])
		}
	}

	return Command{
		Exe:         argv[0],
		Args:        arguments,
		Positionals: positionals,
	}
}
