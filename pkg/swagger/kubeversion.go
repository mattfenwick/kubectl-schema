package swagger

import (
	"fmt"
	"strings"

	"github.com/mattfenwick/collections/pkg/base"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/pkg/errors"
)

type KubeVersion []string

var (
	CompareKubeVersion = slice.CompareSlicePairwise[string]()
)

func NewVersion(v string) (KubeVersion, error) {
	pieces := strings.Split(v, ".")
	if len(pieces) < 3 {
		return nil, errors.Errorf("expected at least 3 pieces, found [%+v]", pieces)
	}
	return pieces, nil
}

func MustVersion(v string) KubeVersion {
	version, err := NewVersion(v)
	utils.DoOrDie(err)
	return version
}

func (v KubeVersion) Compare(b KubeVersion) base.Ordering {
	return CompareKubeVersion(v, b)
}

func (v KubeVersion) ToString() string {
	return strings.Join(v, ".")
}

func (v KubeVersion) SwaggerSpecURL() string {
	return fmt.Sprintf(GithubOpenapiURLTemplate, v.ToString())
}

var (
	GithubOpenapiURLTemplate = "https://raw.githubusercontent.com/kubernetes/kubernetes/v%s/api/openapi-spec/swagger.json"

	// LatestKubePatchVersionStrings records the latest known patch versions for each minor version
	//   these version numbers come from https://github.com/kubernetes/kubernetes/tree/master/CHANGELOG
	LatestKubePatchVersionStrings = []string{
		// there's nothing listed for 1.1
		//"1.2.7", // for some reason, these don't show up on the openapi github specs
		//"1.3.10",
		//"1.4.12",
		"1.5.8",
		"1.6.13",
		"1.7.16",
		"1.8.15",
		"1.9.11",
		"1.10.13",
		"1.11.10",
		"1.12.10",
		"1.13.12",
		"1.14.10",
		"1.15.12",
		"1.16.15",
		"1.17.17",
		"1.18.20",
		"1.19.16",
		"1.20.15",
		"1.21.14",
		"1.22.17",
		"1.23.17",
		"1.24.12",
		"1.25.8",
		"1.26.3",
		"1.27.0-rc.1",
	}

	LatestKubePatchVersions = slice.Map(MustVersion, LatestKubePatchVersionStrings)
)

var (
	defaultKubeVersions = LatestKubePatchVersionStrings[len(LatestKubePatchVersionStrings)-4:]
)
