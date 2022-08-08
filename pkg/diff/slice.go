package diff

import (
	"github.com/mattfenwick/collections/pkg/set"
)

func SliceDiff(as []string, bs []string) *MapDiff {
	aSet, bSet := set.FromSlice(as), set.FromSlice(bs)
	return &MapDiff{
		Added:   bSet.Difference(aSet).ToSlice(), // b - a: not in a, but is now in b
		Removed: aSet.Difference(bSet).ToSlice(), // a - b: was in a, no longer in b
		Same:    aSet.Intersect(bSet).ToSlice(),  // a & b: in both a and b
	}
}
