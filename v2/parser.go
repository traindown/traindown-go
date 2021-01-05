package traindown

import (
	"fmt"
	"unicode"
)

const (
	DateTimeToken TokenType = iota
	FailToken
	LoadToken
	MetaKeyToken
	MetaValueToken
	MovementToken
	NoteToken
	RepToken
	SetToken
	SupersetMovementToken
)

func (tt TokenType) String() string {
	s := "??"

	switch tt {
	case DateTimeToken:
		s = "Date / Time"
	case FailToken:
		s = "Fails"
	case LoadToken:
		s = "Load"
	case MetaKeyToken:
		s = "Metadata Key"
	case MetaValueToken:
		s = "Metadata Value"
	case MovementToken:
		s = "Movement"
	case NoteToken:
		s = "Note"
	case RepToken:
		s = "Reps"
	case SetToken:
		s = "Sets"
	case SupersetMovementToken:
		s = "Supersetted Movement"
	}

	return s
}

func IdleState(l *Lexer) StateFunc {
	r := l.Peek()

	if r == EOFRune {
		return nil
	}

	if isWhitespace(r) {
		return WhitespaceState
	}

	switch r {
	case '@':
		return DateTimeState
	case '#':
		return MetaKeyState
	case '*':
		return NoteState
	default:
		return ValueState
	}
}

func DateTimeState(l *Lexer) StateFunc {
	l.Take("@ ")
	l.Ignore()

	r := l.Next()

	for !isLineTerminator(r) {
		r = l.Next()
	}

	l.Rewind()
	l.Emit(DateTimeToken)

	return IdleState
}

func MetaKeyState(l *Lexer) StateFunc {
	l.Take("# ")
	l.Ignore()

	r := l.Next()

	for r != ':' {
		r = l.Next()
	}

	l.Rewind()
	l.Emit(MetaKeyToken)

	return MetaValueState
}

func MetaValueState(l *Lexer) StateFunc {
	l.Take(": ")
	l.Ignore()

	r := l.Next()

	for !isLineTerminator(r) {
		r = l.Next()
	}

	l.Rewind()
	l.Emit(MetaValueToken)

	return IdleState
}

func MovementState(l *Lexer) StateFunc {
	superset := false

	r := l.Next()

	if r == '+' {
		superset = true
		l.Take(" ")
		l.Ignore()
		r = l.Next()
	}

	if r == '\'' {
		l.Ignore()
		r = l.Next()
	}

	for r != ':' {
		r = l.Next()
	}

	l.Rewind()

	if superset {
		l.Emit(SupersetMovementToken)
	} else {
		l.Emit(MovementToken)
	}

	l.Take(":")
	l.Ignore()

	return IdleState
}

func NoteState(l *Lexer) StateFunc {
	l.Take("* ")
	l.Ignore()

	r := l.Next()

	for !isLineTerminator(r) {
		r = l.Next()
	}

	l.Rewind()
	l.Emit(NoteToken)

	return IdleState
}

func NumberState(l *Lexer) StateFunc {
	l.Take("0123456789.")

	switch l.Peek() {
	case 'f', 'F':
		l.Emit(FailToken)
	case 'r', 'R':
		l.Emit(RepToken)
	case 's', 'S':
		l.Emit(SetToken)
	default:
		l.Emit(LoadToken)
	}

	l.Take("fFrRsS ")
	l.Ignore()

	return IdleState
}

// TODO: Probe for movements starting with a number as well as a load like "bw"
func ValueState(l *Lexer) StateFunc {
	r := l.Next()

	// NOTE: Definitely a super setted movement or a movement name beginning with
	// an escaped number.
	if r == '+' || r == '\'' {
		l.Rewind()
		return MovementState
	}

	if unicode.IsLetter(r) {
		// NOTE: Definitely not a bodyweight load
		if r != 'b' && r != 'B' {
			l.Rewind()
			return MovementState
		}

		p := l.Peek()

		if p != 'w' && p != 'W' {
			l.Rewind()
			return MovementState
		}

		// NOTE: We have a bodyweight load
		for !isWhitespace(r) {
			r = l.Next()
		}

		l.Rewind()
		l.Emit(LoadToken)

		return IdleState
	}

	return NumberState
}

func WhitespaceState(l *Lexer) StateFunc {
	r := l.Next()

	if r == EOFRune {
		return nil
	}

	if !isWhitespace(r) {
		l.Error(fmt.Sprintf("unexpected token %q", r))
		return nil
	}

	for isWhitespace(r) {
		r = l.Next()
	}

	l.Rewind()
	l.Ignore()

	return IdleState
}

func isLineTerminator(r rune) bool {
	if r == EOFRune || r == ';' || r == '\n' || r == '\r' {
		return true
	}

	return false
}

func isWhitespace(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}

	return false
}
