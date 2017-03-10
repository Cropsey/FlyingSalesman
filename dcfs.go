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
	dcfs_solver(e.graph, p.stats, comm)
	//comm.done()
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
		//current_deal = float32(f.Cost) - s.AvgPrice * discount // - NO NO NO
		//current_deal = float32(f.Cost) * s.AvgPrice - s.AvgPrice * discount // (200, 300) = 39639, 51790
		//current_deal = float32(f.Cost) - 0.3 * discount // (200, 300) = 40722, 51625
		current_deal = float32(f.Cost) - 0.6*discount // (200, 300) = 40543, 48493
		//current_deal = float32(f.Cost) - 0.9 * discount // (200, 300) = 40447, 50580, total: 189785
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
