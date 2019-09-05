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
			//indented long word - 20 characters - should not wrap
			[]string{indent + "012345678"},
			[]string{indent + "012345678"},
			20,
		},
		{
			//indented long word - 21 characters - should wrap
			[]string{indent + "0123456789"},
			[]string{indent + "012345678", indent + "9"},
			20,
		},
		{
			//indented really long word - should wrap
			[]string{indent + "0123456789zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"},
			[]string{
				indent + "012345678",
				indent + "9zzzzzzzz",
				indent + "zzzzzzzzz",
				indent + "zzzzzzzzz",
				indent + "zzzzz"},
			20,
		}, {
			//non-indented long word - 20 characters - should not wrap
			[]string{"01234567890123456789"},
			[]string{"01234567890123456789"},
			20,
		},
		{
			//non-indented long word - 21 characters - should wrap
			[]string{"01234567890123456789a"},
			[]string{"01234567890123456789", "a"},
			20,
		},
		{
			//non-indented really long word - should wrap
			[]string{"01234567890123456789zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"},
			[]string{
				"01234567890123456789",
				"zzzzzzzzzzzzzzzzzzzz",
				"zzzzzzzzzzzzzzzzzzzz",
				"zzzzzzzzzzzzzzzzzzzz",
				"zzzz"},
			20,
		},
		{
			//indented normal sentence - 20 characters - should not wrap
			[]string{indent + "it is her"},
			[]string{indent + "it is her"},
			20,
		},
		{
			//indented normal sentence - more than 20 characters - should wrap
			[]string{indent + "it is her favourite thing in the world"},
			[]string{
				indent + "it is her",
				indent + "favourite",
				indent + "thing in",
				indent + "the world",
			},
			20,
		},
		{
			//non-indented normal sentence - 20 characters - should not wrap
			[]string{"it is her fav thingy"},
			[]string{"it is her fav thingy"},
			20,
		},
		{
			//non-indented normal sentence - more than 20 characters - should wrap
			[]string{"it is her favourite thing in the world"},
			[]string{
				"it is her favourite",
				"thing in the world",
			},
			20,
		},
		//TODO further tests
		//lines that are just spaces don't get misidentified as indents and then mangled
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
