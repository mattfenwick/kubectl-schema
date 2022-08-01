package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/set"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/kubectl-schema/pkg/swagger/apiversions"
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

type ShowKindsArgs struct {
	GroupBy            string
	KubeVersions       []string
	IncludeApiVersions []string
	ExcludeApiVersions []string
	IncludeResources   []string
	ExcludeResources   []string
	IncludeAll         bool
	Diff               bool
	// TODO add flag to verify parsing?  by serializing/deserializing to check if it matches input?
}

func (s *ShowKindsArgs) GetGroupBy() ShowKindsGroupBy {
	switch s.GroupBy {
	case "resource":
		return ShowKindsGroupByResource
	case "apiversion", "api-version":
		return ShowKindsGroupByApiVersion
	default:
		panic(errors.Errorf("invalid group by value: %s", s.GroupBy))
	}
}

func SetupShowKindsCommand() *cobra.Command {
	args := &ShowKindsArgs{}

	command := &cobra.Command{
		Use:   "kinds",
		Short: "show available kinds, by api-version and kubernetes version",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunShowKinds(args)
		},
	}

	command.Flags().BoolVar(&args.Diff, "diff", false, "if true, calculate a diff from kube version to kube version.  if true, simply print resources")

	command.Flags().BoolVar(&args.IncludeAll, "include-all", false, "if true, includes all apiversions and resources regardless of includes/excludes.  This is useful for debugging")

	command.Flags().StringVar(&args.GroupBy, "group-by", "resource", "what to group by: valid values are 'resource' and 'api-version'")
	command.Flags().StringSliceVar(&args.KubeVersions, "kube-version", defaultKubeVersions, "kube versions to explain")

	command.Flags().StringSliceVar(&args.ExcludeResources, "resource-exclude", []string{}, "resources to exclude")
	command.Flags().StringSliceVar(&args.IncludeResources, "resource", []string{}, "resources to include")

	command.Flags().StringSliceVar(&args.ExcludeApiVersions, "apiversion-exclude", []string{}, "api versions to exclude")
	command.Flags().StringSliceVar(&args.IncludeApiVersions, "apiversion", []string{}, "api versions to include")

	return command
}

func shouldAllow(s string, allows *set.Set[string], forbids *set.Set[string]) bool {
	return (allows.Len() == 0 || allows.Contains(s)) && !forbids.Contains(s)
}

func RunShowKinds(args *ShowKindsArgs) {
	var include func(apiVersion string, resource string) bool
	if args.IncludeAll {
		include = func(apiVersion string, resource string) bool {
			return true
		}
	} else {
		includeResources := set.NewSet(args.IncludeResources)
		excludeResources := set.NewSet(args.ExcludeResources)
		includeApiVersions := set.NewSet(args.IncludeApiVersions)
		excludeApiVersions := set.NewSet(args.ExcludeApiVersions)

		include = func(apiVersion string, resource string) bool {
			includeApiVersion := shouldAllow(apiVersion, includeApiVersions, excludeApiVersions)
			includeResource := shouldAllow(resource, includeResources, excludeResources)
			return includeApiVersion && includeResource
		}
	}

	fmt.Printf("\n%s\n\n", ShowKinds(args.GetGroupBy(), args.KubeVersions, include, args.Diff))
}

type ShowKindsGroupBy string

const (
	ShowKindsGroupByResource   = "ShowKindsGroupByResource"
	ShowKindsGroupByApiVersion = "ShowKindsGroupByApiVersion"
)

func (s ShowKindsGroupBy) Header() string {
	switch s {
	case ShowKindsGroupByResource:
		return "Resource"
	case ShowKindsGroupByApiVersion:
		return "API version"
	default:
		panic(errors.Errorf("invalid groupBy: %s", s))
	}
}

func ShowKinds(groupBy ShowKindsGroupBy, versions []string, include func(string, string) bool, calculateDiff bool) string {
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
					case ShowKindsGroupByResource:
						table.Add(gvk.Kind, kubeVersion.ToString(), apiVersion)
					case ShowKindsGroupByApiVersion:
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
