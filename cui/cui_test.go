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

func Benchmark_wrapLines(b *testing.B) {
	indent := "           "
	teststring := []string{
		indent + "0123456789\n",
		indent + "a really long line that will prolly be wrapped\n",
		indent + "a l i n e w i t h a l o t o f w o r d s\n",
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		wrapLines(teststring, 20)
	}
}
