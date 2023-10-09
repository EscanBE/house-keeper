package files

import "testing"

func Test_isOrContainsRecursiveFlag(t *testing.T) {
	tests := []struct {
		option string
		want   bool
	}{
		{
			option: "recursive",
			want:   false,
		},
		{
			option: "r",
			want:   false,
		},
		{
			option: "--recursive",
			want:   true,
		},
		{
			option: "-r",
			want:   true,
		},
		{
			option: "-recursive",
			want:   true, // contains 'r' and that's enough
		},
		{
			option: "--r",
			want:   false,
		},
		{
			option: "--recursive-mixed",
			want:   false,
		},
		{
			option: "-rlptgoD",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.option, func(t *testing.T) {
			if got := isOrContainsRecursiveFlag(tt.option); got != tt.want {
				t.Errorf("isOrContainsRecursiveFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
