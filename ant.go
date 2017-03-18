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
	fis []FlightIndex
}

var ants [20]ant
var ANTS = 20

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
		ants[ai].fis = make([]FlightIndex, 0, n)
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
					//printInfo("ant to die", ai, ants[ai].visited, "day", d, "city", ants[ai].city)
					die(ai) // TODO
					break
				}
				//printInfo("FI:", fi)
				feromones[fi] += 1.0
				flight := problem.flights[fi]
				ants[ai].total += flight.Cost
				ants[ai].city = flight.To
				if ants[ai].city == 0 {
					ants[ai].visited = ants[ai].visited[:0]
					ants[ai].fis = ants[ai].fis[:0]
					antsFinished++
				} else {
					ants[ai].visited = append(ants[ai].visited, ants[ai].city)
					ants[ai].fis = append(ants[ai].fis, fi)
				}
			}
			// TODO if ant dolezl do 0 then reset visited
			if antsFinished > ANTS {
				antsFinished = 0
				mf := float32(0.0)
				for fi := range feromones {
					feromones[fi] *= 0.85
					if feromones[fi] > mf { mf = feromones[fi] }
				}
				printInfo("Max feromone:", mf)
				//printInfo("Feromones:", feromones)
			}
		}
		// TODO vyparovani
	}
}

func die(ai int) {
	//printInfo("ant", ai, "dying")
	ants[ai].city = 0
	ants[ai].visited = ants[ai].visited[:0]
	for _, fi := range ants[ai].fis {
		feromones[fi] -= 1.0
	}
	ants[ai].fis = ants[ai].fis[:0]
	// keep current total cost for now; maybe add maximum flight cost or assign current worst running ant total
}

// ants don't fly

func antFlight(problem Problem, graph Graph, visited []City, day Day, city City) (FlightIndex, bool) {
	possible_flights := make([]FlightIndex, 0, MAX_CITIES)
	ft := make([]float32, 0, MAX_CITIES+1) // array of thresholds
	ft = append(ft, 0.0)
	var fsum float32 = 0.0
	for _, fi := range graph.antsGraph[city][day] {
/*
if int(day) == graph.size -1 && problem.flights[fi].From == city && problem.flights[fi].To == 0 {
	printInfo("possible last flight:", problem.flights[fi])
}
*/
		if contains(visited, problem.flights[fi].To) {
			continue
		}
		possible_flights = append(possible_flights, fi)
		fsum += float32(graph.size) + (10.0*float32(graph.size)/float32(ANTS))*feromones[fi]
		ft = append(ft, fsum)
	}
	flightCnt := len(possible_flights)
/*
	printInfo("---------------")
	printInfo("D:", day)
	printInfo("C:", city)
	printInfo("P:", possible_flights)
	//printInfo("F:", feromones)
	printInfo("T:", ft)
*/

	if flightCnt == 0 {
/*
printInfo("no route", "day", day, "city", city)
kam := make([]City, 0, problem.n)
for i:=0;i<problem.n;i++ {
	found := false
	for _, ii := range visited {
		if ii == City(i) { found = true }
	}
	if ! found { kam = append(kam, City(i)) }
}
printInfo("I already visited", len(visited), visited)
printInfo("I might fly to", len(kam), kam)
printInfo("xxx", graph.antsGraph[city][day])
for _, fi := range graph.antsGraph[city][day] {
	printInfo("return flight (ants)", problem.flights[fi])
}
printInfo("xxx", graph.fromDaySortedCost[city][day])
for _, f:= range graph.fromDaySortedCost[city][day] {
	printInfo("return flight (fdsc)", *f)
}
*/
		return 0, false
	}
	r := rand.Float32() * fsum
	//printInfo("R:", r)
	result := flightCnt - 1
	for i, f := range ft {
		if r < f {
			result = i-1
			break
		}
	}
	//printInfo("Res:", result)
//printInfo("possible flights", len(possible_flights), result, ft)
	return possible_flights[result], true
}
