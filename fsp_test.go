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
		problem     Problem
		solution    Solution
	}{
		{
			"empty problem",
			Problem{
				[]Flight{},
				0,
				0,
			},
			Solution{},
		},
		{
			"simple return route",
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
			"route with three stops not in order",
			Problem{
				[]Flight{
					{2, 0, 2, 0},
					{1, 2, 1, 0},
					{0, 1, 0, 0},
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
	}
	done := make(<-chan struct{})
	for _, engine := range engines_all {
		for _, test := range tests {
			ch := engine.Solve(done, test.problem)
			s := <-ch
			if !solutionsEqual(s, test.solution) {
				t.Errorf("Engine %v, test %v: expected '%v', got '%v'",
					engine.Name(),
					test.description,
					test.solution,
					s)
			}
		}
	}
}
