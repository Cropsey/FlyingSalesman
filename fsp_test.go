package fsp

import (
	"testing"
)

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
				FlightStatistics{},
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
				FlightStatistics{},
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
				FlightStatistics{},
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
				FlightStatistics{},
			},
			NewSolution(
				[]Flight{
					{0, 1, 0, 0},
					{1, 2, 1, 0},
					{2, 0, 2, 0},
				}),
		},
	}
	for _, engine := range engines_all {
		for _, test := range tests {
			comm, cm := initComm(len(test.solution.flights))
			go engine.Solve(comm, test.problem)
			s := waitForSolution(cm)
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

type commMaster struct {
	update      chan update
	queryBest   chan int
	receiveBest chan Money
	searchedAll chan int
}

func waitForSolution(cm commMaster) Solution {
	for {
		select {
		case u := <-cm.update:
			return u.s
		case <-cm.queryBest:
			cm.receiveBest <- 0
		case <-cm.searchedAll:
			return *new(Solution)
		}
	}
}

func initComm(i int) (comm, commMaster) {
	cm := commMaster{
		make(chan update, 1),
		make(chan int, 1),
		make(chan Money, 1),
		make(chan int, 1),
	}
	comm := &solutionComm{
		cm.update,
		cm.queryBest,
		cm.receiveBest,
		cm.searchedAll,
		0,
	}
	return comm, cm
}
