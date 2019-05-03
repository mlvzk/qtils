package commandparser_test

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mlvzk/qtils/commandparser"
)

func TestParseCommand(t *testing.T) {
	makeCommand := func(exe string, args map[string]string, positionals []string, arrayed map[string][]string) commandparser.Command {
		return commandparser.Command{
			Exe:         exe,
			Args:        args,
			Positionals: positionals,
			Arrayed:     arrayed,
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
			makeCommand("./main", map[string]string{"key": "value"}, []string{}, map[string][]string{}),
		},
		{
			"with two long keys",
			[]string{"./main", "--key", "value", "--other", "another"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value", "other": "another"}, []string{}, map[string][]string{}),
		},
		{
			"with short key",
			[]string{"./main", "-key", "value"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{}, map[string][]string{}),
		},
		{
			"with two short keys",
			[]string{"./main", "-key", "value", "-other", "another"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value", "other": "another"}, []string{}, map[string][]string{}),
		},
		{
			"with positional",
			[]string{"./main", "value"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"value"}, map[string][]string{}),
		},
		{
			"with two positionals",
			[]string{"./main", "value", "another"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"value", "another"}, map[string][]string{}),
		},
		{
			"with long key and positional",
			[]string{"./main", "--key", "value", "positional"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{"positional"}, map[string][]string{}),
		},
		{
			"with short key and positional",
			[]string{"./main", "-key", "value", "positional"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"key": "value"}, []string{"positional"}, map[string][]string{}),
		},
		{
			"with boolean key",
			[]string{"./main", "-boolean"},
			[]string{"boolean"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"boolean": "1"}, []string{}, map[string][]string{}),
		},
		{
			"with multiple boolean keys",
			[]string{"./main", "-boolean", "-boolean2"},
			[]string{"boolean", "boolean2"},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{"boolean": "1", "boolean2": "1"}, []string{}, map[string][]string{}),
		},
		{
			"with dash positional",
			[]string{"./main", "-"},
			[]string{},
			[]string{},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{"-"}, map[string][]string{}),
		},
		{
			"with short key alias",
			[]string{"./main", "-o", "file"},
			[]string{},
			[]string{},
			map[string][]string{"output": {"o"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}, map[string][]string{}),
		},
		{
			"with long key alias",
			[]string{"./main", "--o", "file"},
			[]string{},
			[]string{},
			map[string][]string{"output": {"o"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}, map[string][]string{}),
		},
		{
			"with multiple aliases to same key",
			[]string{"./main", "--out", "file"},
			[]string{},
			[]string{},
			map[string][]string{"output": {"o", "out"}},
			makeCommand("./main", map[string]string{"output": "file"}, []string{}, map[string][]string{}),
		},
		{
			"sed -n --quiet --silent",
			[]string{"sed", "-n", "--quiet", "--silent"},
			[]string{"quiet"},
			[]string{},
			map[string][]string{"quiet": {"n", "silent"}},
			makeCommand("sed", map[string]string{"quiet": "1"}, []string{}, map[string][]string{}),
		},
		{
			"with arrayed key",
			[]string{"./main", "--key", "value1", "--key", "value2"},
			[]string{},
			[]string{"key"},
			map[string][]string{},
			makeCommand("./main", map[string]string{}, []string{}, map[string][]string{"key": {"value1", "value2"}}),
		},
		{
			"mixed",
			[]string{"./main", "first_pos", "--key", "value", "-other", "another", "-boolean", "-b2", "-arrayed_key", "arrayed1", "-arrayed_key", "arrayed2", "last_pos"},
			[]string{"boolean", "boolean2"},
			[]string{"arrayed_key"},
			map[string][]string{"boolean2": {"b2"}},
			makeCommand(
				"./main",
				map[string]string{"key": "value", "other": "another", "boolean": "1", "boolean2": "1"},
				[]string{"first_pos", "last_pos"},
				map[string][]string{"arrayed_key": {"arrayed1", "arrayed2"}},
			),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := commandparser.New()
			parser.AddBoolean(testCase.booleanKeys...)
			parser.AddArrayed(testCase.arrayedKeys...)

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
