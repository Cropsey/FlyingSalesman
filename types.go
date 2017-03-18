package fsp

import (
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
	From      City
	To        City
	Day       Day
	Cost      Money
	Heuristic Money
    Penalty   float64 
}

type FlightStats struct {
	FlightCount uint16
	BestPrice   Money
	BestDay     Day
	BestDest    City
	AvgPrice    float32
}

type FlightStatistics struct {
	ByDest [][]FlightStats
	ByDay  [][]FlightStats
}

type Problem struct {
	flights []Flight
	start   City
	n       int //size = number of cities/days
	//stats   [][]FlightStats
	stats FlightStatistics
}

func (p Problem) Solve(timeout <-chan time.Time) (Solution, error) {
    sol, err := kickTheEngines(p, timeout)
    /*for _, f := range p.flights {
        if f.Penalty != 0 {
            printInfo(f)
        }
    }*/
    return sol, err
}

func (p Problem) FlightsCnt() int {
	return len(p.flights)
}

func (p Problem) FlightStats() FlightStatistics {
	return p.stats
}

func (p Problem) CitiesCnt() int {
	return p.n
}

func NewProblem(flights []Flight, n int, stats FlightStatistics) Problem {
	return Problem{flights, 0, n, stats}
}

type Solution struct {
	flights   []Flight
	totalCost Money
}

func (s Solution) GetFlights() []Flight {
	return s.flights
}

func (s Solution) GetTotalCost() Money {
	return s.totalCost
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
	return f[i].Day < f[j].Day
}
