package fsp

import "testing"

var engines_all = []FspEngine{
	dunno{},
	one_ordered{},
}

func intSlicesEqual(a, b []int) bool {
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
			"empty problem",
			engines_all,
			Problem{
				[]Flight{},
				[]string{},
			},
			[]int{},
		},
		{
			"simple return route",
			[]FspEngine{one_ordered{}},
			Problem{
				[]Flight{
					{"brq", "lon", 1, 0},
					{"lon", "brq", 2, 0},
				},
				[]string{"brq","lon"},
			},
			[]int{0,1},
		},
	}
	for _, test := range tests {
		for _, engine := range test.engines {
			s := engine.Solve(test.problem)
			if !intSlicesEqual(s, test.solution) {
				t.Error("%v: expected '%v', got '%v'",
					test.description,
					test.solution,
					s)
			}
		}
	}
}
