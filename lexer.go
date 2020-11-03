package traindown

import (
	"fmt"
	"log"
	"strings"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

// Tokens used in parsing Traindown inputs
var Tokens = []string{
	"DATE", "LOAD", "FAILS", "METADATA", "MOVEMENT", "MOVEMENT_SS", "NOTE", "REPS", "SETS",
}

// Token holds information about a token
type Token struct {
	t *lexmachine.Token
}

// TokenMap maps is an enum
var TokenMap map[string]int

func init() {
	TokenMap = make(map[string]int)
	for i, n := range Tokens {
		TokenMap[n] = i
	}
}

// Name getter
func (tok Token) Name() string {
	return Tokens[tok.t.Type]
}

// Type getter
func (tok Token) Type() int {
	return tok.t.Type
}

// Value getter
func (tok Token) Value() string {
	return tok.t.Value.(string)
}

// Start getter
func (tok Token) Start() (int, int) {
	return tok.t.StartLine, tok.t.StartColumn
}

// End getter
func (tok Token) End() (int, int) {
	return tok.t.EndLine, tok.t.EndColumn
}

// String override
func (tok *Token) String() string {
	return fmt.Sprintf("%q %q (From: r%d, c%d To: r%d c%d)", tok.Name(), tok.t.Value, tok.t.StartLine, tok.t.StartColumn, tok.t.EndLine, tok.t.EndColumn)
}

// Lexer type
type Lexer struct {
	l *lexmachine.Lexer
}

// NewLexer returns a new Lexer
func NewLexer() (Lexer, error) {
	var lexer = lexmachine.NewLexer()

	lexer.Add(
		[]byte("@[^\n]*"),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			return scan.Token(
					TokenMap["DATE"],
					strings.TrimSpace(string(match.Bytes)[1:]),
					match),
				nil
		},
	)
	lexer.Add(
		[]byte(`#[^\n|\r][a-zA-Z| |:|[0-9]+`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			lR := strings.Split(string(match.Bytes)[1:], ":")
			l := strings.TrimSpace(lR[0])
			r := strings.TrimSpace(lR[1])
			var kvp strings.Builder
			kvp.WriteString(l)
			kvp.WriteString(": ")
			kvp.WriteString(r)
			return scan.Token(
					TokenMap["METADATA"],
					kvp.String(),
					match),
				nil
		},
	)
	lexer.Add(
		[]byte(`\*[^\n|\r]*`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			return scan.Token(
					TokenMap["NOTE"],
					strings.TrimSpace(string(match.Bytes)[1:]),
					match),
				nil
		},
	)
	lexer.Add(
		[]byte(`[0-9]+[fF]`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			s := string(match.Bytes)
			return scan.Token(
					TokenMap["FAILS"],
					s[:len(s)-1],
					match),
				nil
		},
	)
	lexer.Add(
		[]byte(`[0-9]+[rR]`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			s := string(match.Bytes)
			return scan.Token(
					TokenMap["REPS"],
					s[:len(s)-1],
					match),
				nil
		},
	)
	lexer.Add(
		[]byte(`[0-9]+[sS]`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			s := string(match.Bytes)
			return scan.Token(
					TokenMap["SETS"],
					s[:len(s)-1],
					match),
				nil
		},
	)
	lexer.Add(
		[]byte(`[0-9]*\.?[0-9]+`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			return scan.Token(TokenMap["LOAD"], string(match.Bytes), match), nil
		},
	)
	lexer.Add(
		[]byte(`((\+\s*?)?\w+\s?)+:`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			s := strings.TrimSuffix(string(match.Bytes), ":")

			var tokType int
			if strings.HasPrefix(s, "+") {
				tokType = TokenMap["MOVEMENT_SS"]
				s = strings.TrimPrefix(s, "+")
				s = strings.TrimSpace(s)
			} else {
				tokType = TokenMap["MOVEMENT"]
			}

			return scan.Token(tokType, s, match), nil
		},
	)
	lexer.Add(
		[]byte("( |\t|\n|\r)"),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			return nil, nil
		},
	)

	err := lexer.CompileDFA()
	if err != nil {
		return Lexer{}, err
	}

	return Lexer{lexer}, nil
}

// Scan returns the next token
func (lexer Lexer) Scan(text []byte) ([]*Token, error) {
	scanner, err := lexer.l.Scanner(text)

	var tokens []*Token

	if err != nil {
		return tokens, err
	}

	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {

		if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
			scanner.TC = ui.FailTC
			log.Printf("skipping %v", ui)
		} else if err != nil {
			return tokens, err
		} else {
			token := &Token{tok.(*lexmachine.Token)}
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}
