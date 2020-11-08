package traindown

import (
	"encoding/json"
	"time"
)

// Session is a collection of Movements that occurred.
type Session struct {
	Date        time.Time   `json:"date"`
	DefaultUnit string      `json:"defaultUnit,omitempty"`
	Errors      []error     `json:"errors"`
	Movements   []*Movement `json:"movements"`

	Metadata Metadata `json:"metadata"`
	Notes    []string `json:"notes"`
}

/* Public */

// NewSession spits out a new Session
func NewSession() *Session {
	return &Session{
		Metadata:  make(Metadata),
		Movements: make([]*Movement, 0),
		Notes:     make([]string, 0),
	}
}

func (s Session) String() string {
	ss, _ := json.Marshal(s)
	return string(ss)
}

// Volumes computes the volume performed by unit.
func (s Session) Volumes() map[string]float32 {
	v := make(map[string]float32)

	for _, m := range s.Movements {
		mvs := m.Volumes()
		for mu, mv := range mvs {
			val, ok := v[mu]
			if ok {
				v[mu] = val + mv
			} else {
				v[mu] = mv
			}
		}
	}

	return v
}

/* Private */

func (s *Session) assignSpecial(k string, v string) bool {
	if isUnit(k) {
		s.DefaultUnit = v
		return true
	}
	return false
}
