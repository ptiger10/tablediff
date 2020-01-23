package tablediff

import "fmt"

// equal evaluates whether two tables are identical
func equal(table1 [][]string, table2 [][]string) bool {
	if len(table1) == 0 && len(table2) == 0 {
		return true
	}
	if len(table1) != len(table2) {
		return false
	}
	if len(table1[0]) != len(table2[0]) {
		return false
	}
	for i := 0; i < len(table1); i++ {
		for j := 0; j < len(table1[0]); j++ {
			if table1[i][j] != table2[i][j] {
				return false
			}
		}
	}
	return true
}

// Differences contains the differences between two tables
type Differences struct {
	ExtraRows     int
	ExtraColumns  int
	Modifications [][]string
}

// diff returns the Differences in table2 relative to table1.
func diff(table1 [][]string, table2 [][]string) *Differences {
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
	maxRows := len(table1)
	if len(table2) > len(table1) {
		maxRows = len(table2)
	}
	maxCols := nCols1
	if nCols2 > nCols1 {
		maxCols = nCols2
	}
	mods := make([][]string, maxRows)
	for i := 0; i < maxRows; i++ {
		mods[i] = make([]string, maxCols)
		for j := 0; j < maxCols; j++ {
			var val string
			notInTable1 := len(table1) <= i || nCols1 <= j
			notInTable2 := len(table2) <= i || nCols2 <= j
			if notInTable1 && notInTable2 {
				val = "n/a"
			} else if notInTable1 {
				val = fmt.Sprintf("'' -> %v", table2[i][j])
			} else if notInTable2 {
				val = fmt.Sprintf("%v -> ''", table1[i][j])
			} else {
				if table1[i][j] != table2[i][j] {
					val = fmt.Sprintf("%v -> %v", table1[i][j], table2[i][j])
				}
			}
			mods[i][j] = val
		}
	}
	ret := &Differences{
		ExtraRows:     extraRows,
		ExtraColumns:  extraColumns,
		Modifications: mods,
	}
	return ret
}

// Diff returns the Differences between two tables. If there are no differences, ok returns true.
func Diff(table1 [][]string, table2 [][]string) (diffs *Differences, ok bool) {
	diffs = diff(table1, table2)
	ok = equal(table1, table2)
	return
}
