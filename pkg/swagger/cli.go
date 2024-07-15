package swagger

import (
	"fmt"

	"github.com/mattfenwick/collections/pkg/json"
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
		Long: fmt.Sprintf(`This plugin provides schema inspection utilities, based on the OpenAPI kubernetes swagger specs.
It should not be used with kubernetes spec versions prior to 1.14.0.

This downloads specs to %s.

The data directory can be changed using the %s environment variable; if this variable
is not set, a directory underneath the home directory is created and used.
`, GetSpecsRootDirectory(), DataDirEnvVar),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return utils.SetUpLogger(flags.Verbosity)
		},
	}

	command.PersistentFlags().StringVarP(&flags.Verbosity, "verbosity", "v", "info", "log level; one of [info, debug, trace, warn, error, fatal, panic]")

	command.AddCommand(SetupVersionCommand())
	command.AddCommand(SetupExplainCommand())
	command.AddCommand(SetupCompareResourceCommand())
	command.AddCommand(SetupShowResourcesCommand())
	command.AddCommand(SetupConfigCommand())

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
	fmt.Printf("kubectl-schema version: \n%s\n", jsonString)
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
	command.Flags().StringSliceVar(&args.KubeVersions, "kube-version", []string{defaultKubeVersions[len(defaultKubeVersions)-1]}, "kubernetes spec versions")
	command.Flags().IntVar(&args.Depth, "depth", 0, "number of layers to print; 0 is treated as unlimited")
	command.Flags().StringSliceVar(&args.Paths, "path", []string{}, "paths to search under, components separated by '.'; if empty, all paths are searched")

	return command
}

func SetupCompareResourceCommand() *cobra.Command {
	args := &CompareResourceArgs{}

	command := &cobra.Command{
		Use:   "compare",
		Short: "compare types across kube versions",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunCompareResource(args)
		},
	}

	command.Flags().StringSliceVar(&args.ApiVersions, "api-version", []string{}, "api versions to use; if empty, uses all")

	command.Flags().StringSliceVar(&args.KubeVersions, "kube-version", []string{defaultKubeVersions[0], defaultKubeVersions[len(defaultKubeVersions)-1]}, "two kubernetes versions to compare (must be exactly 2)")
	command.Flags().StringSliceVar(&args.Resources, "resource", []string{"Pod"}, "resources to include; if empty, includes all")

	return command
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

	command.Flags().BoolVar(&args.Diff, "diff", false, "if true, calculate a diff from kube version to kube version.  if false, simply print resources")

	command.Flags().StringVar(&args.GroupBy, "group-by", "resource", "what to group by: valid values are 'resource' and 'api-version'")
	command.Flags().StringSliceVar(&args.KubeVersions, "kube-version", defaultKubeVersions, "kube versions to explain")

	command.Flags().StringSliceVar(&args.Resources, "resource", []string{}, "resources to include; if empty, include all")
	command.Flags().StringSliceVar(&args.ApiVersions, "api-version", []string{}, "api versions to include; if empty, include all")

	return command
}

func SetupConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "show compiled-in configuration",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			fmt.Printf("default kube versions:\n%s\n", json.MustMarshalToString(defaultKubeVersions))
			fmt.Printf("kube patch versions:\n%s\n", json.MustMarshalToString(LatestKubePatchVersionStrings))
		},
	}
	return command
}
