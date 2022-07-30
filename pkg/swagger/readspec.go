package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/file"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	SpecsRootDirectory = "./swagger-data"
)

func ReadSwaggerSpecFromGithub(version KubeVersion) (*KubeSpec, error) {
	path := MakePathFromKubeVersion(version)

	if !file.Exists(path) {
		logrus.Infof("file for version %s not found (path %s); downloading instead", version, path)

		err := os.MkdirAll(SpecsRootDirectory, 0777)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to mkdir %s", SpecsRootDirectory)
		}

		err = utils.GetFileFromURL(version.SwaggerSpecURL(), path)
		if err != nil {
			return nil, err
		}

		// get the keys sorted
		err = json.SortFileOptions(path, false, true)
		if err != nil {
			return nil, err
		}
	}

	spec, err := json.ParseFile[KubeSpec](path)
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
	return fmt.Sprintf("%s/%s-swagger-spec.json", SpecsRootDirectory, version.ToString())
}
