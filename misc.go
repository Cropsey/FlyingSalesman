package fsp

import (
	"fmt"
	"os"
	"time"
)

const MAX_CITIES int = 300
const MAX_FLIGHTS int = 27000000

var BeVerbose bool
var StartTime time.Time

func Cost(flights []Flight) Money {
	var sum Money
	for _, f := range flights {
		sum += f.Cost
	}
	return sum
}

// is solution correct? if not, why?
func correct(p Problem, s Solution) (bool, string) {
	var day Day
	for _, f := range s.flights {
		if day > f.Day {
			return false, "timing"
		}
		day = f.Day + 1
	}
	return true, ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func stops(p Problem) []City {
	m := make(map[City]bool)
	for _, f := range p.flights {
		m[f.From] = true
		m[f.To] = true
	}
	stops := make([]City, 0)
	for c, _ := range m {
		stops = append(stops, c)
	}
	return stops
}

func equal(a, b Flight) bool {
	return a.From == b.From && a.To == b.To && a.Day == b.Day && a.Cost == b.Cost
}

func (p Problem) route2solution(route []City) Solution {
	flights := make([]Flight, 0, p.n)
	var day Day = 0
	for i, current := range route {
		next := route[(i+1)%p.n]
		found_fi := -1
		for fi, flight := range p.flights {
			if flight.Day == day && flight.From == current && flight.To == next {
				found_fi = fi
				break
			}
		}
		if found_fi == -1 {
			panic(fmt.Sprintf("route2solution: flight not possible %v %v %v", day, current, next))
		}
		flights = append(flights, p.flights[found_fi])
		day++
	}
	return NewSolution(flights)
}

func contains(list []City, city City) bool {
	for _, c := range list {
		if c == city {
			return true
		}
	}
	return false
}

func printInfo(args ...interface{}) {
	if BeVerbose {
		args = append(args, "@", time.Since(StartTime))
		fmt.Fprintln(os.Stderr, args...)
	}
}

func expired(timeout <-chan time.Time) bool {
    select {
    case <-timeout: return true
    default: return false
    }
}

