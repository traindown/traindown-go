package traindown

import (
	"fmt"
	"testing"
	"time"
)

func parseCheck(t *testing.T, b []byte, s string) {
	var session *Session
	var err error

	if s == "" {
		session, err = ParseByte(b)
	} else {
		session, err = ParseString(s)
	}

	if err != nil {
		t.Errorf("Failed to parse: %q", err)
	}

	if len(session.Errors) != 0 {
		t.Errorf("Errors on session: %q", session.Errors)
	}

	if session.Date != time.Date(2020, 1, 1, 1, 23, 0, 0, time.UTC) {
		t.Errorf("Failed to parse date")
	}

	m1 := session.Movements[0]
	m1p := m1.Performances
	m1p0 := m1.Performances[0]

	m2 := session.Movements[1]
	m2p := m2.Performances
	m2p0 := m2.Performances[0]
	m2p1 := m2.Performances[1]

	m1p0Meta := Metadata{"performance key": "performance value"}
	m2Meta := Metadata{"movement key": "movement value"}

	if m1.Name != "movement" ||
		len(m1p) != 1 ||
		fmt.Sprint(m1p0.Metadata) != fmt.Sprint(m1p0Meta) ||
		len(m1p0.Notes) != 1 ||
		m1p0.Notes[0] != "performance note" ||
		m1p0.Load != 100 ||
		m1p0.Fails != 1 ||
		m1p0.Reps != 1 ||
		m1p0.Sets != 1 ||
		m1p0.Unit != "movement" {
		t.Errorf("Failed to parse first movement")
	}

	if m2.Name != "another" ||
		m2.SuperSet != true ||
		len(m2p) != 2 ||
		fmt.Sprint(m2.Metadata) != fmt.Sprint(m2Meta) ||
		len(m2.Notes) != 1 ||
		m2.Notes[0] != "movement note" ||
		m2p0.Load != 200.1 ||
		m2p0.Reps != 1 ||
		m2p0.Sets != 1 ||
		m2p0.Unit != "session" ||
		m2p1.Load != 200.2 ||
		m2p1.Reps != 1 ||
		m2p1.Sets != 2 ||
		m2p1.Unit != "performance" {
		t.Errorf("Failed to parse second movement")
	}
}

func TestParse(t *testing.T) {
	text := `
    @ 1/1/20 1:23
    # key: value
    # unit: session

    movement:
      # unit: movement
      100 1r 1f 1s
        * performance note
        # performance key: performance value

    + another:
      * movement note
      # movement key: movement value
      200.1
      200.2 2s
        # unit: performance`
	parseCheck(t, []byte(""), text)
	parseCheck(t, []byte(text), "")
}

func TestParseUnit(t *testing.T) {
	texts := []string{
		"# u: your mom\nmovement:\n100",
		"# U: your mom\nmovement:\n100",
		"# unit: your mom\nmovement:\n100",
		"# Unit: your mom\nmovement:\n100",
		"movement:\n# u: your mom\n100",
		"movement:\n# U: your mom\n100",
		"movement:\n# unit: your mom\n100",
		"movement:\n# Unit: your mom\n100",
		"movement:\n100\n# u: your mom",
		"movement:\n100\n# U: your mom",
		"movement:\n100\n# unit: your mom",
		"movement:\n100\n# Unit: your mom",
	}

	for idx, text := range texts {

		session, err := ParseString(text)

		if err != nil {
			t.Errorf("Failed to parse for %d: %q", idx, err)
			continue
		}

		if len(session.Errors) != 0 {
			t.Errorf("Errors on session for %d: %q", idx, session.Errors)
			continue
		}

		p := session.Movements[0].Performances[0]
		unit := p.Unit

		if unit != "your mom" {
			t.Errorf("Incorrect unit for %d: %q", idx, p)
		}
	}
}
