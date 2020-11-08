package traindown

// SpecialAssignable is an interface for things that make use of the IsUnit check.
type specialAssignable interface {
	assignSpecial(string string) bool
}

// IsUnit returns true if the argument is a unit keyword.
func isUnit(u string) bool {
	if u == "u" || u == "U" || u == "unit" || u == "Unit" {
		return true
	}
	return false
}
