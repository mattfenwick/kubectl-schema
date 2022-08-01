package diff

import "github.com/mattfenwick/kubectl-schema/pkg/utils"

func SliceDiff(as []string, bs []string) *MapDiff {
	aSet, bSet := utils.Set(as), utils.Set(bs)
	var added, removed, same []string
	for key := range aSet {
		if bSet[key] {
			same = append(same, key)
		} else {
			removed = append(removed, key)
		}
	}
	for key := range bSet {
		if !aSet[key] {
			added = append(added, key)
		}
	}
	return &MapDiff{
		Added:   added,
		Removed: removed,
		Same:    same,
	}
}
