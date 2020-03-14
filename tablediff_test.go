package tablediff

import (
	"reflect"
	"testing"
)

func Test_diff(t *testing.T) {
	type args struct {
		table1 [][]string
		table2 [][]string
	}
	tests := []struct {
		name      string
		args      args
		want      *Differences
		wantEqual bool
	}{
		{"no differences", args{
			[][]string{{"foo"}, {"bar"}},
			[][]string{{"foo"}, {"bar"}}},
			nil,
			true,
		},
		{"no differences - both empty", args{
			[][]string{},
			[][]string{}},
			nil,
			true,
		},
		{"1 more row", args{
			[][]string{{"foo"}},
			[][]string{{"foo"}, {"baz"}}},
			&Differences{
				ExtraRows:    1,
				ExtraColumns: 0,
				TableDiffs:   [][]string{{""}, {"''->baz"}},
				Diffs:        "added: [1][0] = baz\n"},
			false,
		},
		{"1 fewer row", args{
			[][]string{{"foo"}, {"baz"}},
			[][]string{{"foo"}}},
			&Differences{
				ExtraRows:    -1,
				ExtraColumns: 0,
				TableDiffs:   [][]string{{""}, {"baz->''"}},
				Diffs:        "removed: [1][0] (previously = baz)\n"},
			false,
		},
		{"1 more column", args{
			[][]string{{"foo"}},
			[][]string{{"foo", "bar"}}},
			&Differences{
				ExtraRows:    0,
				ExtraColumns: 1,
				TableDiffs:   [][]string{{"", "''->bar"}},
				Diffs:        "added: [0][1] = bar\n"},
			false,
		},
		{"1 fewer column", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo"}}},
			&Differences{
				ExtraRows:    0,
				ExtraColumns: -1,
				TableDiffs:   [][]string{{"", "bar->''"}},
				Diffs:        "removed: [0][1] (previously = bar)\n"},
			false,
		},
		{"different values", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo", "baz"}}},
			&Differences{
				ExtraRows:    0,
				ExtraColumns: 0,
				TableDiffs:   [][]string{{"", "bar -> baz"}},
				Diffs:        "modified: [0][1] = bar -> baz\n"},
			false,
		},
		{"different shapes", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo"}, {"baz"}}},
			&Differences{
				ExtraRows:    1,
				ExtraColumns: -1,
				TableDiffs:   [][]string{{"", "bar->''"}, {"''->baz", "n/a"}},
				Diffs:        "removed: [0][1] (previously = bar)\nadded: [1][0] = baz\n"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, equal := Diff(tt.args.table1, tt.args.table2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %#v, want %#v", got, tt.want)
			}
			if equal != tt.wantEqual {
				t.Errorf("Diff() equal = %#v, want %#v", equal, tt.wantEqual)
			}
		})
	}
}
