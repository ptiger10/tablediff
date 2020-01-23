package tablediff

import (
	"reflect"
	"testing"
)

func TestEqual(t *testing.T) {
	type args struct {
		csv1 [][]string
		csv2 [][]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"true", args{
			[][]string{{"foo", "bar"}, {"qux", "quz"}},
			[][]string{{"foo", "bar"}, {"qux", "quz"}}},
			true,
		},
		{"true - both empty", args{
			[][]string{},
			[][]string{}},
			true,
		},
		{"false - different row count", args{
			[][]string{{"foo"}, {"qux"}},
			[][]string{{"foo"}}},
			false,
		},
		{"false - different column count", args{
			[][]string{{"foo", "bar"}, {"qux", "quz"}},
			[][]string{{"foo"}, {"qux"}}},
			false,
		},
		{"false - reordered values", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"bar", "foo"}}},
			false,
		},
		{"false - different values", args{
			[][]string{{"foo"}},
			[][]string{{"bar"}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equal(tt.args.csv1, tt.args.csv2); got != tt.want {
				t.Errorf("equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_diff(t *testing.T) {
	type args struct {
		table1 [][]string
		table2 [][]string
	}
	tests := []struct {
		name string
		args args
		want *Differences
	}{
		{"no differences", args{
			[][]string{{"foo", "bar", "baz"}, {"qux", "quz", "quuz"}},
			[][]string{{"foo", "bar", "baz"}, {"qux", "quz", "quuz"}}},
			&Differences{0, 0, nil},
		},
		{"no differences - both empty", args{
			[][]string{},
			[][]string{}},
			&Differences{0, 0, nil},
		},
		{"1 more row", args{
			[][]string{{"foo"}},
			[][]string{{"foo"}, {"baz"}}},
			&Differences{-1, 0, nil},
		},
		{"1 fewer row", args{
			[][]string{{"foo"}, {"baz"}},
			[][]string{{"foo"}}},
			&Differences{-1, 0, nil},
		},
		{"1 more column", args{
			[][]string{{"foo"}, {"baz"}},
			[][]string{{"foo", "bar"}, {"baz", "qux"}}},
			&Differences{0, 1, nil},
		},
		{"1 fewer column", args{
			[][]string{{"foo", "bar"}, {"baz", "qux"}},
			[][]string{{"foo"}, {"baz"}}},
			&Differences{0, -1, nil},
		},
		{"different values", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo", "baz"}}},
			&Differences{0, 0, []string{"[0,1]: bar -> baz"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := diff(tt.args.table1, tt.args.table2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
