package utils

import (
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

func DoOrDie(err error) {
	if err != nil {
		logrus.Fatalf("%+v", err)
	}
}

func Pointer(s string) *string {
	return &s
}

func AddKey(dict map[string]bool, key string) map[string]bool {
	out := map[string]bool{}
	for k, v := range dict {
		out[k] = v
	}
	out[key] = true
	return out
}

func SortedKeys[K constraints.Ordered, V any](xs map[K]V) []K {
	return slice.Sort(maps.Keys(xs))
}

func StringPrefix(s string, chars int) string {
	if len(s) <= chars {
		return s
	}
	return s[:chars]
}

func Set(xs []string) map[string]bool {
	out := map[string]bool{}
	for _, x := range xs {
		out[x] = true
	}
	return out
}

func CopySlice[A any](s []A) []A {
	newCopy := make([]A, len(s))
	copy(newCopy, s)
	return newCopy
}
