package formats

import "testing"

func TestGetParser(t *testing.T) {
	type test struct {
		fn       string
		expected string
	}
	testCases := []test{
		{"file.jpeg", "JPEG"},
		{"/var/foo/bar/test.jpg", "JPEG"},
		{"C:\\Windows\\desktop.PNG", "PNG"},
		{"somefile", ""},
	}

	for _, test := range testCases {
		t.Run(test.fn, func(t *testing.T) {
			parser, err := GetParser(test.fn)
			if test.expected == "" && err == nil {
				t.Fatalf("error expected, but GetParser() suceeded with parser '%v'", parser)
			}
			if test.expected != "" {
				if err != nil {
					t.Fatalf("success expected, got error: %s", err.Error())
				} else if test.expected != parser.Format() {
					t.Fatalf("expected '%s', got '%s'", test.expected, parser.Format())
				}
			}
		})
	}
}
