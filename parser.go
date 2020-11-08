package traindown

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// ParseByte takes in a Traindown byte slice and returns a pointer to a Session.
func ParseByte(txt []byte) (*Session, error) {
	s, err := parse("", txt)

	if err != nil {
		return &Session{}, err
	}

	return s, nil
}

// ParseString takes in a Traindown string and returns a pointer to a Session.
func ParseString(txt string) (*Session, error) {
	s, err := parse(txt, []byte(""))

	if err != nil {
		return &Session{}, err
	}

	return s, nil
}

func floatValue(s string, t string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)

	if err != nil {
		return 0.0, fmt.Errorf("Failed to parse %q: %q", t, s)
	}

	return float32(f), nil
}

func intValue(s string, t string) (int, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse %q: %q", t, s)
	}

	return i, nil
}

func parse(str string, b []byte) (*Session, error) {
	s := NewSession()

	lexer, err := NewLexer()

	if err != nil {
		return s, err
	}

	var tokens []*Token
	if str != "" {
		tokens, err = lexer.Scan([]byte(str))
	} else {
		tokens, err = lexer.Scan(b)
	}

	if err != nil {
		return s, err
	}

	m := NewMovement()
	p := NewPerformance()
	inSession := true
	inPerformance := false
	mSeq := 0
	pSeq := 0

	for _, tok := range tokens {
		switch tok.Name() {
		case "DATE":
			d, err := dateparse.ParseAny(tok.Value())

			if err != nil {
				s.Errors = append(s.Errors, fmt.Errorf("Failed to parse date: %q. Using today UTC", err))
				s.Date = time.Now()
			} else {
				s.Date = d
			}
		case "FAILS":
			i, err := intValue(tok.Value(), "fails")

			if err != nil {
				s.Errors = append(s.Errors, err)
			}

			p.Fails = i
		case "LOAD":
			if inPerformance {
				p.Sequence = pSeq
				p.maybeInheritUnit(s, m)
				m.Performances = append(m.Performances, p)
				p = NewPerformance()
				pSeq++
			}
			f, err := floatValue(tok.Value(), "load")

			if err != nil {
				s.Errors = append(s.Errors, err)
			}

			p.Load = f
			inPerformance = true
		case "METADATA":
			pair := strings.Split(tok.Value(), ":")
			key := strings.Trim(pair[0], " ")
			value := strings.Trim(pair[1], " ")

			if inSession {
				if !s.assignSpecial(key, value) {
					s.Metadata[key] = value
				}
			} else if inPerformance {
				if !p.assignSpecial(key, value) {
					p.Metadata[key] = value
				}
			} else {
				if !m.assignSpecial(key, value) {
					m.Metadata[key] = value
				}
			}
		case "MOVEMENT", "MOVEMENT_SS":
			inSession = false

			if inPerformance {
				p.Sequence = pSeq
				p.maybeInheritUnit(s, m)
				m.Performances = append(m.Performances, p)
				p = NewPerformance()
				pSeq++
			}
			inPerformance = false

			if m.Name != "" {
				m.Sequence = mSeq
				s.Movements = append(s.Movements, m)
				m = NewMovement()
				mSeq++
				pSeq = 0
			}

			m.Name = tok.Value()

			if tok.Name() == "MOVEMENT_SS" {
				m.SuperSet = true
			}
		case "NOTE":
			if inSession {
				s.Notes = append(s.Notes, tok.Value())
			} else if inPerformance {
				p.Notes = append(p.Notes, tok.Value())
			} else {
				m.Notes = append(m.Notes, tok.Value())
			}
		case "REPS":
			i, err := intValue(tok.Value(), "reps")

			if err != nil {
				s.Errors = append(s.Errors, err)
			}

			p.Reps = i
		case "SETS":
			i, err := intValue(tok.Value(), "sets")

			if err != nil {
				s.Errors = append(s.Errors, err)
			}

			p.Sets = i
		}
	}

	if p.Load != 0.0 {
		p.Sequence = pSeq
		p.maybeInheritUnit(s, m)
		m.Performances = append(m.Performances, p)
	}

	if m.Name != "" {
		m.Sequence = mSeq
		s.Movements = append(s.Movements, m)
	}

	return s, nil
}
