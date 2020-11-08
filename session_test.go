package traindown

import (
	"encoding/json"
	"testing"
)

func TestNewSession(t *testing.T) {
	s := NewSession()

	if s.Metadata == nil ||
		s.Notes == nil ||
		s.Movements == nil {
		t.Fatal("Failed to initialize Movement")
	}
}

func TestStringifySession(t *testing.T) {
	m := NewMovement()
	expect, _ := json.Marshal(m)

	if m.String() != string(expect) {
		t.Fatal("Failed to stringify Movement")
	}
}

func TestAssignSpecialToSession(t *testing.T) {
	s := NewSession()

	assigned := s.assignSpecial("your", "mom")

	if assigned == true {
		t.Errorf("Invalid assignment of special operator")
	}

	keys := []string{"u", "U", "unit", "Unit"}

	for _, k := range keys {
		assigned = s.assignSpecial(k, "your mom")

		if assigned != true || s.DefaultUnit != "your mom" {
			t.Errorf("Failed to assign special operator")
		}
	}
}
