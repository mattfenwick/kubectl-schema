package swagger

import (
	"github.com/mattfenwick/collections/pkg/base"
	"github.com/mattfenwick/collections/pkg/function"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"strings"
)

type SpecType struct {
	AdditionalProperties        *SpecType                `json:"additionalProperties,omitempty"`
	Description                 string                   `json:"description,omitempty"`
	Format                      string                   `json:"format,omitempty"`
	Items                       *SpecType                `json:"items,omitempty"`
	Properties                  map[string]*SpecType     `json:"properties,omitempty"`
	Ref                         string                   `json:"$ref,omitempty"`
	Required                    []string                 `json:"required,omitempty"`
	Type                        string                   `json:"type,omitempty"`
	XKubernetesListMapKeys      []string                 `json:"x-kubernetes-list-map-keys,omitempty"`
	XKubernetesListType         string                   `json:"x-kubernetes-list-type,omitempty"`
	XKubernetesPatchMergeKey    string                   `json:"x-kubernetes-patch-merge-key,omitempty"`
	XKubernetesPatchStrategy    string                   `json:"x-kubernetes-patch-strategy,omitempty"`
	XKubernetesGroupVersionKind []*GVK                   `json:"x-kubernetes-group-version-kind,omitempty"`
	XKubernetesUnions           []map[string]interface{} `json:"x-kubernetes-unions,omitempty"`
}

type KubeSpec struct {
	Definitions map[string]*SpecType `json:"definitions"`
	Info        struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	//Paths map[string]interface{}
	//Security int
	//SecurityDefinitions int
}

func enforceInvariant(specType *SpecType) {
	counts := slice.Filter(function.Id[bool], []bool{specType.Ref != "", specType.Type != ""})
	if len(counts) != 1 && specType.Description == "" {
		logrus.Errorf("INVARIANT violated: %d; %+v", len(counts), specType)
	}
}

func (s *KubeSpec) MustGetDefinition(name string) *SpecType {
	val, ok := s.Definitions[name]
	if !ok {
		panic(errors.Errorf("unable to find definition for %s", name))
	}
	return val
}

func (s *KubeSpec) VisitSpecType(resolvedTypes map[string]*ResolvedType, path Path, specType *SpecType, visit func(path Path, resolved *ResolvedType, circular string)) *ResolvedType {
	enforceInvariant(specType)

	// visit AFTER processing of the type is done
	var resolved *ResolvedType
	defer visit(path, resolved, "")

	if specType.Ref != "" {
		refName := ParseRef(specType.Ref)
		newPath := path.Append(SpecPath{Ref: true})
		if resolvedTypes[refName] != nil {
			// done: NOT circular
			visit(newPath, resolvedTypes[refName], "")
			resolved = resolvedTypes[refName]
		} else if _, ok := resolvedTypes[refName]; ok {
			// in progress: circular
			visit(newPath, nil, refName)
			resolved = &ResolvedType{Circular: refName}
		} else {
			// hasn't been seen yet
			resolvedTypes[refName] = nil
			resolved = s.VisitSpecType(resolvedTypes, newPath, s.MustGetDefinition(refName), visit)
			resolvedTypes[refName] = resolved
		}
	} else {
		switch specType.Type {
		case "":
			logrus.Debugf("skipping empty type: %+v", strings.Join(path.ToStringPieces(), "."))
			resolved = &ResolvedType{Empty: true}
		case "array":
			resolved = &ResolvedType{Array: s.VisitSpecType(resolvedTypes, path.Append(SpecPath{Array: true}), specType.Items, visit)}
		case "object":
			obj := &ResolvedObject{Properties: map[string]*ResolvedType{}}
			for propName, prop := range specType.Properties {
				obj.Properties[propName] = s.VisitSpecType(resolvedTypes, path.Append(SpecPath{ObjectProperty: true}).Append(SpecPath{FieldAccess: propName}), prop, visit)
			}
			if specType.AdditionalProperties != nil {
				obj.AdditionalProperties = s.VisitSpecType(resolvedTypes, path.Append(SpecPath{FieldAccess: "additionalProperties"}), specType.AdditionalProperties, visit)
			}
			resolved = &ResolvedType{Object: obj}
		case "boolean", "string", "integer", "number":
			resolved = &ResolvedType{Primitive: specType.Type}
			logrus.Debugf("found primitive: %s", specType.Type)
		default:
			panic(errors.Errorf("TODO unsupported type %s: %+v, %+v", specType.Type, path, specType))
		}
	}
	return resolved
}

func (s *KubeSpec) Visit(visit func(path Path, resolved *ResolvedType, circular string)) (map[string]*ResolvedType, map[string]map[string]*ResolvedType) {
	resolvedTypes := map[string]*ResolvedType{}
	for defName, def := range s.Definitions {
		resolvedTypes[defName] = nil
		resolved := s.VisitSpecType(resolvedTypes, []SpecPath{{FieldAccess: defName}}, def, visit)
		resolvedTypes[defName] = resolved
	}
	byKindByAPIVersion := map[string]map[string]*ResolvedType{}
	for gvkString, resolved := range resolvedTypes {
		gvk := ParseGVK(gvkString)
		if _, ok := byKindByAPIVersion[gvk.Kind]; !ok {
			byKindByAPIVersion[gvk.Kind] = map[string]*ResolvedType{}
		}
		byKindByAPIVersion[gvk.Kind][gvk.GroupVersion()] = resolved
	}
	return resolvedTypes, byKindByAPIVersion
}

type SpecPath struct {
	FieldAccess    string
	Ref            bool
	ObjectProperty bool
	Array          bool
}

type Path []SpecPath

func (p Path) Append(piece SpecPath) Path {
	return slice.Append(p, []SpecPath{piece})
}

func (p Path) ToStringPieces() []string {
	var elems []string
	for _, piece := range p {
		if piece.FieldAccess != "" {
			elems = append(elems, piece.FieldAccess)
		} else if piece.Ref {
			// nothing to do
		} else if piece.Array {
			elems = append(elems, "[]")
		} else if piece.ObjectProperty {
			// nothing to do ??  TODO decide
		} else {
			panic(errors.Errorf("invalid SpecPath value: %+v", piece))
		}
	}
	return elems
}

func (s *KubeSpec) ResolveStructure() map[string]map[string]*ResolvedType {
	_, byKindByAPIVersion := s.Visit(func(path Path, resolved *ResolvedType, circular string) {
		if circular == "" {
			logrus.Debugf("%+v -- %+v\n", path.ToStringPieces(), resolved)
		} else {
			logrus.Debugf("%+v\n  CIRCULAR %s\n", path.ToStringPieces(), circular)
		}
	})
	return byKindByAPIVersion
}

//func (s *KubeSpec) ResolveGVKs() {
//	gvksByResource := map[string]map[string]*SpecType{}
//	s.Visit(func(path Path, resolved *ResolvedType, circular string) {
//		if circular == "" {
//			fmt.Printf("%+v -- %s\n", path.ToStringPieces(), resolved)
//		} else {
//			fmt.Printf("%+v\n  CIRCULAR %s\n", path.ToStringPieces(), circular)
//		}
//	})
//}

// TODO distinguish between gvk and parsed name
//type ResolvedGVK struct {
//	GVK *GVK
//	Name string
//	Type *ResolvedType
//}

type ResolvedObject struct {
	Properties           map[string]*ResolvedType
	AdditionalProperties *ResolvedType
}

type ResolvedType struct {
	Empty     bool
	Primitive string
	Array     *ResolvedType
	Object    *ResolvedObject
	Circular  string
}

func (r *ResolvedType) Paths(pathContext []string) []*base.Pair[[]string, string] {
	logrus.Debugf("path: %+v", pathContext)

	path := slice.Map(function.Id[string], pathContext)

	var out []*base.Pair[[]string, string]
	if r.Circular != "" {
		out = append(out, base.NewPair(path, r.Circular))
	} else if r.Primitive != "" {
		out = append(out, base.NewPair(path, r.Primitive))
	} else if r.Array != nil {
		out = append(out, base.NewPair(path, "array"))
		out = append(out, r.Array.Paths(slice.Append(path, []string{"[]"}))...)
	} else if r.Object != nil {
		out = append(out, base.NewPair(path, "object"))
		for _, fieldName := range slice.Sort(maps.Keys(r.Object.Properties)) {
			out = append(out, r.Object.Properties[fieldName].Paths(slice.Append(path, []string{fieldName}))...)
		}
		if r.Object.AdditionalProperties != nil {
			out = append(out, r.Object.AdditionalProperties.Paths(slice.Append(path, []string{"additionalProperties"}))...)
		}
	} else if r.Empty {
		out = append(out, base.NewPair(path, "?"))
	} else {
		panic(errors.Errorf("invalid ResolvedType: %+v", r))
	}
	return out
}
