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
	Modifications []string
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
	var mods []string
	for i := 0; i < len(table1); i++ {
		if len(table2) <= i {
			continue
		} else {
			for j := 0; j < nCols1; j++ {
				if nCols2 <= j {
					continue
				} else {
					if table1[i][j] != table2[i][j] {
						mods = append(mods, fmt.Sprintf("[%v,%v]: %v -> %v", i, j, table1[i][j], table2[i][j]))
					}
				}
			}
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
