package tablediff

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/ptiger10/tablewriter"
)

// Differences contains the differences between two tables
type Differences struct {
	TableDiffs [][]string
	Diffs      string
}

func (diffs *Differences) String() string {
	return diffs.Diffs
}

// AsTable returns TableDiffs as an ASCII table.
func (diffs *Differences) AsTable() string {
	var buf bytes.Buffer
	tablewriter.ChangeDefaults(tablewriter.Defaults{MaxColWidth: 50})
	table := tablewriter.NewTable(&buf)
	table.AppendRows(diffs.TableDiffs)
	table.Render()
	return buf.String()
}

// WriteCSV writes `diffs` to `w` in CSV format.
func (diffs *Differences) WriteCSV(w io.Writer) error {
	writer := csv.NewWriter(w)
	writer.Comma = '|'
	return writer.WriteAll(diffs.TableDiffs)
}

// Diff returns whether `got` is equal to `want`, and if not, the Differences between the two.
// The tables are read assuming the major dimension is rows, so [1][3] refers to row 1, column 3 (zero-indexed).
func Diff(got [][]string, want [][]string) (diffs *Differences, equal bool) {
	equal = true
	// check for nil table
	var nCols1, nCols2 int
	if len(got) != 0 {
		nCols1 = len(got[0])
	}
	if len(want) != 0 {
		nCols2 = len(want[0])
	}

	// determine max rows
	maxRows := len(got)
	if len(want) > len(got) {
		maxRows = len(want)
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
			notInGot := len(got) <= i || nCols1 <= k
			notInWant := len(want) <= i || nCols2 <= k
			// not in either table (due to combining different dimensions)
			if notInGot && notInWant {
				val = "n/a"
				equal = false
				// in want, not in got
			} else if notInGot {
				val = fmt.Sprintf("got '', want %v", want[i][k])
				equal = false
				diffString += fmt.Sprintf("[%d][%d]: %v\n", i, k, val)
				// in got, not in want
			} else if notInWant {
				val = fmt.Sprintf("got %v, want ''", got[i][k])
				equal = false
				diffString += fmt.Sprintf("[%d][%d]: %v\n", i, k, val)
				// different in got compared to want
			} else if got[i][k] != want[i][k] {
				val = fmt.Sprintf("got %v, want %v", got[i][k], want[i][k])
				equal = false
				diffString += fmt.Sprintf("[%d][%d]: %v\n", i, k, val)
			}
			tableDiffs[i][k] = val
		}
	}
	if equal {
		return nil, true
	}
	ret := &Differences{
		TableDiffs: tableDiffs,
		Diffs:      diffString,
	}
	return ret, equal
}
