package fsp

import "testing"

func equal(a, b Flight) bool {
	return a.from == b.from && a.to == b.to && a.day == b.day && a.cost == b.cost
}

func check(problem Problem, expected []Flight, t *testing.T) {
	filtered := NewGraph(problem).Filtered()
	if len(filtered) != len(expected) {
		t.Errorf("Expected %d, filtered %d", len(filtered), len(expected))
	}
	for _, e := range expected {
		found := false
		for _, f := range filtered {
			if equal(e, f) {
				found = true
			}
		}
		if !found {
			t.Error("Unable to found expected flight", e)
		}
	}
}

func TestOneDupl(t *testing.T) {
	problem := Problem{
		[]Flight{
			{"xxx", "lon", 1, 100},
			{"brq", "lon", 1, 100},
			{"brq", "lon", 1, 200},
		},
		[]string{"brq", "lon"},
	}
	expect := []Flight{
		{"xxx", "lon", 1, 100},
		{"brq", "lon", 1, 100},
	}
	check(problem, expect, t)
}

func TestNoFilter(t *testing.T) {
	problem := Problem{
		[]Flight{
			{"brq", "lon", 1, 0},
			{"lon", "xxx", 2, 0},
			{"xxx", "brq", 3, 0},
		},
		[]string{"brq", "lon", "xxx"},
	}
	expect := []Flight{
		{"brq", "lon", 1, 0},
		{"lon", "xxx", 2, 0},
		{"xxx", "brq", 3, 0},
	}
	check(problem, expect, t)
}

func TestOneNewGraph(t *testing.T) {
	problem := Problem{
		[]Flight{
			{"brq", "lon", 1, 900},
			{"lon", "xxx", 2, 600},
			{"lon", "xxx", 2, 400},
			{"xxx", "brq", 3, 800},
		},
		[]string{"brq", "lon", "xxx"},
	}
	expect := []Flight{
		{"brq", "lon", 1, 900},
		{"lon", "xxx", 2, 400},
		{"xxx", "brq", 3, 800},
	}
	check(problem, expect, t)
}

func TestMultipleFiler(t *testing.T) {
	problem := Problem{
		[]Flight{
			{"brq", "lon", 1, 700},
			{"brq", "lon", 1, 1000},
			{"brq", "lon", 1, 300},
			{"lon", "xxx", 2, 600},
			{"lon", "xxx", 2, 400},
			{"lon", "xxx", 2, 400},
			{"lon", "xxx", 2, 200},
			{"lon", "brq", 2, 100},
			{"lon", "brq", 2, 200},
			{"lon", "brq", 3, 101},
			{"lon", "brq", 3, 201},
			{"xxx", "brq", 3, 100},
			{"xxx", "brq", 3, 900},
		},
		[]string{"brq", "lon", "xxx"},
	}
	expect := []Flight{
		{"brq", "lon", 1, 300},
		{"lon", "xxx", 2, 200},
		{"lon", "brq", 2, 100},
		{"lon", "brq", 3, 101},
		{"xxx", "brq", 3, 100},
	}
	check(problem, expect, t)
}
