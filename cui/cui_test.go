package cui

import (
	"reflect"
	"testing"
)

func Test_wrapLines_doesnt_break_indents(t *testing.T) {
	indent := "           "
	tables := []struct {
		testinput      []string
		expectedoutput []string
		linelength     int
	}{
		{
			//20 character input - should not wrap
			[]string{indent + "012345678"},
			[]string{indent + "012345678"},
			20,
		},
		{
			//21 character input - should wrap
			[]string{indent + "0123456789"},
			[]string{indent + "0123456789"},
			20,
		},
	}

	for _, table := range tables {
		output := wrapLines(table.testinput, table.linelength)

		if !reflect.DeepEqual(output, table.expectedoutput) {
			t.Errorf("Expected %v, got %v", table.expectedoutput, output)
		}
	}
}
