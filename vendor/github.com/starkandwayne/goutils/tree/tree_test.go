package tree_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/starkandwayne/goutils/tree"
)

var trim = regexp.MustCompile("\t+")

func drawsOk(t *testing.T, msg string, n tree.Node, s string) {
	got := strings.Trim(n.Draw(), "\n")
	want := strings.Trim(trim.ReplaceAllString(s, ""), "\n")
	if got != want {
		t.Errorf("%s failed\nexpected:\n[%s]\ngot:\n[%s]\n", msg, want, got)
	}
}

func pathsOk(t *testing.T, msg string, n tree.Node, want ...string) {
	got := n.Paths("/")
	if len(got) != len(want) {
		got_ := "    - " + strings.Join(got, "\n    - ") + "\n"
		want_ := "    - " + strings.Join(want, "\n    - ") + "\n"
		t.Errorf("%s failed\nexpected %d paths:\n%s\ngot %d paths:\n%s\n", msg, len(want), want_, len(got), got_)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("%s failed\npaths[%d] was incorrect\nexpected: [%s]\n     got: [%s]\n", msg, i, want[i], got[i])
		}
	}
}

func TestDrawing(t *testing.T) {
	drawsOk(t, "{a}",
		tree.New("a"), `
		.
		└── a`)

	drawsOk(t, "{a -> b}",
		tree.New("a", tree.New("b")), `
		.
		└── a
		    └── b`)

	drawsOk(t, "{a -> [b c]}",
		tree.New("a", tree.New("b"), tree.New("c")), `
		.
		└── a
		    ├── b
		    └── c`)

	drawsOk(t, "{a -> b -> c}",
		tree.New("a", tree.New("b", tree.New("c"))), `
		.
		└── a
		    └── b
		        └── c`)

	drawsOk(t, "{a -> [{b -> c} d]}",
		tree.New("a", tree.New("b", tree.New("c")), tree.New("d")), `
		.
		└── a
		    ├── b
		    │   └── c
		    └── d`)

	drawsOk(t, "{a -> [{b -> c -> e} d]}",
		tree.New("a", tree.New("b", tree.New("c", tree.New("e"))), tree.New("d")), `
		.
		└── a
		    ├── b
		    │   └── c
		    │       └── e
		    └── d`)

	drawsOk(t, "multiline node strings",
		tree.New("Alpha\n(first)\n",
			tree.New("Beta\n(second)\n",
				tree.New("Gamma\n(third)\n"),
			),
			tree.New("Delta\n(fourth)\n"),
		), `
		.
		└── Alpha
		    (first)
		    ├── Beta
		    │   (second)
		    │   └── Gamma
		    │       (third)
		    └── Delta
		        (fourth)`)
}

func TestPaths(t *testing.T) {
	pathsOk(t, "{a}",
		tree.New("a"),
		"a")

	pathsOk(t, "{a -> b}",
		tree.New("a", tree.New("b")),
		"a/b")

	pathsOk(t, "{a -> [b c]}",
		tree.New("a", tree.New("b"), tree.New("c")),
		"a/b",
		"a/c")

	pathsOk(t, "{a -> [{b -> c} d]",
		tree.New("a", tree.New("b", tree.New("c")), tree.New("d")),
		"a/b/c",
		"a/d")
}

func TestAppend(t *testing.T) {
	tr := tree.New("a", tree.New("b"))
	drawsOk(t, "{a -> b} before append", tr, `
		.
		└── a
		    └── b`)

	tr.Append(tree.New("c"))
	drawsOk(t, "{a -> [b c]} before append", tr, `
		.
		└── a
		    ├── b
		    └── c`)
}
