package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

func CompareResolvedResources(a *ResolvedType, b *ResolvedType) *utils.JsonDocumentDiffs {
	diffs := &utils.JsonDocumentDiffs{}
	CompareResolvedResourcesHelper(a, b, []string{}, diffs)
	return diffs
}

func CompareResolvedResourcesHelper(a *ResolvedType, b *ResolvedType, pathContext []string, diffs *utils.JsonDocumentDiffs) {
	// make a copy to avoid aliasing
	path := utils.CopySlice(pathContext)

	logrus.Debugf("path: %+v", path)

	if a == nil && b != nil {
		diffs.Add(&utils.JDiff{Type: utils.DiffTypeAdd, Old: a, New: b, Path: path})
	} else if b == nil {
		diffs.Add(&utils.JDiff{Type: utils.DiffTypeRemove, Old: a, New: b, Path: path})
	} else {
		if a.Empty {
			if !b.Empty {
				diffs.Add(&utils.JDiff{Type: utils.DiffTypeChange, Old: a, New: b, Path: path})
			}
		} else if a.Primitive != "" {
			if a.Primitive != b.Primitive {
				diffs.Add(&utils.JDiff{Type: utils.DiffTypeChange, Old: a, New: b, Path: path})
			}
		} else if a.Array != nil {
			if b.Array != nil {
				CompareResolvedResourcesHelper(a.Array, b.Array, append(path, "[]"), diffs)
			} else {
				diffs.Add(&utils.JDiff{Type: utils.DiffTypeChange, Old: a, New: b, Path: path})
			}
		} else if a.Object != nil {
			if b.Object != nil {
				for _, k := range slice.Sort(maps.Keys(a.Object.Properties)) {
					CompareResolvedResourcesHelper(a.Object.Properties[k], b.Object.Properties[k], append(path, fmt.Sprintf(`%s`, k)), diffs)
				}
				for _, k := range slice.Sort(maps.Keys(b.Object.Properties)) {
					if _, ok := a.Object.Properties[k]; !ok {
						diffs.Add(&utils.JDiff{Type: utils.DiffTypeAdd, New: b.Object.Properties[k], Path: append(path, fmt.Sprintf(`%s`, k))})
					}
				}
				// TODO
				//   compare 'required' fields:
				//minLength := len(aVal.Required)
				//if len(bVal.Required) < minLength {
				//	minLength = len(bVal.Required)
				//}
				//for i, aSub := range aVal.Required {
				//	newPath := append(utils.CopySlice(path), "required", fmt.Sprintf("%d", i))
				//	if i >= len(aVal.Required) {
				//		diffs.Add(&utils.JDiff{Type: utils.DiffTypeAdd, New: bVal.Required[i], Path: newPath})
				//	} else if i >= len(bVal.Required) {
				//		diffs.Add(&utils.JDiff{Type: utils.DiffTypeRemove, Old: aSub, Path: newPath})
				//	} else if aSub != bVal.Required[i] {
				//		diffs.Add(&utils.JDiff{Type: utils.DiffTypeChange, Old: aSub, New: bVal.Required[i], Path: newPath})
				//	}
			} else {
				diffs.Add(&utils.JDiff{Type: utils.DiffTypeChange, Old: a, New: b, Path: path})
			}
		} else if a.Circular != "" {
			if a.Circular != b.Circular {
				diffs.Add(&utils.JDiff{Type: utils.DiffTypeChange, Old: a, New: b, Path: path})
			}
		} else {
			panic(errors.Errorf("invalid ResolvedType value: %+v", a))
		}
	}
}
