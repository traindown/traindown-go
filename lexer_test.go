package traindown

import (
	"fmt"
	"testing"
)

type expectation struct {
	name      string
	tokenType int
	value     string
	startLine int
	startCol  int
	endLine   int
	endCol    int
}

func (ex expectation) eq(tok *Token) error {
	if ex.name != tok.Name() {
		return fmt.Errorf("Name mismatch for %q!.\n Expected %q :: got %q", tok, ex.name, tok.Name())
	}

	if ex.tokenType != tok.Type() {
		return fmt.Errorf("Type mismatch for %q!.\n Expected %d :: got %d", tok, ex.tokenType, tok.Type())
	}

	if ex.value != tok.Value() {
		return fmt.Errorf("Value mismatch for %q!.\n Expected %q :: got %q", tok, ex.value, tok.Value())
	}

	sl, sc := tok.Start()
	if ex.startLine != sl {
		return fmt.Errorf("Start Line mismatch for %q!.\n Expected %d :: got %d", tok, ex.startLine, sl)
	}
	if ex.startCol != sc {
		return fmt.Errorf("Start Col mismatch for %q!. Expected %d :: got %d", tok, ex.startCol, sc)
	}

	el, ec := tok.End()
	if ex.endLine != el {
		return fmt.Errorf("End Line mismatch for %q!. Expected %d :: got %d", tok, ex.endLine, el)
	}
	if ex.endCol != ec {
		return fmt.Errorf("End Col mismatch for %q!. Expected %d :: got %d", tok, ex.endCol, ec)
	}

	return nil
}

func TestScan(t *testing.T) {
	lexer, err := NewLexer()

	if err != nil {
		t.Errorf("Failed to init lexer: %q", err.Error())
	}

	text := []byte(`
  @ 2020-01-01 12:34

  * Session Note 123
  # Session Key 123: Session Value 123

  Movement Name 123:
    * Movement Note 123
    # Movement Key 123: Movement Value 123

    100
      * Performance 1 Note
      # Performance 1 Key: Performance 1 Value
    200.2 2r
      * Performance 2 Note
      # Performance 2 Key: Performance 2 Value
    300 3r 3s
      * Performance 3 Note
      # Performance 3 Key: Performance 3 Value
    400.4 4f 4r 4s
      * Performance 4 Note
      # Performance 4 Key: Performance 4 Value

  + Another: 100

  Tabata 20/10: 30s
    # Empty Meta Value:
    # % Effort: 87.2
`)

	var tokens []*Token
	tokens, err = lexer.Scan(text)

	if err != nil {
		t.Errorf("Failed to scan: %q", err.Error())
	}

	expected := []expectation{
		expectation{"DATE", 0, "2020-01-01 12:34", 2, 3, 2, 20},
		expectation{"NOTE", 6, "Session Note 123", 4, 3, 4, 20},
		expectation{
			"METADATA",
			3,
			"Session Key 123: Session Value 123",
			5, 3, 5, 38},
		expectation{"MOVEMENT", 4, "Movement Name 123", 7, 3, 7, 20},
		expectation{"NOTE", 6, "Movement Note 123", 8, 5, 8, 23},
		expectation{
			"METADATA",
			3,
			"Movement Key 123: Movement Value 123",
			9, 5, 9, 42},
		expectation{"LOAD", 1, "100", 11, 5, 11, 7},
		expectation{"NOTE", 6, "Performance 1 Note", 12, 7, 12, 26},
		expectation{
			"METADATA",
			3,
			"Performance 1 Key: Performance 1 Value",
			13, 7, 13, 46},
		expectation{"LOAD", 1, "200.2", 14, 5, 14, 9},
		expectation{"REPS", 7, "2", 14, 11, 14, 12},
		expectation{"NOTE", 6, "Performance 2 Note", 15, 7, 15, 26},
		expectation{
			"METADATA",
			3,
			"Performance 2 Key: Performance 2 Value",
			16, 7, 16, 46},
		expectation{"LOAD", 1, "300", 17, 5, 17, 7},
		expectation{"REPS", 7, "3", 17, 9, 17, 10},
		expectation{"SETS", 8, "3", 17, 12, 17, 13},
		expectation{"NOTE", 6, "Performance 3 Note", 18, 7, 18, 26},
		expectation{
			"METADATA",
			3,
			"Performance 3 Key: Performance 3 Value",
			19, 7, 19, 46},
		expectation{"LOAD", 1, "400.4", 20, 5, 20, 9},
		expectation{"FAILS", 2, "4", 20, 11, 20, 12},
		expectation{"REPS", 7, "4", 20, 14, 20, 15},
		expectation{"SETS", 8, "4", 20, 17, 20, 18},
		expectation{"NOTE", 6, "Performance 4 Note", 21, 7, 21, 26},
		expectation{
			"METADATA",
			3,
			"Performance 4 Key: Performance 4 Value",
			22, 7, 22, 46},
		expectation{"MOVEMENT_SS", 5, "Another", 24, 3, 24, 12},
		expectation{"LOAD", 1, "100", 24, 14, 24, 16},
		expectation{"MOVEMENT", 4, "Tabata 20/10", 26, 3, 26, 15},
		expectation{"SETS", 8, "30", 26, 17, 26, 19},
		expectation{
			"METADATA",
			3,
			"Empty Meta Value: ",
			27, 5, 27, 23},
		expectation{
			"METADATA",
			3,
			"% Effort: 87.2",
			28, 5, 28, 20},
	}

	for idx, ex := range expected {
		err = ex.eq(tokens[idx])

		if err != nil {
			t.Errorf("Mismatch!\n %q", err.Error())
		}
	}
}
