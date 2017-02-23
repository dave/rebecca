package rebecca

import (
	"strconv"
	"testing"
)

func TestExtractSections(t *testing.T) {
	comment := "foo. bar. baz. qux. quz."
	tests := []struct {
		sections string
		expected string
	}{
		{
			sections: "0",
			expected: "foo.",
		},
		{
			sections: "0:1",
			expected: "foo.",
		},
		{
			sections: "0:2",
			expected: "foo. bar.",
		},
		{
			sections: "1:2",
			expected: "bar.",
		},
		{
			sections: "2:4",
			expected: "baz. qux.",
		},
		{
			sections: "4:",
			expected: "quz.",
		},
		{
			sections: "1:",
			expected: "bar. baz. qux. quz.",
		},
		{
			sections: "2:",
			expected: "baz. qux. quz.",
		},
		{
			sections: ":1",
			expected: "foo.",
		},
		{
			sections: ":2",
			expected: "foo. bar.",
		},
		{
			sections: ":4",
			expected: "foo. bar. baz. qux.",
		},
	}
	for _, test := range tests {
		found := extractSections("Spec["+test.sections+"]", test.sections, comment)
		if found != test.expected {
			t.Fatalf("SectionSpec: %s. Expected %s. Found %s.", strconv.Quote(test.sections), strconv.Quote(test.expected), strconv.Quote(found))
		}
	}
}
