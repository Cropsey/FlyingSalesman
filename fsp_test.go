package fsp

import "testing"

var engines_all = []FspEngine{
	dunno{},
	one_ordered{},
	one_places{},
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

func TestSanity(t *testing.T) {
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
			[]FspEngine{one_ordered{}, one_places{}},
			Problem{
				[]Flight{
					{"brq", "lon", 1, 0},
					{"lon", "brq", 2, 0},
				},
				[]string{"brq","lon"},
			},
			[]int{0,1},
		},
		{
			"route with three stops",
			[]FspEngine{one_ordered{}, one_places{}},
			Problem{
				[]Flight{
					{"brq", "lon", 1, 0},
					{"lon", "xxx", 2, 0},
					{"xxx", "brq", 3, 0},
				},
				[]string{"brq","lon", "xxx"},
			},
			[]int{0,1,2},
		},
		{
			"route with three stops not in order and more flights",
			[]FspEngine{one_ordered{}, one_places{}},
			Problem{
				[]Flight{
					{"aaa", "bbb", 1, 0},
					{"aaa", "bbb", 2, 0},
					{"aaa", "bbb", 3, 0},
					{"lon", "xxx", 2, 0}, // 3
					{"bbb", "ccc", 1, 0},
					{"bbb", "ccc", 2, 0},
					{"bbb", "ccc", 3, 0},
					{"xxx", "brq", 3, 0}, // 7
					{"ccc", "aaa", 1, 0},
					{"ccc", "aaa", 2, 0},
					{"ccc", "aaa", 3, 0},
					{"brq", "lon", 1, 0}, // 11
				},
				[]string{"brq","lon", "xxx"},
			},
			[]int{11,3,7},
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
