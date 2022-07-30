package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/kubectl-schema/pkg/swagger/apiversions"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

type ExplainGVKGroupBy string

const (
	ExplainGVKGroupByResource   = "ExplainGVKGroupByResource"
	ExplainGVKGroupByApiVersion = "ExplainGVKGroupByApiVersion"
)

func (e ExplainGVKGroupBy) Header() string {
	switch e {
	case ExplainGVKGroupByResource:
		return "Resource"
	case ExplainGVKGroupByApiVersion:
		return "API version"
	default:
		panic(errors.Errorf("invalid groupBy: %s", e))
	}
}

func ExplainGvks(groupBy ExplainGVKGroupBy, versions []string, include func(string, string) bool, calculateDiff bool) string {
	table := NewPivotTable(groupBy.Header(), versions)
	for _, version := range versions {
		kubeVersion := MustVersion(version)
		logrus.Debugf("kube version: %s", version)

		spec := MustReadSwaggerSpecFromGithub(kubeVersion)
		for name, def := range spec.Definitions {
			if len(def.XKubernetesGroupVersionKind) > 0 {
				logrus.Debugf("%s, %s, %+v\n", name, def.Type, def.XKubernetesGroupVersionKind)
			}
			for _, gvk := range def.XKubernetesGroupVersionKind {
				apiVersion := gvk.GroupVersion()
				if include(apiVersion, gvk.Kind) {
					logrus.Debugf("adding gvk: %s, %s", apiVersion, gvk.Kind)
					switch groupBy {
					case ExplainGVKGroupByResource:
						table.Add(gvk.Kind, kubeVersion.ToString(), apiVersion)
					case ExplainGVKGroupByApiVersion:
						table.Add(apiVersion, kubeVersion.ToString(), gvk.Kind)
					default:
						panic(errors.Errorf("invalid groupBy: %s", groupBy))
					}
				} else {
					logrus.Debugf("skipping gvk: %s, %s", apiVersion, gvk.Kind)
				}
			}
		}
	}

	if calculateDiff {
		return table.FormattedTable(func(rowKey string, values [][]string) []string {
			if len(values) == 0 {
				panic(errors.Errorf("unable to calculate diff for 0 versions"))
			}
			prev := values[0]
			row := []string{rowKey, formatCell(prev)}

			for _, curr := range values[1:] {
				diff := apiversions.SliceDiff(prev, curr)
				diff.Sort()

				var add, remove string
				if len(diff.Added) > 0 {
					add = fmt.Sprintf("add:\n  %s\n\n", strings.Join(slice.Sort(diff.Added), "\n  "))
				}
				if len(diff.Removed) > 0 {
					add = fmt.Sprintf("remove:\n  %s\n\n", strings.Join(slice.Sort(diff.Removed), "\n  "))
				}
				row = append(row, fmt.Sprintf("%s%s", add, remove))

				prev = curr
			}
			return row
		})
	} else {
		return table.FormattedTable(func(rowKey string, values [][]string) []string {
			return slice.Cons(rowKey, slice.Map(formatCell, values))
		})
	}
}

func formatCell(items []string) string {
	return strings.Join(slice.Sort(items), "\n")
}
