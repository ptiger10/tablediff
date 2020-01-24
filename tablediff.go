package tablediff

import "fmt"

// Differences contains the differences between two tables
type Differences struct {
	ExtraRows     int
	ExtraColumns  int
	Modifications [][]string
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
				equal = false
			} else if notInTable1 {
				val = fmt.Sprintf("'' -> %v", table2[i][j])
				equal = false
			} else if notInTable2 {
				val = fmt.Sprintf("%v -> ''", table1[i][j])
				equal = false
			} else {
				if table1[i][j] != table2[i][j] {
					val = fmt.Sprintf("%v -> %v", table1[i][j], table2[i][j])
					equal = false
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
	return ret, equal
}
