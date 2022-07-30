package debug

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/builtin"
	"github.com/mattfenwick/collections/pkg/function"
	"github.com/mattfenwick/collections/pkg/set"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/kubectl-schema/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	goyaml "gopkg.in/yaml.v3"
	"sigs.k8s.io/yaml"
)

type JsonPaths struct {
	Paths [][]string
}

func (j *JsonPaths) Append(path []string) {
	j.Paths = append(j.Paths, slice.Map(function.Id[string], path))
}

func (j *JsonPaths) GetSortedPaths() [][]string {
	return slice.SortBy(slice.CompareSlicePairwiseBy(builtin.CompareOrdered[string]), j.Paths)
}

type Trie struct {
	Values   *set.Set[string]
	Children map[string]*Trie
}

func NewTrie() *Trie {
	return &Trie{
		Values:   set.NewSet[string](nil),
		Children: map[string]*Trie{},
	}
}

func (t *Trie) Add(path []string, value string) {
	if len(path) == 0 {
		t.Values.Add(value)
	} else {
		first, rest := path[0], path[1:]
		if _, ok := t.Children[first]; !ok {
			t.Children[first] = NewTrie()
		}
		t.Children[first].Add(rest, value)
	}
}

func (t *Trie) GetPaths(pathContext []string, paths *JsonPaths) {
	//if t.Values.Len() > 0 {
	//	paths.Append(pathContext)
	//}
	for _, val := range t.Values.ToSlice() {
		paths.Append(append(pathContext, val))
	}
	for key, child := range t.Children {
		child.GetPaths(append(pathContext, key), paths)
	}
}

func BounceMarshalGeneric[A any](in interface{}) (*A, error) {
	yamlBytes, err := goyaml.Marshal(in)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to marshal yaml")
	}
	var out A
	err = yaml.UnmarshalStrict(yamlBytes, &out)
	return &out, errors.Wrapf(err, "unable to unmarshal k8s yaml")
}

func JsonFindPaths(obj interface{}, starterPaths []string) ([][]string, [][]string) {
	paths := &JsonPaths{}
	bouncedObj, err := BounceMarshalGeneric[map[string]interface{}](obj)
	utils.DoOrDie(err)
	JsonFindPathsHelper(*bouncedObj, []string{}, paths)

	trie := NewTrie()
	schemaPaths := &JsonPaths{}
	for _, key := range starterPaths {
		nextLevel := (*bouncedObj)[key].(map[string]interface{})
		for _, val := range nextLevel {
			JsonFindPathsHelper(val, []string{key}, schemaPaths)
		}
	}
	for _, path := range schemaPaths.GetSortedPaths() {
		trie.Add(path[:len(path)-1], path[len(path)-1])
	}
	triePaths := &JsonPaths{}
	trie.GetPaths([]string{}, triePaths)

	return paths.GetSortedPaths(), triePaths.GetSortedPaths()
}

func JsonFindPathsHelper(obj interface{}, pathContext []string, paths *JsonPaths) {
	path := make([]string, len(pathContext))
	copy(path, pathContext)

	logrus.Debugf("path: %+v", path)

	if obj == nil {
		panic(errors.Errorf("unexpected nil at %+v", path))
	} else {
		switch val := obj.(type) {
		case map[string]interface{}:
			for _, k := range maps.Keys(val) {
				JsonFindPathsHelper(val[k], append(path, fmt.Sprintf(`["%s"]`, k)), paths)
			}
		case []interface{}:
			for i := range val {
				newPath := append(path, fmt.Sprintf("[%d]", i))
				JsonFindPathsHelper(val[i], newPath, paths)
			}
		case int:
			paths.Append(append(path, fmt.Sprintf("%T", val)))
		case string:
			paths.Append(append(path, fmt.Sprintf("%T", val)))
		case bool:
			paths.Append(append(path, fmt.Sprintf("%T", val)))
		//case types.Nil: // TODO is this necessary?
		default:
			panic(errors.Errorf("unrecognized type: %+v, %T", path, val))
		}
	}
}
