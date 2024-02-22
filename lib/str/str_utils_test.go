package str

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnderscoreToLowerCamel(tt.args.s); got != tt.want {
				t.Errorf("UnderscoreToLowerCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
