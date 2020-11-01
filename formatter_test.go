package traindown

import (
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	fmtr, err := NewFormatter()

	if err != nil {
		t.Errorf("Failed to init formatter")
	}

	expected := "@ 2020-01-01 1:23\r\n\r\n* key: value\r\n\r\nmovement:\r\n  100 1r 1f 1s\r\n    * performance note\r\n    * performance key: performance value\r\n\r\nanother:\r\n  * movement note\r\n  * movement key: movement value\r\n  200 2s"

	text := `
    @ 2020-01-01 1:23
# key: value

    movement: 100
    1r 1f
             1s
    * performance note
            # performance key:     performance value
    + another:
    * movement note
    # movement key: movement value
    200
  2s`

	var res string
	res, err = fmtr.Format(text)

	if err != nil {
		t.Errorf("Failed formatting: %q", err)
	}

	if res != strings.Trim(expected, " ") {
		t.Errorf("Output mismatch:\n\nGot:\n%q\n\nExpected:\n%q", res, expected)
	}
}
