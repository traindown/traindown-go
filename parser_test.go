package traindown

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
		t.Fatal("Failed to parse date")
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
		m1.Sequence != 0 ||
		len(m1p) != 1 ||
		fmt.Sprint(m1p0.Metadata) != fmt.Sprint(m1p0Meta) ||
		m1p0.Sequence != 0 ||
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
		m2.Sequence != 1 ||
		m2.SuperSet != true ||
		len(m2p) != 2 ||
		fmt.Sprint(m2.Metadata) != fmt.Sprint(m2Meta) ||
		len(m2.Notes) != 1 ||
		m2.Notes[0] != "movement note" ||
		m2p0.Sequence != 0 ||
		m2p0.Load != 200.1 ||
		m2p0.Reps != 1 ||
		m2p0.Sets != 1 ||
		m2p0.Unit != "session" ||
		m2p1.Sequence != 1 ||
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

func TestParseFile(t *testing.T) {
	text, err := ioutil.ReadFile("./testdata")

	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	fmt.Println(text)

	parseCheck(t, text, "")
	parseCheck(t, []byte(""), string(text))
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

func TestParseZeroLoad(t *testing.T) {
	text := `@2022-01-08

Tricep dips: 0 8r 3s

Pull ups: 0 5r 5s
`
	expectedSession := NewSession()
	expectedSession.Date = time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)
	m1 := NewMovement()
	m1.Name = "Tricep dips"
	m1.Sequence = 0
	m1p1 := NewPerformance()
	m1p1.Load = 0
	m1p1.Reps = 8
	m1p1.Sets = 3
	m1.Performances = []*Performance{m1p1}
	m2 := NewMovement()
	m2.Name = "Pull ups"
	m2.Sequence = 1
	m2p1 := NewPerformance()
	m2p1.Load = 0
	m2p1.Reps = 5
	m2p1.Sets = 5
	m2.Performances = []*Performance{m2p1}
	expectedSession.Movements = []*Movement{m1, m2}

	session, err := ParseString(text)
	if err != nil {
		t.Errorf("Failed to parse session: %q", err)
	}
	if diff := cmp.Diff(expectedSession, session); diff != "" {
		t.Errorf("ParseString() mismatch (-want, +got):\n%s", diff)
	}
}
