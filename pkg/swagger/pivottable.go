package swagger

import (
	"github.com/mattfenwick/collections/pkg/set"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"strings"
)

type PivotTable struct {
	FirstColumnHeader string
	Rows              map[string]map[string][]string
	Columns           []string
	columnSet         *set.Set[string]
}

func NewPivotTable(firstColumn string, restColumns []string) *PivotTable {
	columnSet := set.NewSet(restColumns)
	if len(restColumns) != columnSet.Len() {
		panic(errors.Errorf("expected unique columns, found duplicate in %+v", restColumns))
	}
	return &PivotTable{
		FirstColumnHeader: firstColumn,
		Rows:              map[string]map[string][]string{},
		Columns:           restColumns,
		columnSet:         columnSet,
	}
}

func (e *PivotTable) Add(rowKey string, columnKey string, value string) {
	if !e.columnSet.Contains(columnKey) {
		panic(errors.Errorf("invalid column name %s, not found in %+v", columnKey, e.Columns))
	}
	if _, ok := e.Rows[rowKey]; !ok {
		e.Rows[rowKey] = map[string][]string{}
	}
	e.Rows[rowKey][columnKey] = append(e.Rows[rowKey][columnKey], value)
}

func (e *PivotTable) FormattedTable(formatRow func(rowKey string, values [][]string) []string) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)

	table.SetHeader(append([]string{e.FirstColumnHeader}, e.Columns...))
	for _, rowKey := range slice.Sort(maps.Keys(e.Rows)) {
		values := slice.Map(func(c string) []string { return e.Rows[rowKey][c] }, e.Columns)
		table.Append(formatRow(rowKey, values))
	}

	table.Render()
	return tableString.String()
}
