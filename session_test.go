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

func TestVolumesForSession(t *testing.T) {
	s := NewSession()
	m := NewMovement()
	p1 := &Performance{Unit: "your", Load: 100.0, Reps: 10, Sets: 1}
	p2 := &Performance{Unit: "mom", Load: 500.0, Reps: 1, Sets: 2}
	m.Performances = []*Performance{p1, p2}

	s.Movements = []*Movement{m, m}

	volumes := s.Volumes()

	your, ok := volumes["your"]
	if !ok || your != 2000.0 {
		t.Errorf("Failed to compute 'your'. Expected 2000.0. Got %v", your)
	}

	mom, ok := volumes["mom"]
	if !ok || mom != 2000.0 {
		t.Errorf("Failed to compute 'mom'. Expected 2000.0. Got %v", your)
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
