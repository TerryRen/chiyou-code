package str

import (
	"testing"
)

func TestUnderscoreToLowerCamel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "case-1",
			args: args{
				s: "abc_xyz",
			},
			want: "abcXyz",
		},
		{
			name: "case-2",
			args: args{
				s: "abc_xyz_ups",
			},
			want: "abcXyzUps",
		},
		{
			name: "case-3",
			args: args{
				s: "abcXyzUps",
			},
			want: "abcXyzUps",
		},
		{
			name: "case-4",
			args: args{
				s: "A",
			},
			want: "a",
		},
		{
			name: "case-5",
			args: args{
				s: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnderscoreToLowerCamel(tt.args.s); got != tt.want {
				t.Errorf("UnderscoreToLowerCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnderscoreToCamel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "case-1",
			args: args{
				s: "abc_xyz",
			},
			want: "AbcXyz",
		},
		{
			name: "case-2",
			args: args{
				s: "abc_xyz_ups",
			},
			want: "AbcXyzUps",
		},
		{
			name: "case-4",
			args: args{
				s: "a",
			},
			want: "A",
		},
		{
			name: "case-5",
			args: args{
				s: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnderscoreToCapitalizeCamel(tt.args.s); got != tt.want {
				t.Errorf("UnderscoreToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
