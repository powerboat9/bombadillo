package main

import (
	"reflect"
	"testing"
)

func Test_WrapContent_Wrapped_Line_Length(t *testing.T) {
	type fields struct {
		WrappedContent []string
		RawContent     string
		Links          []string
		Location       Url
		ScrollPosition int
		FoundLinkLines []int
		SearchTerm     string
		SearchIndex    int
		FileType       string
		WrapWidth      int
		Color          bool
	}
	type args struct {
		width int
		color bool
	}

	// create a Url for use by the MakePage function
	url, _ := MakeUrl("gemini://rawtext.club")

	tests := []struct {
		name    string
		input   string
		expects []string
		args    args
	}{
		{
			"Short line that doesn't wrap",
			"0123456789\n",
			[]string{
				"0123456789",
				"",
			},
			args{
				10,
				false,
			},
		},
		{
			"Long line wrapped to 10 columns",
			"0123456789 123456789 123456789 123456789 123456789\n",
			[]string{
				"0123456789",
				" 123456789",
				" 123456789",
				" 123456789",
				" 123456789",
				"",
			},
			args{
				10,
				false,
			},
		},
		{
			"Unicode line endings that should not wrap",
			"LF\u000A" +
				"CR+LF\u000D\u000A" +
				"NEL\u0085" +
				"LS\u2028" +
				"PS\u2029",
			[]string{
				"LF",
				"CR+LF",
				"NEL",
				"LS",
				"PS",
				"",
			},
			args{
				10,
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := MakePage(url, tt.input, []string{""})
			p.WrapContent(tt.args.width-1, tt.args.color)
			if !reflect.DeepEqual(p.WrappedContent, tt.expects) {
				t.Errorf("Test failed - %s\nexpects %s\nactual  %s", tt.name, tt.expects, p.WrappedContent)
			}
		})
	}
}
