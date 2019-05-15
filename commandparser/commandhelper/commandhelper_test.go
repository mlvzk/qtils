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

func TestVerifyArgs(t *testing.T) {
	helper := commandhelper.New()
	helper.EatOption(commandhelper.NewOption("key").Required().Build())

	input := map[string]string{}
	expected := []error{errors.New("Missing required argument 'key'")}

	if diff := pretty.Compare(helper.VerifyArgs(input), expected); diff != "" {
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
		commandhelper.NewOption("verbose").Boolean().Alias("v", "loud").Description("Verbose flag").Build(),
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

	if diff := pretty.Compare(got, string(bytes)); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}
