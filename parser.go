package traindown

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// Metadata is key value pairs.
type Metadata map[string]interface{}

type specialAssignable interface {
	assignSpecial(string string) bool
}

func isUnit(u string) bool {
	if u == "u" || u == "U" || u == "unit" || u == "Unit" {
		return true
	}
	return false
}

// Performance is an expression of a movement.
type Performance struct {
	Fails        int     `json:"fails"`
	Load         float32 `json:"load"`
	PercentOfMax float32 `json:"percentOfMax,omitempty"`
	Reps         int     `json:"reps"`
	Sequence     int     `json:"sequence"`
	Sets         int     `json:"sets"`
	Unit         string  `json:"unit"`

	Metadata Metadata `json:"metadata"`
	Notes    []string `json:"notes"`
}

// NewPerformance spits out a new Performance
func NewPerformance() *Performance {
	return &Performance{
		Metadata: make(Metadata),
		Notes:    make([]string, 0),
		Reps:     1,
		Sets:     1,
		Unit:     "unknown unit",
	}
}

func (p Performance) String() string {
	ps, _ := json.Marshal(p)
	return string(ps)
}

// Volume produces a float and a string containing the unit.
func (p Performance) Volume() (float32, string) {
	v := (reps - fails) * sets * load
	return v, p.Unit
}

func (p *Performance) assignSpecial(k string, v string) bool {
	if isUnit(k) {
		p.Unit = v
		return true
	}
	return false
}

func (p *Performance) maybeInheritUnit(s *Session, m *Movement) {
	if p.Unit == "unknown unit" {
		if s.DefaultUnit != "" {
			p.Unit = s.DefaultUnit
		}

		if m.DefaultUnit != "" {
			p.Unit = m.DefaultUnit
		}
	}
}

// Movement is an thing you do, you know?
type Movement struct {
	DefaultUnit string `json:"defaultUnit,omitempty"`
	Name        string `json:"name"`
	Sequence    int    `json:"sequence"`
	SuperSet    bool   `json:"superSet"`

	Performances []*Performance `json:"performances"`

	Metadata Metadata `json:"metadata"`
	Notes    []string `json:"notes"`
}

// NewMovement spits out a new Movement
func NewMovement() *Movement {
	return &Movement{
		Metadata:     make(Metadata),
		Notes:        make([]string, 0),
		Performances: make([]*Performance, 0),
	}
}

func (m Movement) String() string {
	ms, _ := json.Marshal(m)
	return string(ms)
}

func (m *Movement) assignSpecial(k string, v string) bool {
	if isUnit(k) {
		m.DefaultUnit = v
		return true
	}
	return false
}

// Session is a collection of Movements that occurred.
type Session struct {
	Date        time.Time   `json:"date"`
	DefaultUnit string      `json:"defaultUnit,omitempty"`
	Errors      []error     `json:"errors"`
	Movements   []*Movement `json:"movements"`

	Metadata Metadata `json:"metadata"`
	Notes    []string `json:"notes"`
}

// NewSession spits out a new Session
func NewSession() *Session {
	return &Session{
		Metadata:  make(Metadata),
		Movements: make([]*Movement, 0),
		Notes:     make([]string, 0),
	}
}

func (s *Session) assignSpecial(k string, v string) bool {
	if isUnit(k) {
		s.DefaultUnit = v
		return true
	}
	return false
}

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
