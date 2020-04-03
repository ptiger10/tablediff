package tablediff

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_diff(t *testing.T) {
	type args struct {
		got  [][]string
		want [][]string
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
				TableDiffs: [][]string{{""}, {"got '', want baz"}},
				Diffs:      "[1][0]: got '', want baz\n"},
			false,
		},
		{"1 fewer row", args{
			[][]string{{"foo"}, {"baz"}},
			[][]string{{"foo"}}},
			&Differences{
				TableDiffs: [][]string{{""}, {"got baz, want ''"}},
				Diffs:      "[1][0]: got baz, want ''\n"},
			false,
		},
		{"1 more column", args{
			[][]string{{"foo"}},
			[][]string{{"foo", "bar"}}},
			&Differences{
				TableDiffs: [][]string{{"", "got '', want bar"}},
				Diffs:      "[0][1]: got '', want bar\n"},
			false,
		},
		{"1 fewer column", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo"}}},
			&Differences{
				TableDiffs: [][]string{{"", "got bar, want ''"}},
				Diffs:      "[0][1]: got bar, want ''\n"},
			false,
		},
		{"different values", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo", "baz"}}},
			&Differences{
				TableDiffs: [][]string{{"", "got bar, want baz"}},
				Diffs:      "[0][1]: got bar, want baz\n"},
			false,
		},
		{"different shapes", args{
			[][]string{{"foo", "bar"}},
			[][]string{{"foo"}, {"baz"}}},
			&Differences{
				TableDiffs: [][]string{{"", "got bar, want ''"}, {"got '', want baz", "n/a"}},
				Diffs:      "[0][1]: got bar, want ''\n[1][0]: got '', want baz\n"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, equal := Diff(tt.args.got, tt.args.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
			if equal != tt.wantEqual {
				t.Errorf("Diff() equal = %v, want %v", equal, tt.wantEqual)
			}
		})
	}
}

func TestDifferences_WriteCSV(t *testing.T) {
	type fields struct {
		TableDiffs [][]string
		Diffs      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantW   string
		wantErr bool
	}{
		{"pass", fields{
			TableDiffs: [][]string{{"got foo, want bar", "got bar, want ''"}, {"got '', want baz", "n/a"}}},
			"" +
				"got foo, want bar|got bar, want ''\n" +
				"got '', want baz|n/a\n", false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffs := &Differences{
				TableDiffs: tt.fields.TableDiffs,
				Diffs:      tt.fields.Diffs,
			}
			w := &bytes.Buffer{}
			if err := diffs.WriteCSV(w); (err != nil) != tt.wantErr {
				t.Errorf("Differences.WriteCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Differences.WriteCSV() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
