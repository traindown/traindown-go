package traindown

import (
	"encoding/json"
	"testing"
)

func TestNewMovement(t *testing.T) {
	m := NewMovement()

	if m.Metadata == nil ||
		m.Notes == nil ||
		m.Performances == nil {
		t.Fatal("Failed to initialize Movement")
	}
}

func TestStringifyMovement(t *testing.T) {
	m := NewMovement()
	expect, _ := json.Marshal(m)

	if m.String() != string(expect) {
		t.Fatal("Failed to stringify Movement")
	}
}

func TestVolumesForMovement(t *testing.T) {
	m := NewMovement()
	p1 := &Performance{Unit: "your", Load: 100.0, Reps: 10, Sets: 1}
	p2 := &Performance{Unit: "mom", Load: 500.0, Reps: 1, Sets: 2}
	m.Performances = []*Performance{p1, p2}

	volumes := m.Volumes()

	your, ok := volumes["your"]
	if !ok || your != 1000.0 {
		t.Errorf("Failed to compute 'your'. Expected 1000.0. Got %v", your)
	}

	mom, ok := volumes["mom"]
	if !ok || mom != 1000.0 {
		t.Errorf("Failed to compute 'mom'. Expected 1000.0. Got %v", your)
	}
}

func TestAssignSpecialToMovement(t *testing.T) {
	m := NewMovement()

	assigned := m.assignSpecial("your", "mom")

	if assigned == true {
		t.Errorf("Invalid assignment of special operator")
	}

	keys := []string{"u", "U", "unit", "Unit"}

	for _, k := range keys {
		assigned = m.assignSpecial(k, "your mom")

		if assigned != true || m.DefaultUnit != "your mom" {
			t.Errorf("Failed to assign special operator")
		}
	}
}
