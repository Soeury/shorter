package base62

import (
	"testing"
)

func TestChangeToBase62(t *testing.T) {

	tests := []struct {
		name string
		seq  uint64
		want string
	}{
		{
			name: "case:1",
			seq:  1,
			want: "1",
		},

		{
			name: "case:36",
			seq:  36,
			want: "A",
		},

		{
			name: "case:62",
			seq:  62,
			want: "10",
		},

		{
			name: "case:6347",
			seq:  6347,
			want: "1En",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChangeToBase62(tt.seq); got != tt.want {
				t.Errorf("ChangeToBase62() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeToBase10(t *testing.T) {

	tests := []struct {
		name  string
		str62 string
		want  uint64
	}{
		{
			name:  "case:a",
			str62: "a",
			want:  10,
		},

		{
			name:  "case:A",
			str62: "A",
			want:  36,
		},

		{
			name:  "case:10",
			str62: "10",
			want:  62,
		},

		{
			name:  "case:1En",
			str62: "1En",
			want:  6347,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChangeToBase10(tt.str62); got != tt.want {
				t.Errorf("ChangeToBase10() = %v, want %v", got, tt.want)
			}
		})
	}
}
