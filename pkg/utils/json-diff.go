package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"sort"
)

type DiffType string

func (d DiffType) Short() string {
	switch d {
	case DiffTypeAdd:
		return "+"
	case DiffTypeRemove:
		return "-"
	case DiffTypeChange:
		return "<>"
	case DiffTypeSame:
		return " "
	default:
		panic(errors.Errorf("invalid DiffType %s", d))
	}
}

const (
	DiffTypeAdd    DiffType = "DiffTypeAdd"
	DiffTypeRemove DiffType = "DiffTypeRemove"
	DiffTypeChange DiffType = "DiffTypeChange"
	DiffTypeSame   DiffType = "DiffTypeSame"
)

type JsonDocumentDiffs struct {
	Elements []*JDiff
}

func (d *JsonDocumentDiffs) Add(e *JDiff) {
	d.Elements = append(d.Elements, e)
}

type JDiff struct {
	Type DiffType
	Path []string
	Old  interface{}
	New  interface{}
}

func DiffJsonValues(a interface{}, b interface{}) *JsonDocumentDiffs {
	diffs := &JsonDocumentDiffs{}
	JsonDiffHelper(a, b, []string{}, diffs)
	return diffs
}

func JsonDiffHelper(a interface{}, b interface{}, pathContext []string, diffs *JsonDocumentDiffs) {
	// make a copy to avoid aliasing
	//path := CopySlice(pathContext)
	//path := append([]string{}, pathContext...) // TODO this doesn't seem to make a deep copy?
	path := make([]string, len(pathContext))
	copy(path, pathContext)

	logrus.Debugf("path: %+v", path)

	if a == nil && b != nil {
		diffs.Add(&JDiff{Type: DiffTypeAdd, Old: a, New: b, Path: path})
	} else if b == nil {
		diffs.Add(&JDiff{Type: DiffTypeRemove, Old: a, New: b, Path: path})
	} else {
		switch aVal := a.(type) {
		case map[string]interface{}:
			switch bVal := b.(type) {
			case map[string]interface{}:
				aKeys := maps.Keys(aVal)
				sort.Strings(aKeys)
				for _, k := range aKeys {
					JsonDiffHelper(aVal[k], bVal[k], append(path, fmt.Sprintf(`%s`, k)), diffs)
				}
				bKeys := maps.Keys(bVal)
				sort.Strings(bKeys)
				for _, k := range bKeys {
					if _, ok := aVal[k]; !ok {
						diffs.Add(&JDiff{Type: DiffTypeAdd, New: bVal[k], Path: append(path, fmt.Sprintf(`%s`, k))})
					}
				}
			default:
				diffs.Add(&JDiff{Type: DiffTypeChange, Old: aVal, New: bVal, Path: path})
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
						diffs.Add(&JDiff{Type: DiffTypeAdd, New: bVal[i], Path: newPath})
					} else if i >= len(bVal) {
						diffs.Add(&JDiff{Type: DiffTypeRemove, Old: aSub, Path: newPath})
					} else {
						JsonDiffHelper(aSub, bVal[i], newPath, diffs)
					}
				}
			default:
				diffs.Add(&JDiff{Type: DiffTypeChange, Old: aVal, New: bVal, Path: path})
			}
		case int:
			switch bVal := b.(type) {
			case int:
				if aVal != bVal {
					diffs.Add(&JDiff{Type: DiffTypeChange, Old: aVal, New: bVal, Path: path})
				}
			default:
				diffs.Add(&JDiff{Type: DiffTypeChange, Old: aVal, New: bVal, Path: path})
			}
		case string:
			switch bVal := b.(type) {
			case string:
				if aVal != bVal {
					diffs.Add(&JDiff{Type: DiffTypeChange, Old: aVal, New: bVal, Path: path})
				}
			default:
				diffs.Add(&JDiff{Type: DiffTypeChange, Old: aVal, New: bVal, Path: path})
			}
		case bool:
			switch bVal := b.(type) {
			case bool:
				if aVal != bVal {
					diffs.Add(&JDiff{Type: DiffTypeChange, Old: aVal, New: bVal, Path: path})
				}
			default:
				diffs.Add(&JDiff{Type: DiffTypeChange, Old: aVal, New: bVal, Path: path})
			}
		//case types.Nil: // TODO is this necessary?
		default:
			panic(errors.Errorf("unrecognized type: %s, %T, %+v", path, aVal, aVal))
		}
	}
}
