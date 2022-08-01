package diff

import "sort"

type MapDiff struct {
	Added   []string
	Removed []string
	Same    []string
}

func (m *MapDiff) Sort() {
	sort.Strings(m.Added)
	sort.Strings(m.Removed)
	sort.Strings(m.Same)
}
