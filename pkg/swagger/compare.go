package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/set"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"strings"
)

type CompareResourceArgs struct {
	KubeVersions []string
	ApiVersions  []string
	Resources    []string
}

func RunCompareResource(args *CompareResourceArgs) {
	if len(args.KubeVersions) != 2 {
		panic(errors.Errorf("expected 2 kube versions, found %+v", args.KubeVersions))
	}

	allowResource := allower(args.Resources)
	allowApiVersion := allower(args.ApiVersions)

	spec1 := MustReadSwaggerSpecFromGithub(MustVersion(args.KubeVersions[0]))
	kinds1 := spec1.ResolveStructure()
	spec2 := MustReadSwaggerSpecFromGithub(MustVersion(args.KubeVersions[1]))
	kinds2 := spec2.ResolveStructure()

	typeNames := set.FromSlice(maps.Keys(kinds1)).Union(set.FromSlice(maps.Keys(kinds2)))

	for _, typeName := range slice.Sort(typeNames.ToSlice()) {
		if allowResource(typeName) {
			logrus.Debugf("inspecting type %s", typeName)
		} else {
			logrus.Debugf("skipping type %s", typeName)
			continue
		}
		resolved1 := kinds1[typeName]
		resolved2 := kinds2[typeName]
		logrus.Debugf("api versions for kube %s: %+v", args.KubeVersions[0], maps.Keys(resolved1))
		logrus.Debugf("api versions for kube %s: %+v", args.KubeVersions[1], maps.Keys(resolved2))

		for _, apiVersion1 := range maps.Keys(resolved1) {
			if !allowApiVersion(apiVersion1) {
				continue
			}
			type1 := resolved1[apiVersion1]
			for _, apiVersion2 := range maps.Keys(resolved2) {
				if !allowApiVersion(apiVersion2) {
					continue
				}
				type2 := resolved2[apiVersion2]
				fmt.Printf("comparing %s: %s@%s vs. %s@%s\n", typeName, args.KubeVersions[0], apiVersion1, args.KubeVersions[1], apiVersion2)
				for _, e := range CompareResolvedResources(type1, type2).Changes {
					fmt.Printf("  %-20s    %+v\n", e.Kind.Short(), strings.Join(e.Path, "."))
				}
				fmt.Println()
			}
		}
	}
}
