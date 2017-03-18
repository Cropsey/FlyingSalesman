package fsp

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	//"os"
	//"sort"
	//"github.com/pkg/profile"
)

// Freaky engine finding pseudo-ant paths
type AntEngine struct {
	graph Graph
	seed  int
}

var feromones []float32
type ant struct {
	city City
	total Money
	visited []City
}

var ants [100]ant

var antCurrentBest = Money(math.MaxInt32)

func (e AntEngine) Name() string {
	return fmt.Sprintf("%s(%d)", "AntEngine", e.seed)
}

func (e AntEngine) Solve(comm comm, p Problem) {
	//defer profile.Start(/*profile.MemProfile*/).Stop()
	rand.Seed(int64(e.seed) + time.Now().UTC().UnixNano())
	feromones = make([]float32, len(p.flights))
	antInit(p.n)
	antSolver(p, e.graph, comm)
	//comm.done()
}

func antInit(n int) {
	for ai := range ants {
		ants[ai].visited = make([]City, 0, n)
	}
}

func antSolver(problem Problem, graph Graph, comm comm) {
	//solution := make([]Flight, 0, graph.size)
	antsFinished := 0
	for {
		for ai := range ants {
			// TODO shift ants according to their current totatl
			for d := 0; d < graph.size; d++ {
				fi, r := antFlight(problem, graph, ants[ai].visited, Day(d), ants[ai].city)
				if !r {
					die(ai) // TODO
					continue
				}
				printInfo("FI:", fi)
				feromones[fi] += 1.0
				flight := problem.flights[fi]
				ants[ai].total += flight.Cost
				ants[ai].city = flight.To
				if ants[ai].city == 0 {
					ants[ai].visited = ants[ai].visited[:0]
					antsFinished++
				}
				ants[ai].visited = append(ants[ai].visited, ants[ai].city)
				// TODO place feromones
			}
			// TODO if ant dolezl do 0 then reset visited
			if antsFinished > 10 {
				antsFinished = 0
				for fi := range feromones {
					feromones[fi] *= 0.5
				}
				printInfo("Feromones:", feromones)
			}
		}
		// TODO vyparovani
	}
}

func die(ai int) {
	ants[ai].city = 0
	ants[ai].visited = ants[ai].visited[:0]
	// keep current total cost for now; maybe add maximum flight cost or assign current worst running ant total
}

// ants don't fly

func antFlight(problem Problem, graph Graph, visited []City, day Day, city City) (FlightIndex, bool) {
	possible_flights := make([]FlightIndex, 0, MAX_CITIES)
	ft := make([]float32, 0, MAX_CITIES+1) // array of thresholds
	ft = append(ft, 0.0)
	var fsum float32 = 0.0
	for _, fi := range graph.antsGraph[city][day] {
		if contains(visited, problem.flights[fi].To) {
			continue
		}
		possible_flights = append(possible_flights, fi)
		fsum += float32(graph.size) + (float32(graph.size)/25.0)*feromones[fi]
		ft = append(ft, fsum)
	}
	flightCnt := len(possible_flights)
	printInfo("---------------")
	printInfo("P:", possible_flights)
	//printInfo("F:", feromones)
	printInfo("T:", ft)

	if flightCnt == 0 {
		return 0, false
	}
	r := rand.Float32() * fsum
	printInfo("R:", r)
	result := flightCnt - 1
	for i, f := range ft {
		if r < f {
			result = i-1
			break
		}
	}
	printInfo("Res:", result)
//printInfo("possible flights", len(possible_flights), result, ft)
	return possible_flights[result], true
}
