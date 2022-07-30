package debug

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/file"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/kubectl-schema/pkg/swagger"
	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"path"
	"reflect"
)

func setupTestSchemaParserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "test-schema-parser",
		Short: "make sure schema parser handles everything",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, as []string) {
			TestSchemaParser()
		},
	}

	return command
}

func TestSchemaParser() {
	schemaDir := "test-schema"
	for _, version := range swagger.LatestKubePatchVersions {
		CheckSchema(path.Join(schemaDir, version.ToString()), version)
	}
}

func CheckSchema(dir string, version swagger.KubeVersion) {
	specBytes := swagger.MustDownloadSwaggerSpec(version)
	// remove paths
	specMap, err := json.Parse[map[string]interface{}](specBytes)
	utils.DoOrDie(err)
	delete(*specMap, "paths")
	specMapBytes, err := json.MarshalWithOptions(specMap, json.DefaultMarshalOptions)
	utils.DoOrDie(err)
	// carry on
	spec, err := json.Parse[swagger.KubeSpec](specMapBytes)
	utils.DoOrDie(err)

	specString, err := json.MarshalToString(spec)
	utils.DoOrDie(err)
	spec2, err := json.Parse[swagger.KubeSpec]([]byte(specString))
	utils.DoOrDie(err)

	utils.DoOrDie(os.MkdirAll(dir, 0777))
	path1, path2 := path.Join(dir, "spec1.txt"), path.Join(dir, "spec2.txt")

	sortedSpecBytes, err := json.SortOptions(specMapBytes, false, true)
	utils.DoOrDie(err)
	utils.DoOrDie(file.Write(path1, sortedSpecBytes, 0644))

	sortedSpecStringBytes, err := json.SortOptions([]byte(specString), false, true)
	utils.DoOrDie(err)
	utils.DoOrDie(file.Write(path2, sortedSpecStringBytes, 0644))

	//diff, err := utils.CommandRun(exec.Command("git", "diff", "--no-index", "my-spec-1.txt", "my-spec-2.txt"))
	//utils.DoOrDie(err)
	//fmt.Printf("%s\n", diff)
	//spec := MustReadSwaggerSpecFromGithub(version)
	//specString := utils.JsonString(spec)
	//
	//spec2, err := utils.ParseJson[Spec]([]byte(specString))
	//utils.DoOrDie(err)
	//specString2 := utils.JsonString(spec2)

	fmt.Printf("same? %t, %t\n", reflect.DeepEqual(spec, spec2), specString == string(specBytes))
}
