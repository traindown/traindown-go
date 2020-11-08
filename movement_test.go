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
