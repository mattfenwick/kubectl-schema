package swagger

import (
	"github.com/pkg/errors"
)

func RunCompareResource(args *CompareResourceArgs) {
	if len(args.Versions) != 2 {
		panic(errors.Errorf("expected 2 kube versions, found %+v", args.Versions))
	}

	// TODO
	//swaggerSpec1 := MustReadSwaggerSpecFromGithub(MustVersion(args.Versions[0]))
	//swaggerSpec2 := MustReadSwaggerSpecFromGithub(MustVersion(args.Versions[1]))
	//
	//typeNames := map[string]interface{}{}
	//if len(args.TypeNames) > 0 {
	//	for _, name := range args.TypeNames {
	//		typeNames[name] = true
	//	}
	//} else {
	//	for name := range swaggerSpec1.DefinitionsByNameByGroup() {
	//		typeNames[name] = true
	//	}
	//	for name := range swaggerSpec2.DefinitionsByNameByGroup() {
	//		typeNames[name] = true
	//	}
	//}
	//
	//for _, typeName := range utils.SortedKeys(typeNames) {
	//	fmt.Printf("inspecting type %s\n", typeName)
	//
	//	resolved1 := ResolveResource(swaggerSpec1, typeName)
	//	resolved2 := ResolveResource(swaggerSpec2, typeName)
	//
	//	logrus.Infof("group/versions for kube %s: %+v", args.Versions[0], utils.SortedKeys(resolved1))
	//	logrus.Infof("group/versions for kube %s: %+v", args.Versions[1], utils.SortedKeys(resolved2))
	//
	//	for _, groupName1 := range utils.SortedKeys(resolved1) {
	//		type1 := resolved1[groupName1]
	//		for _, groupName2 := range utils.SortedKeys(resolved2) {
	//			type2 := resolved2[groupName2]
	//			fmt.Printf("comparing %s: %s@%s vs. %s@%s\n", typeName, args.Versions[0], groupName1, args.Versions[1], groupName2)
	//			for _, e := range CompareResolvedResources(type1, type2).Elements {
	//				//for _, e := range utils.DiffJsonValues(utils.MustJsonRemarshal(type1), utils.MustJsonRemarshal(type2)).Elements {
	//				if len(e.Path) > 0 && e.Path[len(e.Path)-1] == "description" && args.SkipDescriptions {
	//					logrus.Debugf("skipping description at %+v", e.Path)
	//				} else {
	//					fmt.Printf("  %-20s    %+v\n", e.Type.Short(), strings.Join(e.Path, "."))
	//					if args.PrintValues {
	//						fmt.Printf("  - old: %+v\n  - new: %+v\n", e.Old, e.New)
	//					}
	//				}
	//			}
	//			fmt.Println()
	//		}
	//	}
	//}
}
