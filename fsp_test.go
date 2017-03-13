package fsp

import (
	"math"
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
					{0, 1, 0, 0, 0},
					{1, 0, 1, 0, 0},
				},
				0,
				2,
				FlightStatistics{},
			},
			NewSolution(
				[]Flight{
					{0, 1, 0, 0, 0},
					{1, 0, 1, 0, 0},
				}),
		},
		{
			"route with three stops",
			Problem{
				[]Flight{
					{0, 1, 0, 0, 0},
					{1, 2, 1, 0, 0},
					{2, 0, 2, 0, 0},
				},
				0,
				3,
				FlightStatistics{},
			},
			NewSolution(
				[]Flight{
					{0, 1, 0, 0, 0},
					{1, 2, 1, 0, 0},
					{2, 0, 2, 0, 0},
				}),
		},
		{
			"route with three stops not in order",
			Problem{
				[]Flight{
					{2, 0, 2, 0, 0},
					{1, 2, 1, 0, 0},
					{0, 1, 0, 0, 0},
				},
				0,
				3,
				FlightStatistics{},
			},
			NewSolution(
				[]Flight{
					{0, 1, 0, 0, 0},
					{1, 2, 1, 0, 0},
					{2, 0, 2, 0, 0},
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
	buffer      *Solution
	bufferFree  chan bool
	bufferReady chan int
	queryBest   chan int
	receiveBest chan Money
	searchedAll chan int
}

func waitForSolution(cm commMaster) Solution {
	cm.bufferFree <- true
	for {
		select {
		case <-cm.bufferReady:
			return *cm.buffer
		case <-cm.queryBest:
			cm.receiveBest <- 0
		case <-cm.searchedAll:
			return *cm.buffer
		}
	}
}

func initComm(i int) (comm, commMaster) {
	cm := commMaster{
		&Solution{make([]Flight, i), math.MaxInt32},
		make(chan bool, 1),
		make(chan int, 1),
		make(chan int, 1),
		make(chan Money, 1),
		make(chan int, 1),
	}
	comm := &bufferComm{
		cm.buffer,
		cm.bufferFree,
		cm.bufferReady,
		cm.queryBest,
		cm.receiveBest,
		cm.searchedAll,
		0,
	}
	return comm, cm
}
