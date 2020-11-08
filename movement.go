package traindown

import (
	"encoding/json"
)

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

/* Public */

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

// Volumes computes the volume performed by unit.
func (m Movement) Volumes() map[string]float32 {
	v := make(map[string]float32)
	for _, p := range m.Performances {
		pv, pu := p.Volume()

		val, ok := v[pu]
		if ok {
			v[pu] = val + pv
		} else {
			v[pu] = pv
		}
	}

	return v
}

/* Private */

func (m *Movement) assignSpecial(k string, v string) bool {
	if isUnit(k) {
		m.DefaultUnit = v
		return true
	}
	return false
}
