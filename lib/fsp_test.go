package fsp

import "testing"

var engines_all = []FspEngine{
	dunno{},
}

func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestAll(t *testing.T) {
	tests := []struct {
		description string
		engines     []FspEngine
		problem     Problem
		solution    []int
	}{
		{
			"empty",
			engines_all,
			Problem{
				&[]Flight{},
				&[]string{},
			},
			[]int{},
		},
	}
	for _, test := range tests {
		for _, engine := range test.engines {
			s := engine.Solve(test.problem)
			if !intSliceEqual(s, test.solution) {
				t.Error("%v: expected '%v', got '%v'",
					test.description,
					test.solution,
					s)
			}
		}
	}
}
