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
			&Differences{0, 0, [][]string{{""}, {""}}},
			true,
		},
		{"no differences - both empty", args{
			[][]string{},
			[][]string{}},
			&Differences{0, 0, [][]string{}},
			true,
		},
		{"1 more row", args{
			[][]string{{"foo"}},
			[][]string{{"foo"}, {"baz"}}},
			&Differences{1, 0, [][]string{{""}, {"''->baz"}}},
			false,
		},
		{"1 fewer row", args{
			[][]string{{"foo"}, {"baz"}},
			[][]string{{"foo"}}},
			&Differences{-1, 0, [][]string{{""}, {"baz->''"}}},
			false,
		},
		{"1 more column", args{
			[][]string{{"foo"}},
			[][]string{{"foo", "bar"}}},
			&Differences{0, 1, [][]string{{"", "''->bar"}}},
			false,
		},
		{"1 fewer column", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo"}}},
			&Differences{0, -1, [][]string{{"", "bar->''"}}},
			false,
		},
		{"different values", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo", "baz"}}},
			&Differences{0, 0, [][]string{{"", "bar->baz"}}},
			false,
		},
		{"different shapes", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo"}, {"baz"}}},
			&Differences{1, -1, [][]string{{"", "bar->''"}, {"''->baz", "n/a"}}},
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
