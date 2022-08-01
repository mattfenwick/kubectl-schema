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

func RunCompareResource(args *CompareResourceArgs) {
	if len(args.Versions) != 2 {
		panic(errors.Errorf("expected 2 kube versions, found %+v", args.Versions))
	}

	allowedResources := set.NewSet(args.TypeNames)
	allowResource := func(name string) bool {
		return len(args.TypeNames) == 0 || allowedResources.Contains(name)
	}

	spec1 := MustReadSwaggerSpecFromGithub(MustVersion(args.Versions[0]))
	kinds1 := spec1.ResolveStructure()
	spec2 := MustReadSwaggerSpecFromGithub(MustVersion(args.Versions[1]))
	kinds2 := spec2.ResolveStructure()

	typeNames := set.NewSet(maps.Keys(kinds1))
	typeNames.Union(set.NewSet(maps.Keys(kinds2)))

	for _, typeName := range slice.Sort(typeNames.ToSlice()) {
		if allowResource(typeName) {
			logrus.Debugf("inspecting type %s", typeName)
		} else {
			logrus.Debugf("skipping type %s", typeName)
			continue
		}
		resolved1 := kinds1[typeName]
		resolved2 := kinds2[typeName]
		logrus.Debugf("group/versions for kube %s: %+v", args.Versions[0], maps.Keys(resolved1))
		logrus.Debugf("group/versions for kube %s: %+v", args.Versions[1], maps.Keys(resolved2))

		CompareSingleResource(typeName, resolved1, resolved2)
	}
}

func CompareSingleResource(typeName string, resolved1, resolved2 map[string]*ResolvedType) {
	for _, groupName1 := range maps.Keys(resolved1) {
		type1 := resolved1[groupName1]
		for _, groupName2 := range maps.Keys(resolved2) {
			type2 := resolved2[groupName2]
			//fmt.Printf("comparing %s: %s@%s vs. %s@%s\n", typeName, args.Versions[0], groupName1, args.Versions[1], groupName2)
			for _, e := range CompareResolvedResources(type1, type2).Elements {
				fmt.Printf("  %-20s    %+v\n", e.Type.Short(), strings.Join(e.Path, "."))
			}
			fmt.Println()
		}
	}
}
