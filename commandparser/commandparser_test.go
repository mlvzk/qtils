package commandparser_test

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mlvzk/qtils/commandparser"
)

type option struct {
	key     string
	aliases []string
	arrayed bool
	boolean bool
}

func (o option) GetKey() string {
	return o.key
}

func (o option) GetAliases() []string {
	return o.aliases
}

func (o option) IsArrayed() bool {
	return o.arrayed
}

func (o option) IsBoolean() bool {
	return o.boolean
}

var _ commandparser.Option = option{}

func TestParseCommand(t *testing.T) {
	cases := []struct {
		name        string
		argv        []string
		options     []option
		expected    *commandparser.Command
		expectedErr error
	}{
		{
			"with invalid key",
			[]string{"./main", "--invalid"},
			[]option{},
			nil,
			commandparser.NewInvalidKeyError("invalid"),
		},
		{
			"with long key",
			[]string{"./main", "--key", "value"},
			[]option{option{key: "key"}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"key": "value"},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with two long keys",
			[]string{"./main", "--key", "value", "--other", "another"},
			[]option{option{key: "key"}, option{key: "other"}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"key": "value", "other": "another"},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with short key",
			[]string{"./main", "-key", "value"},
			[]option{option{key: "key"}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"key": "value"},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with two short keys",
			[]string{"./main", "-key", "value", "-other", "another"},
			[]option{option{key: "key"}, option{key: "other"}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"key": "value", "other": "another"},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with positional",
			[]string{"./main", "value"},
			[]option{},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{},
				Positionals: []string{"value"},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with two positionals",
			[]string{"./main", "value", "another"},
			[]option{},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{},
				Positionals: []string{"value", "another"},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with long key and positional",
			[]string{"./main", "--key", "value", "positional"},
			[]option{option{key: "key"}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"key": "value"},
				Positionals: []string{"positional"},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with short key and positional",
			[]string{"./main", "-key", "value", "positional"},
			[]option{option{key: "key"}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"key": "value"},
				Positionals: []string{"positional"},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with boolean key",
			[]string{"./main", "-boolean"},
			[]option{option{key: "boolean", boolean: true}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{"boolean": true},
			},
			nil,
		},
		{
			"with multiple boolean keys",
			[]string{"./main", "-boolean", "-boolean2"},
			[]option{option{key: "boolean", boolean: true}, option{key: "boolean2", boolean: true}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{"boolean": true, "boolean2": true},
			},
			nil,
		},
		{
			"with dash positional",
			[]string{"./main", "-"},
			[]option{},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{},
				Positionals: []string{"-"},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with short key alias",
			[]string{"./main", "-o", "file"},
			[]option{option{key: "output", aliases: []string{"o"}}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"output": "file"},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with long key alias",
			[]string{"./main", "--o", "file"},
			[]option{option{key: "output", aliases: []string{"o"}}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"output": "file"},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with multiple aliases to same key",
			[]string{"./main", "--out", "file"},
			[]option{option{key: "output", aliases: []string{"o", "out"}}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"output": "file"},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"sed -n --quiet --silent",
			[]string{"sed", "-n", "--quiet", "--silent"},
			[]option{option{key: "quiet", boolean: true, aliases: []string{"n", "silent"}}},
			&commandparser.Command{
				Exe:         "sed",
				Args:        map[string]string{},
				Positionals: []string{},
				Arrayed:     map[string][]string{},
				Booleans:    map[string]bool{"quiet": true},
			},
			nil,
		},
		{
			"with arrayed key",
			[]string{"./main", "--key", "value1", "--key", "value2"},
			[]option{option{key: "key", arrayed: true}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{},
				Positionals: []string{},
				Arrayed:     map[string][]string{"key": []string{"value1", "value2"}},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"with same key both boolean and arrayed",
			[]string{"./main", "--verbose", "--verbose"},
			[]option{option{key: "verbose", boolean: true, arrayed: true}},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{},
				Positionals: []string{},
				Arrayed:     map[string][]string{"verbose": []string{"1", "1"}},
				Booleans:    map[string]bool{},
			},
			nil,
		},
		{
			"mixed",
			[]string{"./main", "first_pos", "-v", "-v", "-v", "--key", "value", "-other", "another", "-boolean", "-b2", "-arrayed_key", "arrayed1", "-arrayed_key", "arrayed2", "last_pos"},
			[]option{
				option{key: "key"},
				option{key: "other"},
				option{key: "arrayed_key", arrayed: true},
				option{key: "boolean", boolean: true},
				option{key: "boolean2", boolean: true, aliases: []string{"b2"}},
				option{key: "verbose", boolean: true, arrayed: true, aliases: []string{"v"}},
			},
			&commandparser.Command{
				Exe:         "./main",
				Args:        map[string]string{"key": "value", "other": "another"},
				Positionals: []string{"first_pos", "last_pos"},
				Arrayed:     map[string][]string{"arrayed_key": {"arrayed1", "arrayed2"}, "verbose": {"1", "1", "1"}},
				Booleans:    map[string]bool{"boolean": true, "boolean2": true},
			},
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := commandparser.New()

			for _, option := range testCase.options {
				parser.AddOption(option)
			}

			command, err := parser.Parse(testCase.argv)
			if diff := pretty.Compare(testCase.expectedErr, err); diff != "" {
				t.Errorf("%s Error diff:\n%s", t.Name(), diff)
			}
			if diff := pretty.Compare(testCase.expected, command); diff != "" {
				t.Errorf("%s Command diff:\n%s", t.Name(), diff)
			}
		})
	}
}
