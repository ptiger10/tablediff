package tablediff

import (
	"fmt"
)

func ExampleDifferences() {
	table1 := [][]string{{"foo", "bar", "qux"}}
	table2 := [][]string{{"foo", "baz", "qux"}}
	diffs, ok := Diff(table1, table2)
	if !ok {
		fmt.Println(diffs)
	}
	// Output:
	// [0][1]: got bar, want baz
}

func ExampleDifferences_AsTable() {
	table1 := [][]string{{"foo", "bar", "qux"}}
	table2 := [][]string{{"foo", "baz", "qux"}}
	diffs, ok := Diff(table1, table2)
	if !ok {
		fmt.Println(diffs.AsTable())
	}
	// Output:
	// +--+-------------------+--+
	// |  | got bar, want baz |  |
	// +--+-------------------+--+
}
