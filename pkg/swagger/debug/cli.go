package debug

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/file"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/kubectl-schema/pkg/swagger"
	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

func SetupSwaggerDebugCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "swagger-debug",
		Short: "debug kube swagger spec",
		Args:  cobra.ExactArgs(0),
	}

	command.AddCommand(setupParseCommand())
	command.AddCommand(setupAnalyzeSchemaCommand())
	command.AddCommand(setupTestSchemaParserCommand())

	return command
}

type ParseArgs struct {
	Version string
}

func setupParseCommand() *cobra.Command {
	args := &ParseArgs{}

	command := &cobra.Command{
		Use:   "parse",
		Short: "parse and serialize openapi spec for comparison (test command)",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunParse(args)
		},
	}

	command.Flags().StringVar(&args.Version, "version", "1.23.0", "kubernetes version")

	return command
}

func RunParse(args *ParseArgs) {
	spec := swagger.MustReadSwaggerSpecFromGithub(swagger.MustVersion(args.Version))

	for name, t := range spec.Definitions {
		for propName, prop := range t.Properties {
			logrus.Debugf("%s, %s: %+v\n<<>>\n", name, propName, prop.Items)
		}
	}

	// must do weird marshal/unmarshal/marshal dance to get struct keys sorted
	bytes, err := json.MarshalWithOptions(spec, &json.MarshalOptions{EscapeHTML: true, Indent: true, Sort: true})
	utils.DoOrDie(err)

	fmt.Printf("%s\n", bytes)
}

type AnalyzeSchemaArgs struct {
	Version string
	All     bool
}

func setupAnalyzeSchemaCommand() *cobra.Command {
	args := &AnalyzeSchemaArgs{}

	command := &cobra.Command{
		Use:   "analyze-schema",
		Short: "analyze shape of openapi schema",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			RunAnalyzeSchema(args)
		},
	}

	command.Flags().StringVar(&args.Version, "version", "1.24.2", "kubernetes version")
	command.Flags().BoolVar(&args.All, "all", false, "if true, analyze all latest versions")

	return command
}

func RunAnalyzeSchema(args *AnalyzeSchemaArgs) {
	if args.All {
		RunAnalyzeSchemaLatest()
		return
	}

	path := swagger.MakePathFromKubeVersion(swagger.MustVersion(args.Version))
	specObj, err := json.ParseFile[map[string]interface{}](path)
	utils.DoOrDie(err)
	//spec := MustReadSwaggerSpecFromGithub(args.Version) // TODO

	//starterPaths := []string{"paths", "definitions"}
	starterPaths := []string{"definitions"}
	paths, schemaPaths := JsonFindPaths(specObj, starterPaths)
	for _, p := range paths {
		if false {
			fmt.Printf("%s\n", strings.Join(p, " "))
		}
	}
	for _, p := range schemaPaths {
		fmt.Printf("%s\n", strings.Join(p, " "))
	}
}

func RunAnalyzeSchemaLatest() {
	for _, version := range swagger.LatestKubePatchVersions {
		path := fmt.Sprintf("test-schema/%s.txt", version)
		specBytes := swagger.MustDownloadSwaggerSpec(version)
		specObj, err := json.Parse[map[string]interface{}](specBytes)
		utils.DoOrDie(err)
		//spec := MustReadSwaggerSpecFromGithub(args.Version) // TODO

		//starterPathsToInspect := []string{"paths", "definitions"}
		starterPathsToInspect := []string{"definitions"}
		schemaPaths, dedupedPaths := JsonFindPaths(specObj, starterPathsToInspect)
		for _, p := range schemaPaths {
			if false {
				fmt.Printf("%s\n", strings.Join(p, " "))
			}
		}
		lines := slice.Map(func(xs []string) string { return strings.Join(xs, " ") }, dedupedPaths)
		err = file.Write(path, []byte(strings.Join(lines, "\n")), 0644)
		utils.DoOrDie(err)
		//for _, p := range dedupedPaths {
		//	fmt.Printf("%s\n", strings.Join(p, " "))
		//}
	}
}
