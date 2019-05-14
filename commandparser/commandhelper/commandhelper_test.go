package commandhelper_test

import (
	"errors"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mlvzk/qtils/commandparser/commandhelper"
)

func TestVerifyArgs(t *testing.T) {
	helper := commandhelper.New()
	helper.EatOption(commandhelper.NewOption("key").Required().Build())

	input := map[string]string{}
	expected := []error{errors.New("Missing required argument 'key'")}

	if diff := pretty.Compare(helper.VerifyArgs(input), expected); diff != "" {
		t.Errorf("%s diff:\n%s", t.Name(), diff)
	}
}
