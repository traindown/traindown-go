package traindown

import (
	"encoding/json"
	"math/rand"
	"testing"
)

func TestNewPerformance(t *testing.T) {
	p := NewPerformance()

	if p.Metadata == nil ||
		p.Notes == nil ||
		p.Reps != 1 ||
		p.Sets != 1 ||
		p.Unit != "unknown unit" {
		t.Fatal("Failed to initialize Performance")
	}
}

func TestStringifyPerformance(t *testing.T) {
	p := NewPerformance()
	expect, _ := json.Marshal(p)

	if p.String() != string(expect) {
		t.Fatal("Failed to stringify Performance")
	}
}

func TestVolume(t *testing.T) {
	fails := rand.Intn(100)
	load := rand.Float32() * 1000
	reps := rand.Intn(100)
	sets := rand.Intn(100)
	unit := "your mom"

	p := Performance{
		Fails: fails,
		Load:  load,
		Reps:  reps,
		Sets:  sets,
		Unit:  unit}

	expected := (float32(reps) - float32(fails)) * float32(sets) * load

	v, u := p.Volume()

	if v != expected {
		t.Errorf("Unexpected volume: %v", v)
	}

	if u != unit {
		t.Errorf("Unexpected unit: %q", u)
	}
}

func TestAssignSpecialToPerformance(t *testing.T) {
	p := NewPerformance()

	assigned := p.assignSpecial("your", "mom")

	if assigned == true {
		t.Errorf("Invalid assignment of special operator")
	}

	keys := []string{"u", "U", "unit", "Unit"}

	for _, k := range keys {
		assigned = p.assignSpecial(k, "your mom")

		if assigned != true || p.Unit != "your mom" {
			t.Errorf("Failed to assign special operator")
		}
	}
}

func TestMaybeInheritUnit(t *testing.T) {
	du := "your mom"
	m := &Movement{DefaultUnit: ""}
	s := &Session{DefaultUnit: du}

	p := NewPerformance()
	p1 := Performance{}

	p.maybeInheritUnit(s, m)
	p1.maybeInheritUnit(s, m)

	if p.Unit != du || p1.Unit != du {
		t.Errorf("Failed to inherit from Session")
	}

	p = NewPerformance()
	p1 = Performance{}

	m.DefaultUnit = du
	s.DefaultUnit = ""

	p.maybeInheritUnit(s, m)
	p1.maybeInheritUnit(s, m)

	if p.Unit != du || p1.Unit != du {
		t.Errorf("Failed to inherit from Movement")
	}

	p.Unit = "not your mom"
	p1.Unit = "not your mom"
	s.DefaultUnit = du

	p.maybeInheritUnit(s, m)
	p1.maybeInheritUnit(s, m)

	if p.Unit == du || p1.Unit == du {
		t.Errorf("Incorrectly overwrote a set unit")
	}
}
