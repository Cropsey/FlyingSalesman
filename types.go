package fsp

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

type City uint32
type Money uint32

func (m Money) String() string {
	return fmt.Sprintf("%d", m)
}

type Day uint16

type Flight struct {
	from City
	to   City
	day  Day
	cost Money
}

func NewFlight(from, to uint32, day uint16, cost int) Flight {
	flight := new(Flight)
	flight.from = City(from)
	flight.to = City(to)
	flight.day = Day(day)
	flight.cost = Money(cost)
	return *flight
}

type Problem struct {
	flights []Flight
	start   City
	cities  []string
}

type taskData struct {
	graph   Graph
	problem Problem
	timeout <-chan time.Time
}

func (p Problem) Solve(timeout <-chan time.Time) (Solution, error) {
	graph := NewGraph(p)
	task := &taskData{graph, p, timeout}
	if len(p.flights) > 2 {
		return KickTheEngines(task)
	} else {
		return DFS(task)
	}
}

func NewProblem(flights []Flight, cities []string) Problem {
	return Problem{flights, 0, cities[:]}
}

type Solution struct {
	flights   []Flight
	totalCost Money
	cities    []string
}

func NewSolution(flights []Flight, cities []string) Solution {
	sort.Sort(ByDay(flights))
	return Solution{flights, Cost(flights), cities}
}

type ByDay []Flight

func (f ByDay) Len() int {
	return len(f)
}
func (f ByDay) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f ByDay) Less(i, j int) bool {
	return f[i].day < f[j].day
}

func (s Solution) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(s.totalCost.String())
	buffer.WriteString("\n")
	for _, f := range s.flights {
		from := s.cities[f.from]
		to := s.cities[f.to]
		flight := fmt.Sprintf("%s %s %d %d\n", from, to, f.day, f.cost)
		buffer.WriteString(flight)
	}
	return buffer.String()
}

type Engine interface {
	Solve(done <-chan struct{}, problem Problem) <-chan Solution
}
