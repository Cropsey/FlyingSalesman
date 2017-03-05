package fsp

import "testing"

var engines_all = []Engine{
	dunno{},
	One{},
}

func solutionsEqual(a, b Solution) bool {
	if a.totalCost != b.totalCost {
		return false
	}
	for i, _ := range a.flights {
		if !equal(a.flights[i], b.flights[i]) {
			return false
		}
	}
	return true
}

func TestSanity(t *testing.T) {
	tests := []struct {
		description string
		engines     []Engine
		problem     Problem
		solution    Solution
	}{
		{
			"empty problem",
			engines_all,
			Problem{
				[]Flight{},
				"",
			},
			Solution{},
		},
		{
			"simple return route",
			[]Engine{One},
			Problem{
				[]Flight{
					{"brq", "lon", 1, 0},
					{"lon", "brq", 2, 0},
				},
				"brq",
			},
			NewSolution(
				[]Flight{
					{"brq", "lon", 1, 0},
					{"lon", "brq", 2, 0},
				}),
		},
		{
			"route with three stops",
			[]Engine{One},
			Problem{
				[]Flight{
					{"brq", "lon", 1, 0},
					{"lon", "xxx", 2, 0},
					{"xxx", "brq", 3, 0},
				},
				"brq",
			},
			NewSolution(
				[]Flight{
					{"brq", "lon", 1, 0},
					{"lon", "xxx", 2, 0},
					{"xxx", "brq", 3, 0},
				}),
		},
		{
			"route with three stops not in order and more flights",
			[]Engine{One},
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
				"brq",
			},
			NewSolution(
				[]Flight{
					{"lon", "xxx", 2, 0}, // 3
					{"xxx", "brq", 3, 0}, // 7
					{"brq", "lon", 1, 0}, // 11
				}),
		},
	}
	done := make(<-chan struct{})
	for _, test := range tests {
		for _, engine := range test.engines {
			ch := engine.Solve(done, test.problem)
			s := <-ch
			if !solutionsEqual(s, test.solution) {
				t.Errorf("%v: expected '%v', got '%v'",
					test.description,
					test.solution,
					s)
			}
		}
	}
}
