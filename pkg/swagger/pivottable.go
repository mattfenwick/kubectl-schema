package swagger

import (
	"strings"

	"github.com/mattfenwick/collections/pkg/set"
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

type PivotTable struct {
	FirstColumnHeader string
	Rows              map[string]map[string][]string
	Columns           []string
	columnSet         *set.Set[string]
}

func NewPivotTable(firstColumn string, restColumns []string) *PivotTable {
	columnSet := set.FromSlice(restColumns)
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

func (e *PivotTable) ToRawTable(formatRow func(rowKey string, values [][]string) []string) *RawTable {
	headers := append([]string{e.FirstColumnHeader}, e.Columns...)
	var rows [][]string

	for _, rowKey := range slice.Sort(maps.Keys(e.Rows)) {
		values := slice.Map(func(c string) []string { return e.Rows[rowKey][c] }, e.Columns)
		rows = append(rows, formatRow(rowKey, values))
	}

	return NewRawTable(headers, rows)
}

type RawTable struct {
	Headers []string
	Rows    [][]string
}

func NewRawTable(headers []string, rows [][]string) *RawTable {
	for i, r := range rows {
		if len(headers) != len(r) {
			panic(errors.Errorf("mismatch between length of headers and of row %d: %d vs. %d", i, len(headers), len(r)))
		}
	}
	return &RawTable{Headers: headers, Rows: rows}
}

func (r *RawTable) ToMarkdownTable() string {
	rows := [][]string{
		r.Headers,
		slice.Map(func(h string) string { return "---" }, r.Headers),
	}
	rows = append(rows, r.Rows...)
	lines := slice.Map(func(row []string) string {
		return "| " + strings.Join(row, " | ") + " |"
	}, rows)
	return strings.Join(lines, "\n")
}

func (r *RawTable) ToFormattedTable() string {
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
