package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/file"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

const (
	DataDirEnvVar = "KUBECTL_SCHEMA_DATA_DIRECTORY"
)

func getHomeDir() string {
	home, err := os.UserHomeDir()
	utils.DoOrDie(errors.Wrapf(err, "unable to get home dir"))
	return home
}

func GetSpecsRootDirectory() string {
	if dataDir, ok := os.LookupEnv(DataDirEnvVar); ok {
		return dataDir
	}
	return path.Join(getHomeDir(), ".kubectl-schema")
}

func ReadSwaggerSpecFromGithub(version KubeVersion) (*KubeSpec, error) {
	specPath := MakePathFromKubeVersion(version)

	if !file.Exists(specPath) {
		logrus.Infof("file for version %s not found (path %s); downloading instead", version, specPath)

		dataDir := GetSpecsRootDirectory()
		err := os.MkdirAll(dataDir, 0777)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to mkdir %s", dataDir)
		}

		err = utils.GetFileFromURL(version.SwaggerSpecURL(), specPath)
		if err != nil {
			return nil, err
		}

		// get the keys sorted
		err = json.SortFileOptions(specPath, false, true)
		if err != nil {
			return nil, err
		}
	}

	spec, err := json.ParseFile[KubeSpec](specPath)
	utils.DoOrDie(err)

	return spec, nil
}

func MustReadSwaggerSpecFromGithub(version KubeVersion) *KubeSpec {
	spec, err := ReadSwaggerSpecFromGithub(version)
	utils.DoOrDie(err)
	return spec
}

func MustDownloadSwaggerSpec(version KubeVersion) []byte {
	bytes, err := utils.GetURL(version.SwaggerSpecURL())
	utils.DoOrDie(err)
	return bytes
}

func MakePathFromKubeVersion(version KubeVersion) string {
	return fmt.Sprintf("%s/%s-swagger-spec.json", GetSpecsRootDirectory(), version.ToString())
}
