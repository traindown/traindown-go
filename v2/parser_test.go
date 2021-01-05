package traindown

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateTime(t *testing.T) {
	expecteds := []string{
		"2020-01-01",
		" 1/1/20 1:23",
		"  Thursday",
	}

	var sb strings.Builder

	for _, e := range expecteds {
		sb.WriteString(fmt.Sprintf("@%v\n", e))
	}

	l := NewLexer(sb.String(), IdleState)

	l.Start()

	for _, e := range expecteds {
		tok, done := l.NextToken()

		assert.Equal(t, false, done, "Failed to finish parsing")
		assert.Equal(t, DateTimeToken, tok.Type, "Incorrect token")
		assert.Equal(t, strings.Trim(e, " "), tok.Value, "Incorrect value")
	}
}

func TestMetadata(t *testing.T) {
	expecteds := []string{
		" Whatever 123!!':     VALUE",
		"  OASKD:123asdf123",
		":",
		" : ",
	}

	var sb strings.Builder

	for _, e := range expecteds {
		sb.WriteString(fmt.Sprintf("#%v\n", e))
	}

	l := NewLexer(sb.String(), IdleState)

	l.Start()

	for _, e := range expecteds {
		kv := strings.Split(e, ":")
		key := kv[0]
		value := kv[1]

		tok, done := l.NextToken()

		assert.Equal(t, false, done, "Failed to finish parsing")
		assert.Equal(t, MetaKeyToken, tok.Type, "Incorrect token")
		assert.Equal(t, strings.Trim(key, " "), tok.Value, "Incorrect value")

		tok, done = l.NextToken()

		assert.Equal(t, false, done, "Failed to finish parsing")
		assert.Equal(t, MetaValueToken, tok.Type, "Incorrect token")
		assert.Equal(t, strings.Trim(value, " "), tok.Value, "Incorrect value")
	}
}

func TestNote(t *testing.T) {
	expecteds := []string{
		"This is 123 it's a baby boy!!!???////...",
		"    UR MOM    ",
	}

	var sb strings.Builder

	for _, e := range expecteds {
		sb.WriteString(fmt.Sprintf("*%v\n", e))
	}

	l := NewLexer(sb.String(), IdleState)

	l.Start()

	for _, e := range expecteds {
		tok, done := l.NextToken()

		assert.Equal(t, false, done, "Failed to finish parsing")
		assert.Equal(t, NoteToken, tok.Type, "Incorrect token")
		// NOTE: We leave a ragged right edge
		assert.Equal(t, strings.TrimLeft(e, " "), tok.Value, "Incorrect value")
	}
}

func TestValue(t *testing.T) {
	src := `
    Squats 123: 123.45 100.5r
    '30 second pullup: bw+25 50r
    '2/3 Squat: 500
    + 1/4 Turds:
      * Hard af
      # Bands: 2 red, 2 blue
      123 100r 5s`

	expecteds := []string{
		"[Movement] Squats 123",
		"[Load] 123.45",
		"[Reps] 100.5",
		"[Movement] 30 second pullup",
		"[Load] bw+25",
		"[Reps] 50",
		"[Movement] 2/3 Squat",
		"[Load] 500",
		"[Supersetted Movement] 1/4 Turds",
		"[Note] Hard af",
		"[Metadata Key] Bands",
		"[Metadata Value] 2 red, 2 blue",
		"[Load] 123",
		"[Reps] 100",
		"[Sets] 5",
	}

	l := NewLexer(src, IdleState)

	l.Start()

	for _, e := range expecteds {
		tok, done := l.NextToken()

		assert.Equal(t, false, done, "Failed to finish parsing")
		assert.Equal(t, e, tok.String(), "Incorrect token")
	}
}
