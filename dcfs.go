package fsp

import (
	//"fmt"
	"math"
	//"os"
	"sort"
)

// Depth + Cheapest First Search engine
// a variant of greedy DFS using cheapest next flight first with heuristics based on average price for same flights on different days
type Dcfs struct {
	graph Graph
}

func (e Dcfs) Name() string {
	return "Dcfs"
}

var dcfsCurrentBest = Money(math.MaxInt32)

func (e Dcfs) Solve(comm comm, p Problem) {
	/*stops := stops(p)
	if len(stops) < 2 {
		comm.sendSolution(Solution{})
		return
	}
	visited := make([]City, 1, len(stops))
	visited[0] = 0
	to_visit := append(stops[1:], stops[0])*/
	//solution := make([]Flight, 0, e.graph.size)
	//comm.sendSolution(NewSolution(dcfs_solver(solution, City(0), e.graph, p.stats)))
	//comm.sendSolution(NewSolution(dcfs(solution, visited, to_visit, flights)))
	dcfs_solver(e.graph, p.stats, comm)
	comm.done()
}

type EvaluatedFlight struct {
	flight Flight
	value  float32
}

type byValue []EvaluatedFlight

func (f byValue) Len() int {
	return len(f)
}

func (f byValue) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f byValue) Less(i, j int) bool {
	return f[i].value < f[j].value
}

func insertSorted(slice []EvaluatedFlight, node EvaluatedFlight) []EvaluatedFlight {
	l := len(slice)
	if l == 0 {
		return []EvaluatedFlight{node}
	}
	i := sort.Search(l, func(i int) bool { return slice[i].value > node.value })
	//fmt.Println(i)
	if i == 0 {
		return append([]EvaluatedFlight{node}, slice...)
	}
	if i == -1 {
		return append(slice[0:l], node)
	}
	tail := append([]EvaluatedFlight{node}, slice[i:]...)
	return append(slice[0:i], tail...)
}

func dcfs_solver(graph Graph, stats [][]FlightStats, comm comm) /*[]Flight*/ {

	printInfo("starting dcfs solver")
	visited := make([]City, 0, MAX_CITIES)
	solution := make([]Flight, 0, graph.size)
	home := City(0)
	day := Day(0)
	price := Money(0)
	/*
		for i := 0; i < graph.size; i++ {
			fmt.Fprintln(os.Stderr, "I am at", current, "and can go to:")
			//best_deal := float32(-math.MaxUint16)
			best_deal := float32(math.MaxUint16)
			var best_flight Flight
			var current_deal float32
			possible_flights := make([]EvaluatedFlight, 0, MAX_CITIES)
			for _, f := range graph.data[current][day] {
				if contains(visited, f.To) {
					continue
				}
				s := stats[current][f.To]
				discount := s.AvgPrice - float32(f.Cost)
				//if discount > best_deal {
				//if float32(f.Cost) < best_deal {
				current_deal = float32(f.Cost) - 1.0 * discount
				possible_flights = append(possible_flights, EvaluatedFlight{f, current_deal}) //TODO: create some kind of sorted insert
				if current_deal < best_deal {
					//best_deal = discount
					//best_deal = float32(f.Cost)
					best_deal = current_deal
					best_flight = f
					//FIXME: if there is no flight possible, it will fail here
				}
				fmt.Fprintln(os.Stderr, " - ", f.To,
					" for ", f.Cost, "avg is", s.AvgPrice,
					" it saves ", discount, "money")
			}
			sort.Sort(byValue(possible_flights))
			if best_deal < 0 {
				fmt.Fprintln(os.Stderr, "OMG, it's a trap all flights are overpriced here!!!")
			}
			fmt.Fprintln(os.Stderr, "going to", best_flight.To)
			fmt.Fprintln(os.Stderr)
			day += 1
			current = best_flight.To
			visited = append(visited, best_flight.To)
			solution = append(solution, best_flight)

		} */
	dcfs_iterate(solution, day, home, visited, graph, stats, price, comm)
}

func dcfs_iterate(partial []Flight, day Day, current City,
	visited []City, graph Graph, stats [][]FlightStats, price Money, comm comm) {

	if price > dcfsCurrentBest {
		// we have already got worse than best result, give it up, bro
		return
	}
	if int(day) == graph.size {
		if price < dcfsCurrentBest {
			dcfsCurrentBest = price
			comm.sendSolution(NewSolution(partial))
		}
		return
	}
	//fmt.Fprintln(os.Stderr, "I am at", current, "day is", day)
	var current_deal float32
	possible_flights := make([]EvaluatedFlight, 0, MAX_CITIES)
	for _, f := range graph.data[current][day] {
		if contains(visited, f.To) {
			continue
		}
		s := stats[current][f.To]
		discount := s.AvgPrice - float32(f.Cost)
		discount_rate := discount / float32(f.Cost)
		if discount_rate < -0.5 {
			// no discount, no deal, bro
			continue
		}
		current_deal = float32(f.Cost) - 0.25*discount
		//possible_flights = append(possible_flights, EvaluatedFlight{f, current_deal})
		possible_flights = insertSorted(possible_flights, EvaluatedFlight{f, current_deal})
	}
	//sort.Sort(byValue(possible_flights))
	for _, f := range possible_flights {
		dcfs_iterate(append(partial, f.flight),
			day+1,
			f.flight.To,
			append(visited, f.flight.To),
			graph, stats,
			price+f.flight.Cost,
			comm)
	}
	return //[]Flight{}
}
