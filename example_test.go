package tablediff

import (
	"fmt"
)

func ExampleDifferences() {
	diffs := &Differences{
		ExtraRows:     1,
		ExtraColumns:  0,
		Modifications: [][]string{{"foo -> bar"}},
	}
	fmt.Println(diffs)
	// Output:
	// +------------+
	// | foo -> bar |
	// +------------+
}
