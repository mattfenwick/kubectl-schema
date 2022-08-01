package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/kubectl-schema/pkg/diff"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var (
	defaultKubeVersions = []string{
		"1.22.12",
		"1.23.9",
		"1.24.3",
		"1.25.0-alpha.3",
	}
)

type ShowResourcesArgs struct {
	GroupBy      string
	KubeVersions []string
	ApiVersions  []string
	Resources    []string
	Diff         bool
	// TODO add flag to verify parsing?  by serializing/deserializing to check if it matches input?
}

func (s *ShowResourcesArgs) GetGroupBy() ShowResourcesGroupBy {
	switch s.GroupBy {
	case "resource":
		return ShowResourcesGroupByResource
	case "apiversion", "api-version":
		return ShowResourcesGroupByApiVersion
	default:
		panic(errors.Errorf("invalid group by value: %s", s.GroupBy))
	}
}

func SetupShowResourcesCommand() *cobra.Command {
	args := &ShowResourcesArgs{}

	command := &cobra.Command{
		Use:   "resources",
		Short: "show available resources, by api-version and kubernetes version",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunShowResources(args)
		},
	}

	command.Flags().BoolVar(&args.Diff, "diff", false, "if true, calculate a diff from kube version to kube version.  if true, simply print resources")

	command.Flags().StringVar(&args.GroupBy, "group-by", "resource", "what to group by: valid values are 'resource' and 'api-version'")
	command.Flags().StringSliceVar(&args.KubeVersions, "kube-version", defaultKubeVersions, "kube versions to explain")

	command.Flags().StringSliceVar(&args.Resources, "resource", []string{}, "resources to include; if empty, include all")
	command.Flags().StringSliceVar(&args.ApiVersions, "api-version", []string{}, "api versions to include; if empty, include all")

	return command
}

func RunShowResources(args *ShowResourcesArgs) {
	fmt.Printf("\n%s\n\n", ShowResources(args.GetGroupBy(), args.KubeVersions, apiVersionAndResourceAllower(args.ApiVersions, args.Resources), args.Diff))
}

type ShowResourcesGroupBy string

const (
	ShowResourcesGroupByResource   ShowResourcesGroupBy = "ShowResourcesGroupByResource"
	ShowResourcesGroupByApiVersion ShowResourcesGroupBy = "ShowResourcesGroupByApiVersion"
)

func (s ShowResourcesGroupBy) Header() string {
	switch s {
	case ShowResourcesGroupByResource:
		return "Resource"
	case ShowResourcesGroupByApiVersion:
		return "API version"
	default:
		panic(errors.Errorf("invalid groupBy: %s", s))
	}
}

func ShowResources(groupBy ShowResourcesGroupBy, versions []string, include func(string, string) bool, calculateDiff bool) string {
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
					case ShowResourcesGroupByResource:
						table.Add(gvk.Kind, kubeVersion.ToString(), apiVersion)
					case ShowResourcesGroupByApiVersion:
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
				cellDiff := diff.SliceDiff(prev, curr)
				cellDiff.Sort()

				var add, remove string
				if len(cellDiff.Added) > 0 {
					add = fmt.Sprintf("add:\n  %s\n\n", strings.Join(slice.Sort(cellDiff.Added), "\n  "))
				}
				if len(cellDiff.Removed) > 0 {
					add = fmt.Sprintf("remove:\n  %s\n\n", strings.Join(slice.Sort(cellDiff.Removed), "\n  "))
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
