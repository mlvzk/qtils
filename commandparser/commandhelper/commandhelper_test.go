package commandhelper_test

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mlvzk/qtils/commandparser/commandhelper"
)

var update = flag.Bool("update", false, "update .golden files")

func TestVerify(t *testing.T) {
	helper := commandhelper.New()
	helper.EatOption(
		commandhelper.NewOption("key").Required(),
		commandhelper.NewOption("arrayed_key_required").Arrayed().Required(),
		commandhelper.NewOption("arrayed_key").Arrayed().Validate(func(value string) error {
			return errors.New("this shouldn't happen")
		}),
		commandhelper.NewOption("validated_key").Validate(func(value string) error {
			return errors.New("expected error from validated_key.Validate")
		}),
		commandhelper.NewOption("accepts_kv").ValidateBind(commandhelper.ValidateKeyValue("=")),
		commandhelper.NewOption("accepts_kv_invalid").ValidateBind(commandhelper.ValidateKeyValue(" = ")),
	)

	inputArgs := map[string]string{
		"accepts_kv":         "k=v",
		"accepts_kv_invalid": "test",
	}
	inputArrayed := map[string][]string{}
	expected := []error{
		errors.New("missing required argument 'key'"),
		errors.New("missing required argument 'arrayed_key_required'"),
		errors.New("expected error from validated_key.Validate"),
		errors.New("invalid value of key 'accepts_kv_invalid'; value does not contain the ' = ' delimiter"),
	}

	if diff := pretty.Compare(expected, helper.Verify(inputArgs, inputArrayed)); diff != "" {
		t.Errorf("%s diff:\n%s", t.Name(), diff)
	}
}

func TestHelp(t *testing.T) {
	helper := commandhelper.New()
	helper.SetName("test")
	helper.SetVersion("v0.0.0")
	helper.AddUsage("test [OPTIONS] <positionals...>", "test [OPTIONS]")
	helper.AddAuthor("mlvzk", "tester")

	helper.EatOption(
		commandhelper.NewOption("key").Description("Simple key").Default("test"),
		commandhelper.NewOption("port").Description("Port the server should listen on").Required(),
		commandhelper.NewOption("verbose").Boolean().Arrayed().Alias("v", "loud").Description("Verbose flag"),
		commandhelper.NewOption("multiline_description").Description(`this is a
multiline description`),
	)

	got := helper.Help()

	golden := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		file, err := os.Create(golden)
		if err != nil {
			t.Fatalf(err.Error())
		}

		io.Copy(file, strings.NewReader(got))
		file.Close()
	}

	file, err := os.Open(golden)
	if err != nil {
		t.Fatalf(err.Error())
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := strings.Replace(string(bytes), "\r\n", "\n", -1) // fix for windows CI

	if diff := pretty.Compare(expected, got); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}

func BenchmarkHelp(b *testing.B) {
	helper := commandhelper.New()
	helper.SetName("test")
	helper.SetVersion("v0.0.0")
	helper.AddUsage("test [OPTIONS] <positionals...>", "test [OPTIONS]")
	helper.AddAuthor("mlvzk", "tester")

	helper.EatOption(
		commandhelper.NewOption("key").Description("Simple key").Default("test"),
		commandhelper.NewOption("port").Description("Port the server should listen on").Required(),
		commandhelper.NewOption("verbose").Boolean().Arrayed().Alias("v", "loud").Description("Verbose flag"),
		commandhelper.NewOption("multiline_description").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description1").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description2").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description3").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description4").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description5").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description6").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description7").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description8").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description9").Description(`this is a
multiline description`),
		commandhelper.NewOption("multiline_description10").Description(`this is a
multiline description`),
	)

	b.ResetTimer()
	for i := 0; i < 1000; i++ {
		helper.Help()
	}
}

func TestValidateInt(t *testing.T) {
	if err := commandhelper.ValidateInt("key123")("abcd"); err == nil {
		t.Errorf("Expected error but got nil")
	} else {
		want := commandhelper.NewInvalidValue("key123", "value must be an integer")
		if err.Error() != want.Error() {
			t.Errorf("got != want -> '%+v' != '%+v'", err, want)
		}
	}

	if err := commandhelper.ValidateInt("key123")("1234"); err != nil {
		t.Errorf("Expected error to be nil, but got: '%+v'", err)
	}
}

func TestValidateSelection(t *testing.T) {
	if err := commandhelper.ValidateSelection("apple", "pear", "pie")("key123")("unknown"); err == nil {
		t.Errorf("Expected error but got nil")
	} else {
		want := commandhelper.NewInvalidValue("key123", "value must be one of: apple, pear, pie")
		if err.Error() != want.Error() {
			t.Errorf("got != want -> '%+v' != '%+v'", err, want)
		}
	}

	if err := commandhelper.ValidateSelection("apple", "pear", "pie")("key123")("pie"); err != nil {
		t.Errorf("Expected error to be nil, but got: '%+v'", err)
	}
}
