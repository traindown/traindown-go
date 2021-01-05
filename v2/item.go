package traindown

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemDate
	itemLoad
	itemFails
	itemMeta
	itemMove
	itemMoveSS
	itemNote
	itemReps
	itemSets
)

type item struct {
	typ itemType
	val string
}

// TODO: Make this not suck
func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	case itemDate:
		return i.val
	case itemLoad:
		return i.val
	case itemFails:
		return i.val
	case itemMeta:
		return i.val
	case itemMove:
		return i.val
	case itemMoveSS:
		return i.val
	case itemNote:
		return i.val
	case itemReps:
		return i.val
	case itemSets:
		return i.val
	}

	return "asdf"
}
