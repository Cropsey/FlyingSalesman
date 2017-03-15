package fsp

import (
	"fmt"
	"math"
	//"os"
	"sort"
	//"github.com/pkg/profile"
	"sync"
)

// Reverse node heuristics and DFS
type Bhdfs struct {
	graph Graph
	skip  int
}

var BhdfsResultsCounter uint32
var bhdfsCurrentBest = Money(math.MaxInt32)

func (e Bhdfs) Name() string {
	return fmt.Sprintf("%s(%d)", "Bhdfs", e.skip)
}

func (e Bhdfs) Solve(comm comm, p Problem) {
	//defer profile.Start(/*profile.MemProfile*/).Stop()
	bhdfsSolver(e.graph, p.stats, comm, e.skip)
	//comm.done()
}

func bhdfsSolver(graph Graph, stats FlightStatistics, comm comm, skip int) /*[]Flight*/ {

	printInfo("starting bhdfs solver", skip)
	visited := make([]City, 0, MAX_CITIES)
	solution := make([]Flight, 0, graph.size)
	home := City(0)
	day := Day(0)
	price := Money(0)
	var once sync.Once
	once.Do(func() {
		bhdfsEvaluate(graph)
		printInfo("bhdfs evaluation completed")
	})
	bhdfsIterate(solution, day, home, visited, graph, stats, price, comm, skip)
}

func bhdfsEvaluate(g Graph) {
	// evaluate each node in graph with best/worst price to reach final destination
	for day := g.size - 1; day >= 0; day-- {
		for i, flights := range g.dayFromData[day] {
			if day == g.size-1 {
				for j, f /*lights2*/ := range flights {
					//for _, f := range flights2 {
					g.dayFromData[day][i][j].Heuristic = f.Cost
					//}
				}
				continue
			}
			//for _, flights2 := range flights{
			for j, f := range flights { //flights on day day, from j
				best := Money(math.MaxInt32)
				worst := Money(0)
				for _, f2 := range g.dayFromData[day+1][f.To] {
					if f2.To == f.From {
						// avoid short cycles (how to avoid long ones?)
						// printInfo("short cycle--", f, f2)
						continue
					}

					if f2.Cost < best {
						best = f2.Heuristic
					}
					if f2.Cost > worst {
						worst = f2.Heuristic
					}

					//printInfo("candidate on", day, f, f2)
				}
				//f.Heuristic = best
				g.dayFromData[day][i][j].Heuristic = (worst+best)/2 + f.Cost
				//printInfo("day", day, "from", i, "worst", worst, f)
			}
			//}
			//printInfo(flights)
		}
		//printInfo("Day:", day, )
	}
	/*
		for i, x := range g.dayFromData {
			for j, y := range x {
				for _, z := range y {
					printInfo("[", i, j, "->", z.To, "]:", *z)
					}
				}
			}
	*/
}

func bhdfsInsertSortedFlight(slice []EvaluatedFlight, node EvaluatedFlight) []EvaluatedFlight {
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
	//tail := append([]EvaluatedFlight{node}, slice[i:]...)
	return append(slice[0:i], append([]EvaluatedFlight{node}, slice[i:l]...)...)
}

func bhdfsIterate(partial []Flight, day Day, current City,
	visited []City, graph Graph, stats FlightStatistics, price Money, comm comm, skip int) {

	if price >= bhdfsCurrentBest {
		// we have already got worse than best result, give it up, bro
		BhdfsResultsCounter++
		return
	}
	if int(day) == graph.size {
		BhdfsResultsCounter++
		//if price < bhdfsCurrentBest {
		//bhdfsCurrentBest = price
		bhdfsCurrentBest = comm.sendSolution(NewSolution(partial))
		//}
		return
	}
	//fmt.Fprintln(os.Stderr, "I am at", current, "day is", day)
	var current_deal float32
	//var current_deal int32
	possible_flights := make([]EvaluatedFlight, 0, MAX_CITIES)
	for _, f := range graph.fromDaySortedCost[current][day] {
		if contains(visited, f.To) {
			continue
		}
		s := stats.ByDest[current][f.To]
		discount := s.AvgPrice - float32(f.Cost)
		discount_rate := discount / float32(f.Cost)
		//if discount_rate < -0.3 {
		if f.Cost > 650 && discount_rate < -0.3 {
			// no discount, no deal, bro
			continue
		}
		current_deal = float32(f.Cost+f.Heuristic/2) - 0.6*discount
		//printInfo(f)

		//possible_flights = append(possible_flights, EvaluatedFlight{*f, current_deal})
		possible_flights = bhdfsInsertSortedFlight(possible_flights, EvaluatedFlight{*f, current_deal})
	}
	//sort.Sort(byValue(possible_flights))
	for i, f := range possible_flights {
		if day == 0 && skip > i {
			skip--
			continue
		}
		bhdfsIterate(append(partial, f.flight),
			day+1,
			f.flight.To,
			append(visited, f.flight.To),
			//bhdfsInsertVisited(visited, f.flight.To),
			graph, stats,
			price+f.flight.Cost,
			comm, skip)
	}
	return //[]Flight{}
}
