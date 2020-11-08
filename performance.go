package traindown

import (
	"encoding/json"
)

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

/* Public */

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
	v := (float32(p.Reps) - float32(p.Fails)) * float32(p.Sets) * p.Load
	return v, p.Unit
}

/* Private */

func (p *Performance) assignSpecial(k string, v string) bool {
	if isUnit(k) {
		p.Unit = v
		return true
	}
	return false
}

func (p *Performance) maybeInheritUnit(s *Session, m *Movement) {
	if p.Unit == "unknown unit" || p.Unit == "" {
		if s.DefaultUnit != "" {
			p.Unit = s.DefaultUnit
		}

		if m.DefaultUnit != "" {
			p.Unit = m.DefaultUnit
		}
	}
}
