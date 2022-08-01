package diff

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"sort"
)

type JsonDiff struct {
	Changes []*Node
}

func (d *JsonDiff) Add(e *Node) {
	d.Changes = append(d.Changes, e)
}

func CompareJson(a interface{}, b interface{}) *JsonDiff {
	diffs := &JsonDiff{}
	CompareJsonHelper(a, b, []string{}, diffs)
	return diffs
}

func CompareJsonHelper(a interface{}, b interface{}, pathContext []string, diffs *JsonDiff) {
	// make a copy to avoid aliasing
	//path := CopySlice(pathContext)
	//path := append([]string{}, pathContext...) // TODO this doesn't seem to make a deep copy?
	path := make([]string, len(pathContext))
	copy(path, pathContext)

	logrus.Debugf("path: %+v", path)

	if a == nil && b != nil {
		diffs.Add(&Node{Kind: KindAdd, Old: a, New: b, Path: path})
	} else if b == nil {
		diffs.Add(&Node{Kind: KindRemove, Old: a, New: b, Path: path})
	} else {
		switch aVal := a.(type) {
		case map[string]interface{}:
			switch bVal := b.(type) {
			case map[string]interface{}:
				aKeys := maps.Keys(aVal)
				sort.Strings(aKeys)
				for _, k := range aKeys {
					CompareJsonHelper(aVal[k], bVal[k], append(path, fmt.Sprintf(`%s`, k)), diffs)
				}
				bKeys := maps.Keys(bVal)
				sort.Strings(bKeys)
				for _, k := range bKeys {
					if _, ok := aVal[k]; !ok {
						diffs.Add(&Node{Kind: KindAdd, New: bVal[k], Path: append(path, fmt.Sprintf(`%s`, k))})
					}
				}
			default:
				diffs.Add(&Node{Kind: KindChange, Old: aVal, New: bVal, Path: path})
			}
		case []interface{}:
			switch bVal := b.(type) {
			case []interface{}:
				minLength := len(aVal)
				if len(bVal) < minLength {
					minLength = len(bVal)
				}
				for i, aSub := range aVal {
					newPath := append(path, fmt.Sprintf("%d", i))
					if i >= len(aVal) {
						diffs.Add(&Node{Kind: KindAdd, New: bVal[i], Path: newPath})
					} else if i >= len(bVal) {
						diffs.Add(&Node{Kind: KindRemove, Old: aSub, Path: newPath})
					} else {
						CompareJsonHelper(aSub, bVal[i], newPath, diffs)
					}
				}
			default:
				diffs.Add(&Node{Kind: KindChange, Old: aVal, New: bVal, Path: path})
			}
		case int:
			switch bVal := b.(type) {
			case int:
				if aVal != bVal {
					diffs.Add(&Node{Kind: KindChange, Old: aVal, New: bVal, Path: path})
				}
			default:
				diffs.Add(&Node{Kind: KindChange, Old: aVal, New: bVal, Path: path})
			}
		case string:
			switch bVal := b.(type) {
			case string:
				if aVal != bVal {
					diffs.Add(&Node{Kind: KindChange, Old: aVal, New: bVal, Path: path})
				}
			default:
				diffs.Add(&Node{Kind: KindChange, Old: aVal, New: bVal, Path: path})
			}
		case bool:
			switch bVal := b.(type) {
			case bool:
				if aVal != bVal {
					diffs.Add(&Node{Kind: KindChange, Old: aVal, New: bVal, Path: path})
				}
			default:
				diffs.Add(&Node{Kind: KindChange, Old: aVal, New: bVal, Path: path})
			}
		//case types.Nil: // TODO is this necessary?
		default:
			panic(errors.Errorf("unrecognized type: %s, %T, %+v", path, aVal, aVal))
		}
	}
}
