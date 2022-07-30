package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/collections/pkg/set"
	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func RunRootSchemaCommand() {
	command := SetupRootSchemaCommand()
	utils.DoOrDie(errors.Wrapf(command.Execute(), "run root schema command"))
}

type RootSchemaFlags struct {
	Verbosity string
}

func SetupRootSchemaCommand() *cobra.Command {
	flags := &RootSchemaFlags{}
	command := &cobra.Command{
		Use:   "schema",
		Short: "schema inspection utilities",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.SetUpLogger(flags.Verbosity)
		},
	}

	command.PersistentFlags().StringVarP(&flags.Verbosity, "verbosity", "v", "info", "log level; one of [info, debug, trace, warn, error, fatal, panic]")

	command.AddCommand(SetupVersionCommand())
	command.AddCommand(setupExplainResourceCommand())
	command.AddCommand(setupCompareResourceCommand())
	command.AddCommand(setupExplainGvkCommand())

	return command
}

var (
	version   = "development"
	gitSHA    = "development"
	buildTime = "development"
)

func SetupVersionCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "print out version information",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunVersionCommand()
		},
	}

	return command
}

func RunVersionCommand() {
	jsonString, err := json.MarshalToString(map[string]string{
		"Version":   version,
		"GitSHA":    gitSHA,
		"BuildTime": buildTime,
	})
	utils.DoOrDie(err)
	fmt.Printf("KubeUtils version: \n%s\n", jsonString)
}

var (
	//defaultExcludeResources = []string{"WatchEvent", "DeleteOptions"}
	//defaultIncludeResources = []string{
	//	"Service",
	//	"ClusterRole",
	//	"ClusterRoleBinding",
	//	"ConfigMap",
	//	"CronJob",
	//	"CustomResourceDefinition",
	//	"Deployment",
	//	"Ingress",
	//	"Job",
	//	"Role",
	//	"RoleBinding",
	//	"Secret",
	//	"ServiceAccount",
	//	"StatefulSet",
	//}
	//
	//defaultExcludeApiVersions = []string{}
	//defaultIncludeApiVersions = []string{
	//	"v1",
	//	"apps.v1",
	//	"batch.v1",
	//}

	defaultKubeVersions = []string{
		"1.18.20",
		"1.19.16",
		"1.20.15",
		"1.21.14",
		"1.22.12",
		"1.23.9",
		"1.24.3",
		"1.25.0-alpha.3",
	}
)

type ExplainGVKArgs struct {
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

func (e *ExplainGVKArgs) GetGroupBy() ExplainGVKGroupBy {
	switch e.GroupBy {
	case "resource":
		return ExplainGVKGroupByResource
	case "apiversion", "api-version":
		return ExplainGVKGroupByApiVersion
	default:
		panic(errors.Errorf("invalid group by value: %s", e.GroupBy))
	}
}

func setupExplainGvkCommand() *cobra.Command {
	args := &ExplainGVKArgs{}

	command := &cobra.Command{
		Use:   "gvk",
		Short: "explain gvks from a swagger spec",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunExplainGvks(args)
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

func RunExplainGvks(args *ExplainGVKArgs) {
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

	fmt.Printf("\n%s\n\n", ExplainGvks(args.GetGroupBy(), args.KubeVersions, include, args.Diff))
}

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

type CompareResourceArgs struct {
	Versions []string
	//GroupVersions []string // TODO ?
	TypeNames        []string
	SkipDescriptions bool
	PrintValues      bool
}

func setupCompareResourceCommand() *cobra.Command {
	args := &CompareResourceArgs{}

	command := &cobra.Command{
		Use:   "compare",
		Short: "compare types across kube versions",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunCompareResource(args)
		},
	}

	//command.Flags().StringSliceVar(&args.GroupVersions, "group-version", []string{}, "group/versions to look for type under; looks under all if not specified")
	//utils.DoOrDie(command.MarkFlagRequired("group-version"))

	command.Flags().StringSliceVar(&args.Versions, "version", []string{"1.18.19", "1.23.0"}, "kubernetes versions")
	command.Flags().StringSliceVar(&args.TypeNames, "type", []string{"Pod"}, "types to compare")
	command.Flags().BoolVar(&args.SkipDescriptions, "skip-descriptions", true, "if true, skip comparing descriptions (since these often change for non-functional reasons)")
	command.Flags().BoolVar(&args.PrintValues, "print-values", false, "if true, print values (in addition to just the path and change type)")

	return command
}
