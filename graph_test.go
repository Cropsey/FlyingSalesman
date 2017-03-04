package fsp

import "testing"

func check(problem Problem, expected []Flight, t *testing.T) {
	graph := NewGraph(problem)
	filtered := graph.Filtered()
	if len(filtered) != len(expected) {
		t.Errorf("Expected %d, filtered %d", len(expected), len(filtered))
	}
	for _, e := range expected {
		found := false
		for _, f := range filtered {
			if equal(e, f) {
				found = true
			}
		}
		if !found {
			t.Error("Unable to find expected flight", e)
		}
	}
}

func TestOneDupl(t *testing.T) {
	problem := NewProblem(
		[]Flight{
			{0, 1, 1, 100},
			{1, 0, 1, 100},
			{1, 0, 1, 200},
		},
		[]string{"brq", "lon"},
	)
	expect := []Flight{
		{0, 1, 1, 100},
		{1, 0, 1, 100},
	}
	check(problem, expect, t)
}

func TestNoFilter(t *testing.T) {
	problem := NewProblem(
		[]Flight{
			{0, 1, 1, 0},
			{1, 2, 2, 0},
			{2, 0, 3, 0},
		},
		[]string{"brq", "lon", "xxx"},
	)
	expect := []Flight{
		{0, 1, 1, 0},
		{1, 2, 2, 0},
		{2, 0, 3, 0},
	}
	check(problem, expect, t)
}

func TestOneNewGraph(t *testing.T) {
	problem := NewProblem(
		[]Flight{
			{0, 1, 1, 900},
			{1, 2, 2, 600},
			{1, 2, 2, 400},
			{2, 0, 3, 800},
		},
		[]string{"brq", "lon", "xxx"},
	)
	expect := []Flight{
		{0, 1, 1, 900},
		{1, 2, 2, 400},
		{2, 0, 3, 800},
	}
	check(problem, expect, t)
}

func TestMultipleFiler(t *testing.T) {
	problem := NewProblem(
		[]Flight{
			{0, 1, 1, 700},
			{0, 1, 1, 1000},
			{0, 1, 1, 300},
			{1, 2, 2, 600},
			{1, 2, 2, 400},
			{1, 2, 2, 400},
			{1, 2, 2, 200},
			{1, 0, 2, 100},
			{1, 0, 2, 200},
			{1, 0, 3, 101},
			{1, 0, 3, 201},
			{2, 0, 3, 100},
			{2, 0, 3, 900},
		},
		[]string{"brq", "lon", "xxx"},
	)
	expect := []Flight{
		{0, 1, 1, 300},
		{1, 2, 2, 200},
		{1, 0, 2, 100},
		{1, 0, 3, 101},
		{2, 0, 3, 100},
	}
	check(problem, expect, t)
}
