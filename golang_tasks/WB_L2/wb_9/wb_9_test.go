package main

import (
	"testing"
)

func TestUnpackingString(t *testing.T) {
	tests := []struct {
		input       string
		want        string
		expectError bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", true},
		{`qwe\4\5`, "qwe45", false},
		{`qwe\45`, "qwe44444", false},
	}

	for _, tt := range tests {
		got, err := UnpackingString(tt.input)
		if (err != nil) != tt.expectError {
			t.Errorf("UnpackingString(%q) error = %v, wantErr %v", tt.input, err, tt.expectError)
		}
		if got != tt.want {
			t.Errorf("UnpackingString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
