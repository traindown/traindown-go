package traindown

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// EOFRune is the end of file rune.
const EOFRune rune = -1

// StateFunc takes in a pointer to Lexer and returns a new StateFunc.
type StateFunc func(*Lexer) StateFunc

// Lexer is the core struct that manages state through the StateFuncs.
type Lexer struct {
	bufferSize      int
	Err             error
	ErrorHandler    func(e string)
	rewind          vcr
	source          string
	start, position int
	startState      StateFunc
	tokens          chan Token
}

// NewLexer creates a returns a Lexer.
func NewLexer(src string, start StateFunc) *Lexer {
	return &Lexer{
		bufferSize: 10,
		source:     src,
		startState: start,
		start:      0,
		position:   0,
		rewind:     newVCR(),
	}
}

// Current returns the value being being analyzed at this moment.
func (l Lexer) Current() string {
	return l.source[l.start:l.position]
}

// Emit pushes into the tokens channel a new token with the Current value.
func (l *Lexer) Emit(t TokenType) {
	tok := Token{
		Type:  t,
		Value: l.Current(),
	}

	l.tokens <- tok

	l.start = l.position

	l.rewind.clear()
}

// Error provides a mechanism for parsers to signal to the Lexer an error has
// occurred.
func (l *Lexer) Error(e string) {
	if l.ErrorHandler != nil {
		l.Err = errors.New(e)
		l.ErrorHandler(e)
	} else {
		panic(e)
	}
}

// Ignore dumps the VCR and sets the current beginning position to the cursor.
func (l *Lexer) Ignore() {
	l.rewind.clear()

	l.start = l.position
}

// Next increments to the next rune and returns the value.
func (l *Lexer) Next() rune {
	var r rune
	var s int

	str := l.source[l.position:]

	if len(str) == 0 {
		r, s = EOFRune, 0
	} else {
		r, s = utf8.DecodeRuneInString(str)
	}

	l.position += s
	l.rewind.push(r)

	return r
}

// NextToken returns the next token from the lexer along with a done bool.
func (l *Lexer) NextToken() (*Token, bool) {
	if tok, ok := <-l.tokens; ok {
		return &tok, false
	}

	return nil, true
}

// Peek is Next + Rewind returning the rune that will come next.
func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Rewind()

	return r
}

// Rewind pops the last rune from the VCR and sets position there guarding
// against rewinding into a previous token.
func (l *Lexer) Rewind() {
	r := l.rewind.pop()

	if r > EOFRune {
		size := utf8.RuneLen(r)
		l.position -= size

		if l.position < l.start {
			l.position = l.start
		}
	}
}

// SetBufferSize is a setter for the buffer capacity.
func (l *Lexer) SetBufferSize(s int) {
	l.bufferSize = s
}

// Start starts the Lexer which populates the tokens channel.
func (l *Lexer) Start() {
	l.tokens = make(chan Token, l.bufferSize)

	state := l.startState

	go func(l *Lexer) {
		for state != nil {
			state = state(l)
		}

		close(l.tokens)
	}(l)
}

// Take accepts a string of runes that it will match against to increment
// forward only when the current rune is in the match set.
func (l *Lexer) Take(chars string) {
	r := l.Next()

	for strings.ContainsRune(chars, r) {
		r = l.Next()
	}

	l.Rewind() // last next wasn't a match
}
