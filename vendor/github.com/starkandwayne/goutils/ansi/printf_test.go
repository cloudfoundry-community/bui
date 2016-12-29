package ansi

import (
	"testing"
)

func TestColorizer(t *testing.T) {
	var tests = []struct {
		In  string
		Out string
		Not string
	}{
		{"@k{color}", "\033[00;30mcolor\033[00m", "color"},
		{"@K{COLOR}", "\033[01;30mCOLOR\033[00m", "COLOR"},

		{"@r{color}", "\033[00;31mcolor\033[00m", "color"},
		{"@R{COLOR}", "\033[01;31mCOLOR\033[00m", "COLOR"},

		{"@g{color}", "\033[00;32mcolor\033[00m", "color"},
		{"@G{COLOR}", "\033[01;32mCOLOR\033[00m", "COLOR"},

		{"@y{color}", "\033[00;33mcolor\033[00m", "color"},
		{"@Y{COLOR}", "\033[01;33mCOLOR\033[00m", "COLOR"},

		{"@b{color}", "\033[00;34mcolor\033[00m", "color"},
		{"@B{COLOR}", "\033[01;34mCOLOR\033[00m", "COLOR"},

		{"@m{color}", "\033[00;35mcolor\033[00m", "color"},
		{"@p{color}", "\033[00;35mcolor\033[00m", "color"},
		{"@M{COLOR}", "\033[01;35mCOLOR\033[00m", "COLOR"},
		{"@P{COLOR}", "\033[01;35mCOLOR\033[00m", "COLOR"},

		{"@c{color}", "\033[00;36mcolor\033[00m", "color"},
		{"@C{COLOR}", "\033[01;36mCOLOR\033[00m", "COLOR"},

		{"@w{color}", "\033[00;37mcolor\033[00m", "color"},
		{"@W{COLOR}", "\033[01;37mCOLOR\033[00m", "COLOR"},

		{"@k{black} and @r{red}", "\033[00;30mblack\033[00m and \033[00;31mred\033[00m", "black and red"},
		{"error: @R{%s}", "error: \033[01;31m%s\033[00m", "error: %s"},

		{"@*{RAINBOW}", "\033[01;31mR\033[00m\033[01;33mA\033[00m\033[01;32mI\033[00m\033[01;36mN\033[00m\033[01;34mB\033[00m\033[01;35mO\033[00m\033[01;31mW\033[00m", "RAINBOW"},

		{"@s@d@l@f", "@s@d@l@f", "@s@d@l@f"},
		{"host error: %s", "host error: %s", "host error: %s"},
	}

	for _, test := range tests {
		Color(true)
		if colorize(test.In) != test.Out {
			t.Errorf("colorize(`%s`) was `%s`, not `%s`", test.In, colorize(test.In), test.Out)
		}
		Color(false)
		if colorize(test.In) != test.Not {
			t.Errorf("colorize(`%s`) was `%s`, not `%s`", test.In, colorize(test.In), test.Not)
		}
	}
}
