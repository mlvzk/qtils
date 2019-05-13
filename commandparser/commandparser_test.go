package commandparser_test

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mlvzk/qtils/commandparser"
	"github.com/mlvzk/qtils/commandparser/commandhelper"
)

func TestParseCommand(t *testing.T) {
	makeCommand := func(exe string, args map[string]string, positionals []string, arrayed map[string][]string, booleans map[string]bool) commandparser.Command {
		return commandparser.Command{
			Exe:         exe,
			Args:        args,
			Positionals: positionals,
			Arrayed:     arrayed,
			Booleans:    booleans,
		}
	}

	cases := []struct {
		name        string
		argv        []string
		booleanKeys []string
		arrayedKeys []string
		aliases     map[string][]string
		expected    commandparser.Command
	}{
		{
			"with long key",
			[]string{"./main", "--key", "value"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with two long keys",
			[]string{"./main", "--key", "value", "--other", "another"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value", "other": "another"}, []string{}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with short key",
			[]string{"./main", "-key", "value"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with two short keys",
			[]string{"./main", "-key", "value", "-other", "another"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value", "other": "another"}, []string{}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with positional",
			[]string{"./main", "value"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"value"}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with two positionals",
			[]string{"./main", "value", "another"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"value", "another"}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with long key and positional",
			[]string{"./main", "--key", "value", "positional"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{"positional"}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with short key and positional",
			[]string{"./main", "-key", "value", "positional"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{"positional"}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with boolean key",
			[]string{"./main", "-boolean"},
			[]string{"boolean"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{}, map[string][]string{}, map[string]bool{"boolean": true}),
		},
		{
			"with multiple boolean keys",
			[]string{"./main", "-boolean", "-boolean2"},
			[]string{"boolean", "boolean2"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{}, map[string][]string{}, map[string]bool{"boolean": true, "boolean2": true}),
		},
		{
			"with dash positional",
			[]string{"./main", "-"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"-"}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with short key alias",
			[]string{"./main", "-o", "file"},
			[]string{},
			[]string{},
			map[string][]string{"output": {"o"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with long key alias",
			[]string{"./main", "--o", "file"},
			[]string{},
			[]string{},
			map[string][]string{"output": {"o"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}, map[string][]string{}, map[string]bool{}),
		},
		{
			"with multiple aliases to same key",
			[]string{"./main", "--out", "file"},
			[]string{},
			[]string{},
			map[string][]string{"output": {"o", "out"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}, map[string][]string{}, map[string]bool{}),
		},
		{
			"sed -n --quiet --silent",
			[]string{"sed", "-n", "--quiet", "--silent"},
			[]string{"quiet"},
			[]string{},
			map[string][]string{"quiet": {"n", "silent"}},
			makeCommand("sed", map[string]string{}, []string{}, map[string][]string{}, map[string]bool{"quiet": true}),
		},
		{
			"with arrayed key",
			[]string{"./main", "--key", "value1", "--key", "value2"},
			[]string{},
			[]string{"key"},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{}, map[string][]string{"key": {"value1", "value2"}}, map[string]bool{}),
		},
		{
			"with same key both boolean and arrayed",
			[]string{"./main", "--verbose", "--verbose"},
			[]string{"verbose"},
			[]string{"verbose"},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{}, map[string][]string{"verbose": {"1", "1"}}, map[string]bool{}),
		},
		{
			"mixed",
			[]string{"./main", "first_pos", "-v", "-v", "-v", "--key", "value", "-other", "another", "-boolean", "-b2", "-arrayed_key", "arrayed1", "-arrayed_key", "arrayed2", "last_pos"},
			[]string{"boolean", "boolean2", "verbose"},
			[]string{"arrayed_key", "verbose"},
			map[string][]string{"boolean2": {"b2"}, "verbose": {"v"}},
			makeCommand(
				"./main",
				map[string]string{"key": "value", "other": "another"},
				[]string{"first_pos", "last_pos"},
				map[string][]string{"arrayed_key": {"arrayed1", "arrayed2"}, "verbose": {"1", "1", "1"}},
				map[string]bool{"boolean": true, "boolean2": true},
			),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := commandparser.New()

			options := map[string]commandhelper.OptionBuilder{}

			for _, booleanKey := range testCase.booleanKeys {
				if _, exists := options[booleanKey]; !exists {
					options[booleanKey] = commandhelper.NewOption(booleanKey)
				}

				options[booleanKey].Boolean()
			}
			for _, arrayedKey := range testCase.arrayedKeys {
				if _, exists := options[arrayedKey]; !exists {
					options[arrayedKey] = commandhelper.NewOption(arrayedKey)
				}

				options[arrayedKey].Arrayed()
			}
			for key, values := range testCase.aliases {
				if _, exists := options[key]; !exists {
					options[key] = commandhelper.NewOption(key)
				}

				options[key].Alias(values...)
			}

			for _, option := range options {
				parser.AddOption(option.Build())
			}

			command := parser.Parse(testCase.argv)
			if diff := pretty.Compare(command, testCase.expected); diff != "" {
				t.Errorf("%s diff:\n%s", t.Name(), diff)
			}
		})
	}
}
