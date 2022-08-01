package diff

import "github.com/pkg/errors"

type Kind string

func (d Kind) Short() string {
	switch d {
	case KindAdd:
		return "+"
	case KindRemove:
		return "-"
	case KindChange:
		return "<>"
	case KindSame:
		return " "
	default:
		panic(errors.Errorf("invalid Kind %s", d))
	}
}

const (
	KindAdd    Kind = "KindAdd"
	KindRemove Kind = "KindRemove"
	KindChange Kind = "KindChange"
	KindSame   Kind = "KindSame"
)

type Node struct {
	Kind Kind
	Path []string
	Old  interface{}
	New  interface{}
}
