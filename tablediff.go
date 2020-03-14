package tablediff

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

// Differences contains the differences between two tables
type Differences struct {
	ExtraRows    int
	ExtraColumns int
	TableDiffs   [][]string
	Diffs        string
}

func (diffs *Differences) String() string {
	return diffs.Diffs
}

// AsTable returns TableDiffs as an ASCII table.
func (diffs *Differences) AsTable() string {
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.AppendBulk(diffs.TableDiffs)
	table.Render()
	return buf.String()
}

func (diffs *Differences) Write(w io.Writer) error {
	writer := csv.NewWriter(w)
	return writer.WriteAll(diffs.TableDiffs)
}

// Diff returns the Differences in table2 relative to table1 and whether the two tables are equal.
// The tables are read assuming the major dimension is rows, so [1][3] refers to row 1, column 3 (zero-indexed).
func Diff(table1 [][]string, table2 [][]string) (diffs *Differences, equal bool) {
	equal = true
	// check for nil table
	var nCols1, nCols2 int
	if len(table1) != 0 {
		nCols1 = len(table1[0])
	}
	if len(table2) != 0 {
		nCols2 = len(table2[0])
	}
	extraRows := len(table2) - len(table1)
	extraColumns := nCols2 - nCols1

	// determine max rows
	maxRows := len(table1)
	if len(table2) > len(table1) {
		maxRows = len(table2)
	}
	// determine max columns
	maxCols := nCols1
	if nCols2 > nCols1 {
		maxCols = nCols2
	}
	var diffString string
	tableDiffs := make([][]string, maxRows)
	for i := 0; i < maxRows; i++ {
		tableDiffs[i] = make([]string, maxCols)
		for k := 0; k < maxCols; k++ {
			var val string
			notInTable1 := len(table1) <= i || nCols1 <= k
			notInTable2 := len(table2) <= i || nCols2 <= k
			// not in either table (due to combining different dimensions)
			if notInTable1 && notInTable2 {
				val = "n/a"
				equal = false
				// added relative to table 1
			} else if notInTable1 {
				val = fmt.Sprintf("''->%v", table2[i][k])
				equal = false
				diffString += fmt.Sprintf("added: [%d][%d] = %v\n", i, k, table2[i][k])
				// removed relative to table 1
			} else if notInTable2 {
				val = fmt.Sprintf("%v->''", table1[i][k])
				equal = false
				diffString += fmt.Sprintf("removed: [%d][%d] (previously = %v)\n", i, k, table1[i][k])
				// modified from table 1 to table 2
			} else if table1[i][k] != table2[i][k] {
				val = fmt.Sprintf("%v -> %v", table1[i][k], table2[i][k])
				equal = false
				diffString += fmt.Sprintf("modified: [%d][%d] = %v -> %v\n", i, k, table1[i][k], table2[i][k])
			}
			tableDiffs[i][k] = val
		}
	}
	if equal {
		return nil, true
	}
	ret := &Differences{
		ExtraRows:    extraRows,
		ExtraColumns: extraColumns,
		TableDiffs:   tableDiffs,
		Diffs:        diffString,
	}
	return ret, equal
}
