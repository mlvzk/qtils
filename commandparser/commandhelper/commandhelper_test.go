package commandhelper_test

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mlvzk/qtils/commandparser/commandhelper"
)

func TestVerifyArgs(t *testing.T) {
	helper := commandhelper.New()
	helper.EatOption(commandhelper.NewOption("key").Build())

	input := map[string]string{"key": "1", "invalidkey": "1"}
	expected := []string{"invalidkey"}

	if diff := pretty.Compare(helper.VerifyArgs(input), expected); diff != "" {
		t.Errorf("%s diff:\n%s", t.Name(), diff)
	}
}
