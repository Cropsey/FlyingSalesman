package fsp

import (
	"fmt"
	"os"
	"math"
)

// Depth + Cheapest First Search engine
// finds single naive solution by using the cheapest available flight from a city
type Dcfs struct{
	graph Graph
}

func (e Dcfs) Name() string {
	return "Dcfs"
}

func (e Dcfs) Solve(comm comm, p Problem) {
	stops := stops(p)
	if len(stops) < 2 {
		comm.sendSolution(Solution{})
		return
	}
	visited := make([]City, 1, len(stops))
	visited[0] = 0
	to_visit := append(stops[1:], stops[0])
	solution := make([]Flight, 0, len(stops))
	comm.sendSolution(NewSolution(dcfs_solver(solution, City(0), to_visit, e.graph, p.stats)))
	//comm.sendSolution(NewSolution(dcfs(solution, visited, to_visit, flights)))
}

/*
func dcfs(partial []Flight, visited, to_visit []City, flights []Flight) []Flight {
	if len(to_visit) == 0 {
		return partial
	}
	for _, f := range flights {
		if f.From == visited[len(visited)-1] {
			if si := indexOf(to_visit, f.To); si != -1 {
				solution := dcfs(append(partial, f),
					append(visited, f.To),
					append(to_visit[:si], to_visit[si+1:]...),
					flights)
				if len(solution) != 0 {
					// soluton found, yaaaay!
					return solution
				}
			}
		}
	}
	// no solution
	return []Flight{}
}
*/

func dcfs_solver(solution []Flight, home City, to_visit []City, graph Graph, stats [][]FlightStats) []Flight {

	fmt.Fprintln(os.Stderr,"dcfs solver")
	visited := make([]City, 0, MAX_CITIES)
	current := home
	day := 0
	for i := 0; i < graph.size; i++ {
		fmt.Fprintln(os.Stderr, "I am at", current, "and can go to:")
		best_deal := float32(-math.MaxUint16)
		var best_flight Flight
		for _, f := range graph.data[current][day] {
			if contains(visited, f.To) {
				continue
			}
			s := stats[current][f.To]
			discount := s.AvgPrice - float32(f.Cost)
			if discount > best_deal {
				best_deal = discount
				best_flight = f
				//FIXME: if there is no flight possible, it will fail here
			}
			fmt.Fprintln(os.Stderr, " - ", f.To,
			" for ", f.Cost, "avg is", s.AvgPrice,
			" it saves ", discount, "money")
		}
		fmt.Fprintln(os.Stderr, "going to", best_flight.To)
		fmt.Fprintln(os.Stderr)
		day += 1
		current = best_flight.To
		visited = append(visited, best_flight.To)
		solution = append(solution, best_flight)

	}
	return solution
}

func next_flight(to_visit []City, flights []Flight) Flight{

	return flights[0]
}
