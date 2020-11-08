package traindown

import (
	"testing"
)

func TestIsUnit(t *testing.T) {
	us := []string{"u", "U", "unit", "Unit"}

	for _, u := range us {
		if !isUnit(u) {
			t.Errorf("Incorrect response for %v", u)
		}
	}
}
