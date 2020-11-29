package main

import (
	"reflect"
	"testing"
)

func TestSplitSections(t *testing.T) {
	specs := []struct {
		input string
		exp   []string
	}{
		{
			input: `
some stuff before
hello that is nice
# Hello
another section
## and another
another section
`,
			exp: []string{
				"\nsome stuff before\nhello that is nice\n",
				"# Hello\nanother section\n",
				"## and another\nanother section\n",
			},
		},
		{
			input: "Hello",
			exp:   []string{"Hello\n"},
		},
	}

	for _, spec := range specs {
		if got := splitSections(spec.input); !reflect.DeepEqual(got, spec.exp) {
			t.Errorf("Got %+v, expected %+v", got, spec.exp)
		}
	}
}
