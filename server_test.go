package main

import (
	"testing"
)

// TestSum tests the Sum function
func TestSum(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
		x    int
		y    int
		want int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed numbers", -2, 3, 1},
	}

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.x, tt.y); got != tt.want {
				t.Errorf("Sum(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.want)
			}
		})
	}
}
