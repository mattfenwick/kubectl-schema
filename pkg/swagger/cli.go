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
	command.AddCommand(setupCompareResourceCommand())
	command.AddCommand(SetupShowResourcesCommand())

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
