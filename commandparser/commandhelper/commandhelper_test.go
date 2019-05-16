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
		commandhelper.NewOption("key").Required().Build(),
		commandhelper.NewOption("arrayed_key_required").Arrayed().Required().Build(),
		commandhelper.NewOption("arrayed_key").Arrayed().Validate(func(value string) error {
			return errors.New("this shouldn't happen")
		}).Build(),
		commandhelper.NewOption("validated_key").Validate(func(value string) error {
			return errors.New("expected error from validated_key.Validate")
		}).Build(),
	)

	inputArgs := map[string]string{}
	inputArrayed := map[string][]string{}
	expected := []error{
		errors.New("missing required argument 'key'"),
		errors.New("missing required argument 'arrayed_key_required'"),
		errors.New("expected error from validated_key.Validate"),
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
		commandhelper.NewOption("key").Description("Simple key").Default("test").Build(),
		commandhelper.NewOption("port").Description("Port the server should listen on").Required().Build(),
		commandhelper.NewOption("verbose").Boolean().Arrayed().Alias("v", "loud").Description("Verbose flag").Build(),
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

	if diff := pretty.Compare(string(bytes), got); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}
