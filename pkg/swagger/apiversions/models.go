package apiversions

import (
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/mattfenwick/kubectl-schema/pkg/diff"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"reflect"
	"sort"
	"strings"
)

type ResourceDiff struct {
	Changed map[string]*diff.MapDiff
}

func (r *ResourceDiff) SortedChangedKeys() []string {
	return slice.Sort(maps.Keys(r.Changed))
}

func (r *ResourceDiff) Table(includes map[string]bool, skips map[string]bool) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetHeader([]string{"Kind", "Added", "Removed", "Same"})
	for _, kind := range r.SortedChangedKeys() {
		change := r.Changed[kind]
		change.Sort()
		if (len(includes) == 0 || includes[kind]) && !skips[kind] && (len(change.Added) != 0 || len(change.Removed) != 0) {
			table.Append([]string{
				kind,
				strings.Join(change.Added, "\n"),
				strings.Join(change.Removed, "\n"),
				strings.Join(change.Same, "\n"),
			})
			logrus.Debugf("kind %s; added: %+v, removed: %+v, same: %+v\n", kind, change.Added, change.Removed, change.Same)
		}
	}
	table.Render()
	return tableString.String()
}

type ResourcesTable struct {
	Version string
	Kinds   map[string][]string
	Headers []string
	Rows    [][]string
}

func (r *ResourcesTable) SortedKinds() []string {
	return slice.Sort(maps.Keys(r.Kinds))
}

func (r *ResourcesTable) SimpleTable() string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetHeader([]string{"Kind", "API Version"})
	for _, kind := range r.SortedKinds() {
		versions := r.Kinds[kind]
		sort.Strings(versions)
		for _, apiVersion := range versions {
			table.Append([]string{kind, apiVersion})
		}
	}
	table.Render()
	return tableString.String()
}

func NewResourcesTable(version string, headers []string, rows [][]string) (*ResourcesTable, error) {
	if !reflect.DeepEqual(headers, []string{"NAME", "SHORTNAMES", "APIVERSION", "NAMESPACED", "KIND", "VERBS"}) {
		return nil, errors.Errorf("invalid headers: %+v", headers)
	}
	table := &ResourcesTable{Version: version, Kinds: map[string][]string{}, Headers: headers, Rows: rows}
	for _, row := range rows {
		table.Kinds[row[4]] = append(table.Kinds[row[4]], row[2])
	}
	return table, nil
}

func (r *ResourcesTable) Diff(other *ResourcesTable) *ResourceDiff {
	changed := map[string]*diff.MapDiff{}

	for ak, av := range r.Kinds {
		bv, ok := other.Kinds[ak]
		if ok {
			changed[ak] = diff.SliceDiff(av, bv)
		} else {
			changed[ak] = &diff.MapDiff{Removed: av}
		}
	}
	for bk, bv := range other.Kinds {
		if _, ok := r.Kinds[bk]; !ok {
			changed[bk] = &diff.MapDiff{Added: bv}
		}
	}
	return &ResourceDiff{
		Changed: changed,
	}
}

func (r *ResourcesTable) KindResourcesTable() string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetHeader(r.Headers)
	table.AppendBulk(r.Rows)
	table.Render()
	return tableString.String()
}
