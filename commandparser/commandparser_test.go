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
		expected    commandparser.Command
	}{
		{
			"with long key",
			[]string{"./main", "--key", "value"},
			[]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{}),
		},
		{
			"with two long keys",
			[]string{"./main", "--key", "value", "--other", "another"},
			[]string{},
			makeCommand("./main", map[string]string{"key": "value", "other": "another"}, []string{}),
		},
		{
			"with short key",
			[]string{"./main", "-key", "value"},
			[]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{}),
		},
		{
			"with two short keys",
			[]string{"./main", "-key", "value", "-other", "another"},
			[]string{},
			makeCommand("./main", map[string]string{"key": "value", "other": "another"}, []string{}),
		},
		{
			"with positional",
			[]string{"./main", "value"},
			[]string{},
			makeCommand("./main", map[string]string{}, []string{"value"}),
		},
		{
			"with two positionals",
			[]string{"./main", "value", "another"},
			[]string{},
			makeCommand("./main", map[string]string{}, []string{"value", "another"}),
		},
		{
			"with long key and positional",
			[]string{"./main", "--key", "value", "positional"},
			[]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{"positional"}),
		},
		{
			"with short key and positional",
			[]string{"./main", "-key", "value", "positional"},
			[]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{"positional"}),
		},
		{
			"with boolean key",
			[]string{"./main", "-boolean"},
			[]string{"boolean"},
			makeCommand("./main", map[string]string{"boolean": "1"}, []string{}),
		},
		{
			"with multiple boolean keys",
			[]string{"./main", "-boolean", "-boolean2"},
			[]string{"boolean", "boolean2"},
			makeCommand("./main", map[string]string{"boolean": "1", "boolean2": "1"}, []string{}),
		},
		{
			"with dash positional",
			[]string{"./main", "-"},
			[]string{},
			makeCommand("./main", map[string]string{}, []string{"-"}),
		},
		{
			"mixed",
			[]string{"./main", "value", "--key", "value", "-other", "another", "-boolean", "last_pos"},
			[]string{"boolean"},
			makeCommand("./main", map[string]string{"key": "value", "other": "another", "boolean": "1"}, []string{"value", "last_pos"}),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := commandparser.CommandParser{}
			parser.AddBoolean(testCase.booleanKeys...)

			command := parser.Parse(testCase.argv)
			if diff := pretty.Compare(command, testCase.expected); diff != "" {
				t.Errorf("%s diff:\n%s", t.Name(), diff)
			}
		})
	}
}
