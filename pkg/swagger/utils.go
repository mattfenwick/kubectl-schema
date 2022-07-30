package swagger

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type GVK struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

func (g *GVK) GroupVersion() string {
	if g.Group == "" {
		return g.Version
	}
	return fmt.Sprintf("%s.%s", g.Group, g.Version)
}

func (g *GVK) ToString() string {
	return fmt.Sprintf("%s.%s", g.GroupVersion(), g.Kind)
}

func ParseRef(ref string) string {
	pieces := strings.Split(ref, "/")
	if len(pieces) != 3 {
		panic(errors.Errorf("unable to parse ref: expected 3 pieces, found %d (%s)", len(pieces), ref))
	}
	return pieces[2]
}

func ParseGVK(gvk string) *GVK {
	split := strings.Split(gvk, ".")
	if len(split) < 3 {
		panic(errors.Errorf("invalid gvk string: %s", gvk))
	}
	return &GVK{
		Group:   strings.Join(split[:len(split)-2], "."),
		Version: split[len(split)-2],
		Kind:    split[len(split)-1],
	}
}
