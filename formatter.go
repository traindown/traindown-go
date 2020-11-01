package traindown

import (
	"fmt"
	"strings"
)

// Formatter formats Traindown documents
type Formatter struct {
	l *Lexer
}

// NewFormatter returns a pointer to a Formatter
func NewFormatter() (*Formatter, error) {
	lexer, err := NewLexer()

	if err != nil {
		return &Formatter{}, err
	}

	return &Formatter{&lexer}, nil
}

func spacer(inSession bool, inPerformance bool) string {
	var spacer string
	if inSession {
		spacer = ""
	} else if inPerformance {
		spacer = "    "
	} else {
		spacer = "  "
	}
	return spacer
}

// Format takes a Traindown string and returns a prettier version of it.
func (f Formatter) Format(txt string) (string, error) {
	tokens, err := f.l.Scan([]byte(txt))

	if err != nil {
		return "", fmt.Errorf("Failed to format: %q", err)
	}

	var s strings.Builder

	inSession := true
	inPerformance := false

	for _, tok := range tokens {
		switch tok.Name() {
		case "DATE":
			s.WriteString("@ ")
			s.WriteString(tok.Value())
			s.WriteString("\r\n")
		case "FAILS":
			s.WriteString(" ")
			s.WriteString(tok.Value())
			s.WriteString("f")
		case "LOAD":
			inPerformance = true
			s.WriteString("\r\n")
			s.WriteString("  ")
			s.WriteString(tok.Value())
		case "METADATA":
			s.WriteString("\r\n")
			s.WriteString(spacer(inSession, inPerformance))
			s.WriteString("* ")
			s.WriteString(tok.Value())
		case "MOVEMENT", "MOVEMENT_SS":
			inSession = false
			inPerformance = false
			s.WriteString("\r\n\r\n")
			if tok.Value() == "MOVEMENT_SS" {
				s.WriteString("+ ")
			}
			s.WriteString(tok.Value())
			s.WriteString(":")
		case "NOTE":
			s.WriteString("\r\n")
			s.WriteString(spacer(inSession, inPerformance))
			s.WriteString("* ")
			s.WriteString(tok.Value())
		case "REPS":
			s.WriteString(" ")
			s.WriteString(tok.Value())
			s.WriteString("r")
		case "SETS":
			s.WriteString(" ")
			s.WriteString(tok.Value())
			s.WriteString("s")
		}
	}

	return s.String(), nil
}
