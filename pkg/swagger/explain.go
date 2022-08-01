package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"strings"
)

type ExplainArgs struct {
	Format       string
	ApiVersions  []string
	Resources    []string
	KubeVersions []string
	Depth        int
	Paths        []string
}

func SetupExplainCommand() *cobra.Command {
	args := &ExplainArgs{}

	command := &cobra.Command{
		Use:   "explain",
		Short: "explain resources from a swagger spec",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunExplain(args)
		},
	}

	command.Flags().StringVar(&args.Format, "format", "condensed", "output format; possible values: table, condensed")
	command.Flags().StringSliceVar(&args.ApiVersions, "api-version", []string{}, "api versions to look for resource under; looks under all if not specified")
	command.Flags().StringSliceVar(&args.Resources, "resource", []string{}, "kubernetes resources to explain")
	command.Flags().StringSliceVar(&args.KubeVersions, "kube-version", []string{"1.23.0"}, "kubernetes spec versions")
	command.Flags().IntVar(&args.Depth, "depth", 0, "number of layers to print; 0 is treated as unlimited")
	command.Flags().StringSliceVar(&args.Paths, "path", []string{}, "paths to search under, components separated by '.'; if empty, all paths are searched")

	return command
}

func RunExplain(args *ExplainArgs) {
	allowApiVersion := allower(args.ApiVersions)
	allowResource := allower(args.Resources)

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
		typesByKindByApiVersion := spec.ResolveStructure()

		for _, resourceName := range slice.Sort(maps.Keys(typesByKindByApiVersion)) {
			if !allowResource(resourceName) {
				continue
			}
			typesByApiVersion := map[string]*ResolvedType{}
			for apiVersion, resolvedType := range typesByKindByApiVersion[resourceName] {
				if allowApiVersion(apiVersion) {
					typesByApiVersion[apiVersion] = resolvedType
				}
			}

			switch args.Format {
			case "table":
				for apiVersion, resolvedType := range typesByApiVersion {
					fmt.Printf("%s %s:\n", apiVersion, resourceName)
					fmt.Printf("%s\n\n", TableResource(resolvedType, allowPath))
				}
			case "condensed":
				fmt.Printf("%s:\n", resourceName)
				for apiVersion, resolvedType := range typesByApiVersion {
					fmt.Printf("%s\n\n", CondensedResource(apiVersion, resolvedType, allowPath))
				}
			default:
				panic(errors.Errorf("invalid output format: %s", args.Format))
			}
		}
	}
}

func TableResource(resolvedType *ResolvedType, allowPath func([]string) bool) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetColMinWidth(1, 100)
	table.SetHeader([]string{"Kind", "Field"})
	for _, pair := range resolvedType.Paths([]string{}) {
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

func CondensedResource(apiVersion string, resolvedType *ResolvedType, allowPath func([]string) bool) string {
	lines := []string{apiVersion + ":"}
	for _, pair := range resolvedType.Paths([]string{}) {
		path, vType := pair.Fst, pair.Snd
		if len(path) > 0 && allowPath(path) {
			prefix := strings.Repeat("  ", len(path)-1)
			typeString := fmt.Sprintf("%s%s", prefix, path[len(path)-1])
			lines = append(lines, fmt.Sprintf("%-60s    %s", typeString, vType))
		}
	}
	return strings.Join(lines, "\n")
}
