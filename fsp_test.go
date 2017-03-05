package fsp

import "testing"

var engines_all = []Engine{
	One{},
	Mitm{},
}

func solutionsEqual(a, b Solution) bool {
	if a.totalCost != b.totalCost {
		return false
	}
	if len(a.flights) != len(b.flights) {
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
				0,
				0,
			},
			Solution{},
		},
		{
			"simple return route",
			engines_all,
			Problem{
				[]Flight{
					{0, 1, 0, 0},
					{1, 0, 1, 0},
				},
				0,
				2,
			},
			NewSolution(
				[]Flight{
					{0, 1, 0, 0},
					{1, 0, 1, 0},
				}),
		},
		{
			"route with three stops",
			engines_all,
			Problem{
				[]Flight{
					{0, 1, 0, 0},
					{1, 2, 1, 0},
					{2, 0, 2, 0},
				},
				0,
				3,
			},
			NewSolution(
				[]Flight{
					{0, 1, 0, 0},
					{1, 2, 1, 0},
					{2, 0, 2, 0},
				}),
		},
		{
			"route with three stops not in order and more flights",
			engines_all,
			Problem{
				[]Flight{
					{2, 3, 0, 0},
					{2, 3, 1, 0},
					{2, 3, 2, 0},
					{1, 5, 1, 0}, // 3
					{3, 4, 0, 0},
					{3, 4, 1, 0},
					{3, 4, 2, 0},
					{5, 0, 2, 0}, // 7
					{4, 2, 0, 0},
					{4, 2, 1, 0},
					{4, 2, 2, 0},
					{0, 1, 0, 0}, // 11
				},
				0,
				6,
			},
			NewSolution(
				[]Flight{
					{1, 5, 1, 0}, // 3
					{5, 0, 2, 0}, // 7
					{0, 1, 0, 0}, // 11
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
