package tablediff

import (
	"encoding/csv"
	"fmt"
	"io"
)

// Differences contains the differences between two tables
type Differences struct {
	ExtraRows     int
	ExtraColumns  int
	Modifications [][]string
}

func (diffs *Differences) Write(w io.Writer) error {
	writer := csv.NewWriter(w)
	return writer.WriteAll(diffs.Modifications)
}

// Diff returns the Differences in table2 relative to table1 and whether the two tables are equal.
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
	mods := make([][]string, maxRows)
	for i := 0; i < maxRows; i++ {
		mods[i] = make([]string, maxCols)
		for k := 0; k < maxCols; k++ {
			var val string
			notInTable1 := len(table1) <= i || nCols1 <= k
			notInTable2 := len(table2) <= i || nCols2 <= k
			if notInTable1 && notInTable2 {
				val = "n/a"
				equal = false
			} else if notInTable1 {
				val = fmt.Sprintf("''->%v", table2[i][k])
				equal = false
			} else if notInTable2 {
				val = fmt.Sprintf("%v->''", table1[i][k])
				equal = false
			} else {
				if table1[i][k] != table2[i][k] {
					val = fmt.Sprintf("%v->%v", table1[i][k], table2[i][k])
					equal = false
				}
			}
			mods[i][k] = val
		}
	}
	if equal {
		return nil, true
	}
	ret := &Differences{
		ExtraRows:     extraRows,
		ExtraColumns:  extraColumns,
		Modifications: mods,
	}
	return ret, equal
}
