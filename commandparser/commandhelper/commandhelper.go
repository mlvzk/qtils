package commandhelper

import (
	"errors"

	"github.com/mlvzk/qtils/color"
	"github.com/mlvzk/qtils/commandparser"
	"github.com/mlvzk/qtils/util"
)

type Helper struct {
	name        string
	version     string
	authors     []string
	options     []string
	usages      []string
	description map[string]string
	defaults    map[string]string
	required    map[string]bool
	optionSpecs []OptionSpec
}

func New() *Helper {
	return &Helper{
		"",
		"v0.0.0",
		[]string{},
		[]string{},
		[]string{},
		map[string]string{},
		map[string]string{},
		map[string]bool{},
		[]OptionSpec{},
	}
}

func (helper *Helper) SetName(name string) {
	helper.name = name
}

func (helper *Helper) SetVersion(version string) {
	helper.version = version
}

func (helper *Helper) AddAuthor(author ...string) {
	helper.authors = append(helper.authors, author...)
}

func (helper *Helper) AddUsage(usage ...string) {
	helper.usages = append(helper.usages, usage...)
}

func (helper *Helper) EatOption(options ...OptionSpec) []commandparser.Option {
	parserOptions := make([]commandparser.Option, len(options))
	for i, option := range options {
		helper.optionSpecs = append(helper.optionSpecs, option)
		helper.options = append(helper.options, option.GetKey())
		helper.defaults[option.GetKey()] = option.GetDefault()
		helper.description[option.GetKey()] = option.GetDescription()
		helper.required[option.GetKey()] = option.IsRequired()
		parserOptions[i] = options[i]
	}

	return parserOptions
}

func (helper *Helper) FillDefaults(args map[string]string) map[string]string {
	newArgs := map[string]string{}

	for k, v := range args {
		newArgs[k] = v
	}

	for k, v := range helper.defaults {
		if _, isSet := newArgs[k]; !isSet {
			newArgs[k] = v
		}
	}

	return newArgs
}

func (helper *Helper) VerifyArgs(args map[string]string) []error {
	errs := []error{}

	for requiredKey, requiredValue := range helper.required {
		if requiredValue == false {
			continue
		}

		if _, exists := args[requiredKey]; !exists {
			errs = append(errs, errors.New("Missing required argument '"+requiredKey+"'"))
		}
	}

	return errs
}

// returns longest **individual** columns
func getLongestColumns(rows [][]string) (longest []int) {
	for _, r := range rows {
		for i, c := range r {
			if len(longest) < (i + 1) {
				longest = append(longest, len(c))
			}

			if len(c) > longest[i] {
				longest[i] = len(c)
			}
		}
	}
	return
}

func join(strings []string, sep string) string {
	var result string

	for i, v := range strings {
		result += v
		if i != len(strings)-1 {
			result += sep
		}
	}

	return result
}

func (helper *Helper) Help() string {
	var result string

	result = color.Green(helper.name) + " " + helper.version

	if len(helper.usages) != 0 {
		result += "\n\n" + color.Section("USAGE:")
		for _, usage := range helper.usages {
			result += "\n\t" + usage
		}
	}

	if len(helper.optionSpecs) != 0 {
		var optionRows [][]string

		result += "\n\n" + color.Section("OPTIONS:")
		for _, option := range helper.optionSpecs {
			aliasesStr := join(option.GetAliases(), ", -")
			if len(option.GetAliases()) != 0 {
				aliasesStr = ", -" + join(option.GetAliases(), ", -")
			}

			firstColumn := "--" + option.GetKey() + aliasesStr
			secondColumn := option.GetDescription()
			if option.GetDefault() != "" {
				secondColumn += color.Info(" (default: " + option.GetDefault() + ")")
			}
			if option.IsRequired() {
				secondColumn += color.Important(" (required)")
			}
			if option.IsArrayed() {
				secondColumn += color.Info(" (accepts multiple)")
			}
			if !option.IsBoolean() {
				firstColumn += " <value>"
			}

			optionRows = append(optionRows, []string{firstColumn, secondColumn})
		}

		longestColumns := getLongestColumns(optionRows)
		padCol := func(column string, index int) string {
			return util.RightPad(column, " ", longestColumns[index]+4)
		}
		for _, row := range optionRows {
			result += "\n\t" + color.Green(padCol(row[0], 0)) + row[1]
		}
	}

	if len(helper.authors) != 0 {
		if len(helper.authors) == 1 {
			result += "\n\n" + color.Section("Author: ") + helper.authors[0]
		} else {
			result += "\n\n" + color.Section("Authors: ") + join(helper.authors, ", ")
		}
	}

	return result + "\n"
}

type OptionSpec interface {
	commandparser.Option
	GetDefault() string
	GetDescription() string
	IsRequired() bool
}

type OptionBuilder interface {
	Alias(key ...string) OptionBuilder
	Default(value string) OptionBuilder
	Description(value string) OptionBuilder
	Required() OptionBuilder
	Arrayed() OptionBuilder
	Boolean() OptionBuilder
	// TODO: implement Validate(func)
	Build() OptionSpec
}

type Option struct {
	key          string
	defaultValue string
	description  string
	aliases      []string
	required     bool
	arrayed      bool
	boolean      bool
}

func NewOption(key string) OptionBuilder {
	return &Option{
		key:          key,
		defaultValue: "",
		aliases:      []string{},
		required:     false,
		arrayed:      false,
		boolean:      false,
	}
}

func (option *Option) GetKey() string {
	return option.key
}

func (option *Option) Alias(key ...string) OptionBuilder {
	option.aliases = append(option.aliases, key...)

	return option
}

func (option *Option) GetAliases() []string {
	return option.aliases
}

func (option *Option) Default(value string) OptionBuilder {
	option.defaultValue = value

	return option
}

func (option *Option) GetDefault() string {
	return option.defaultValue
}

func (option *Option) Description(description string) OptionBuilder {
	option.description = description

	return option
}

func (option *Option) GetDescription() string {
	return option.description
}

func (option *Option) Required() OptionBuilder {
	option.required = true

	return option
}

func (option *Option) IsRequired() bool {
	return option.required
}

func (option *Option) Arrayed() OptionBuilder {
	option.arrayed = true

	return option
}

func (option *Option) IsArrayed() bool {
	return option.arrayed
}

func (option *Option) Boolean() OptionBuilder {
	option.boolean = true

	return option
}

func (option *Option) IsBoolean() bool {
	return option.boolean
}

func (option *Option) Build() OptionSpec {
	return option
}
