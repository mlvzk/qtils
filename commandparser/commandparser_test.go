package commandparser_test

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mlvzk/qtils/commandparser"
)

func TestParseCommand(t *testing.T) {
	makeCommand := func(exe string, args map[string]string, positionals []string) commandparser.Command {
		return commandparser.Command{
			Exe:         exe,
			Args:        args,
			Positionals: positionals,
		}
	}

	cases := []struct {
		name        string
		argv        []string
		booleanKeys []string
		aliases     map[string][]string
		expected    commandparser.Command
	}{
		{
			"with long key",
			[]string{"./main", "--key", "value"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{}),
		},
		{
			"with two long keys",
			[]string{"./main", "--key", "value", "--other", "another"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value", "other": "another"}, []string{}),
		},
		{
			"with short key",
			[]string{"./main", "-key", "value"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{}),
		},
		{
			"with two short keys",
			[]string{"./main", "-key", "value", "-other", "another"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value", "other": "another"}, []string{}),
		},
		{
			"with positional",
			[]string{"./main", "value"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"value"}),
		},
		{
			"with two positionals",
			[]string{"./main", "value", "another"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"value", "another"}),
		},
		{
			"with long key and positional",
			[]string{"./main", "--key", "value", "positional"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{"positional"}),
		},
		{
			"with short key and positional",
			[]string{"./main", "-key", "value", "positional"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{"positional"}),
		},
		{
			"with boolean key",
			[]string{"./main", "-boolean"},
			[]string{"boolean"},
			map[string][]string{},
			makeCommand("./main", map[string]string{"boolean": "1"}, []string{}),
		},
		{
			"with multiple boolean keys",
			[]string{"./main", "-boolean", "-boolean2"},
			[]string{"boolean", "boolean2"},
			map[string][]string{},
			makeCommand("./main", map[string]string{"boolean": "1", "boolean2": "1"}, []string{}),
		},
		{
			"with dash positional",
			[]string{"./main", "-"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"-"}),
		},
		{
			"with short key alias",
			[]string{"./main", "-o", "file"},
			[]string{},
			map[string][]string{"output": {"o"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}),
		},
		{
			"with long key alias",
			[]string{"./main", "--o", "file"},
			[]string{},
			map[string][]string{"output": {"o"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}),
		},
		{
			"with multiple aliases to same key",
			[]string{"./main", "--out", "file"},
			[]string{},
			map[string][]string{"output": {"o", "out"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}),
		},
		{
			"sed -n --quiet --silent",
			[]string{"sed", "-n", "--quiet", "--silent"},
			[]string{"quiet"},
			map[string][]string{"quiet": {"n", "silent"}},
			makeCommand("sed", map[string]string{"quiet": "1"}, []string{}),
		},
		{
			"mixed",
			[]string{"./main", "first_pos", "--key", "value", "-other", "another", "-boolean", "-b2", "last_pos"},
			[]string{"boolean", "boolean2"},
			map[string][]string{"boolean2": {"b2"}},
			makeCommand(
				"./main",
				map[string]string{"key": "value", "other": "another", "boolean": "1", "boolean2": "1"},
				[]string{"first_pos", "last_pos"},
			),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := commandparser.New()
			parser.AddBoolean(testCase.booleanKeys...)

			for key, values := range testCase.aliases {
				parser.AddAliases(key, values...)
			}

			command := parser.Parse(testCase.argv)
			if diff := pretty.Compare(command, testCase.expected); diff != "" {
				t.Errorf("%s diff:\n%s", t.Name(), diff)
			}
		})
	}
}
