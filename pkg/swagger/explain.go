package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/set"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"strings"
)

type ExplainResourceArgs struct {
	Format        string
	GroupVersions []string
	TypeNames     []string
	KubeVersions  []string
	Depth         int
	Paths         []string
}

func setupExplainResourceCommand() *cobra.Command {
	args := &ExplainResourceArgs{}

	command := &cobra.Command{
		Use:   "explain",
		Short: "explain types from a swagger spec",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunExplainResource(args)
		},
	}

	command.Flags().StringVar(&args.Format, "format", "condensed", "output format")
	command.Flags().StringSliceVar(&args.GroupVersions, "group-version", []string{}, "group/versions to look for type under; looks under all if not specified")
	command.Flags().StringSliceVar(&args.TypeNames, "type", []string{}, "kubernetes types to explain")
	command.Flags().StringSliceVar(&args.KubeVersions, "version", []string{"1.23.0"}, "kubernetes spec versions")
	command.Flags().IntVar(&args.Depth, "depth", 0, "number of layers to print; 0 is treated as unlimited")
	command.Flags().StringSliceVar(&args.Paths, "path", []string{}, "paths to search under, components separated by '.'; if empty, all paths are searched")

	return command
}

func RunExplainResource(args *ExplainResourceArgs) {
	allowedGVs := set.NewSet(args.GroupVersions)
	allowGV := func(gv string) bool {
		return len(args.GroupVersions) == 0 || allowedGVs.Contains(gv)
	}
	allowedResources := set.NewSet(args.TypeNames)
	allowResource := func(name string) bool {
		return len(args.TypeNames) == 0 || allowedResources.Contains(name)
	}
	allowDepth := func(prefix int, depth int) bool {
		if args.Depth == 0 {
			// always allow if maxDepth is unset
			return true
		}
		return (depth - prefix) <= args.Depth
	}
	allowedPaths := slice.Map(func(p string) []string { return strings.Split(p, ".") }, args.Paths)
	allowPath := func(path []string) bool {
		if len(allowedPaths) == 0 {
			return true
		}
		for _, prefix := range allowedPaths {
			if IsPrefixOf(prefix, path) && allowDepth(len(prefix), len(path)) {
				return true
			}
		}
		return false
	}

	//table := NewPivotTable("?", args.KubeVersions)

	for _, kubeVersion := range args.KubeVersions {
		fmt.Printf("%s\n", kubeVersion)
		spec := MustReadSwaggerSpecFromGithub(MustVersion(kubeVersion))
		resolvedGVKs := spec.ResolveStructure()

		for _, name := range slice.Sort(maps.Keys(resolvedGVKs)) {
			if !allowResource(name) {
				continue
			}
			gvks := map[string]*ResolvedType{}
			for gv, kind := range resolvedGVKs[name] {
				if allowGV(gv) {
					gvks[gv] = kind
				}
			}

			switch args.Format {
			case "table":
				for gv, kind := range gvks {
					fmt.Printf("%s %s:\n", gv, name)
					fmt.Printf("%s\n\n", TableResource(kind, allowPath))
				}
			case "condensed":
				fmt.Printf("%s:\n", name)
				for gv, kind := range gvks {
					fmt.Printf("%s\n\n", CondensedResource(gv, kind, allowPath))
				}
			default:
				panic(errors.Errorf("invalid output format: %s", args.Format))
			}
		}
	}
}

func TableResource(kind *ResolvedType, allowPath func([]string) bool) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetColMinWidth(1, 100)
	table.SetHeader([]string{"Type", "Field"})
	for _, pair := range kind.Paths([]string{}) {
		path, vType := pair.Fst, pair.Snd
		if allowPath(path) {
			table.Append([]string{strings.Join(path, "."), vType})
		}
	}
	table.Render()
	return tableString.String()
}

func IsPrefixOf[A comparable](xs []A, ys []A) bool {
	for i := 0; i < len(xs); i++ {
		if i >= len(ys) || xs[i] != ys[i] {
			return false
		}
	}
	return true
}

func CondensedResource(gv string, kind *ResolvedType, allowPath func([]string) bool) string {
	lines := []string{gv + ":"}
	for _, pair := range kind.Paths([]string{}) {
		path, vType := pair.Fst, pair.Snd
		if len(path) > 0 && allowPath(path) {
			prefix := strings.Repeat("  ", len(path)-1)
			typeString := fmt.Sprintf("%s%s", prefix, path[len(path)-1])
			lines = append(lines, fmt.Sprintf("%-60s    %s", typeString, vType))
		}
	}
	return strings.Join(lines, "\n")
}
