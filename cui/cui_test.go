package cui

import (
	"reflect"
	"testing"
)

// tests related to issue 31
func Test_wrapLines_space_preservation(t *testing.T) {
	tables := []struct {
		testinput      []string
		expectedoutput []string
		linelength     int
	}{
		{
			//normal sentence - 20 characters - should not wrap
			[]string{"it is her fav thingy"},
			[]string{"it is her fav thingy"},
			20,
		},
		{
			//normal sentence - more than 20 characters - should wrap with a space at the end of the first line
			[]string{"it is her favourite thing in the world"},
			[]string{
				"it is her favourite ",
				"thing in the world",
			},
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
	teststring := []string{
		"0123456789",
		"a really long line that will prolly be wrapped",
		"a l i n e w i t h a l o t o f w o r d s",
		"onehugelongwordaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		wrapLines(teststring, 20)
	}
}