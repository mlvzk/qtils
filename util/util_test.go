package util_test

import (
	"testing"

	"github.com/mlvzk/qtils/util"
)

func TestLeftPad(t *testing.T) {
	testCases := []struct {
		input  string
		padStr string
		padLen int
		want   string
	}{
		{
			"aaa",
			"0",
			3,
			"aaa",
		},
		{
			"aa",
			"0",
			3,
			"0aa",
		},
		{
			"a",
			"0",
			3,
			"00a",
		},
		{
			"",
			"0",
			3,
			"000",
		},

		{
			"",
			"",
			3,
			"",
		},
		{
			"aaaa",
			"0",
			3,
			"aaaa",
		},
	}

	for _, testCase := range testCases {
		got := util.LeftPad(testCase.input, testCase.padStr, testCase.padLen)

		if got != testCase.want {
			t.Errorf("got != want -> %s != %s", got, testCase.want)
		}
	}
}

func TestRightPad(t *testing.T) {
	testCases := []struct {
		input  string
		padStr string
		padLen int
		want   string
	}{
		{
			"aaa",
			"0",
			3,
			"aaa",
		},
		{
			"aa",
			"0",
			3,
			"aa0",
		},
		{
			"a",
			"0",
			3,
			"a00",
		},
		{
			"",
			"0",
			3,
			"000",
		},
		{
			"",
			"",
			3,
			"",
		},
		{
			"aaaa",
			"0",
			3,
			"aaaa",
		},
	}

	for _, testCase := range testCases {
		got := util.RightPad(testCase.input, testCase.padStr, testCase.padLen)

		if got != testCase.want {
			t.Errorf("got != want -> %s != %s", got, testCase.want)
		}
	}
}
