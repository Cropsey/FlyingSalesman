package fsp

import (
	"bytes"
	"fmt"
	"sort"
)

type City string
type Money int

func (m Money) String() string {
	return fmt.Sprintf("%d", m)
}

type Day int

type Flight struct {
	from City
	to   City
	day  Day
	cost Money
}

func NewFlight(from, to string, day, cost int) Flight {
	flight := new(Flight)
	flight.from = City(from)
	flight.to = City(to)
	flight.day = Day(day)
	flight.cost = Money(cost)
	return *flight
}
func (f Flight) String() string {
	return fmt.Sprintf("%s %s %d %d\n", f.from, f.to, f.day, f.cost)
}

type Problem struct {
	flights []Flight
	start   City
}
func NewProblem(src string, flights []Flight) Problem {
    return Problem{flights, City(src)}
}

type Solution struct {
	flights   []Flight
	totalCost Money
}

func NewSolution(flights []Flight) Solution {
	sort.Sort(ByDay(flights))
	return Solution{flights, Cost(flights)}
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
		buffer.WriteString(f.String())
	}
	return buffer.String()
}

type Engine interface {
	Solve(done <-chan struct{}, problem Problem) <-chan Solution
}
