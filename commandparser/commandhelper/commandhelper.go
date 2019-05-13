package commandhelper

import (
	"fmt"

	"github.com/mlvzk/qtils/commandparser"
)

type Helper struct {
	name     string
	version  string
	authors  []string
	options  []string
	details  map[string]string
	defaults map[string]string
	required map[string]bool
}

func New() *Helper {
	return &Helper{
		"",
		"v0.0.0",
		[]string{},
		[]string{},
		map[string]string{},
		map[string]string{},
		map[string]bool{},
	}
}

func (helper *Helper) SetName(name string) {
	helper.name = name
}

func (helper *Helper) SetVersion(version string) {
	helper.version = version
}

func (helper *Helper) AddAuthors(author ...string) {
	helper.authors = append(helper.authors, author...)
}

func (helper *Helper) EatOption(options ...OptionSpec) []commandparser.Option {
	for _, option := range options {
		helper.options = append(helper.options, option.GetKey())
		helper.defaults[option.GetKey()] = option.GetDefault()
		helper.details[option.GetKey()] = option.GetDescription()
		helper.required[option.GetKey()] = option.IsRequired()
	}

	parserOptions := make([]commandparser.Option, len(options))
	for i := range options {
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

argLoop:
	for arg := range args {
		for _, option := range helper.options {
			if arg == option {
				continue argLoop
			}
		}

		errs = append(errs, fmt.Errorf("Invalid argument '%s'", arg))
	}

	for requiredKey, requiredValue := range helper.required {
		if requiredValue == false {
			continue
		}

		if _, exists := args[requiredKey]; !exists {
			errs = append(errs, fmt.Errorf("Missing required argument '%s'", requiredKey))
		}
	}

	return errs
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
