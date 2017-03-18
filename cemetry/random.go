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

// Freaky engine finding pseudo-random paths
type RandomEngine struct {
	graph Graph
	seed  int
}

var RandomEngineResultsCounter uint32
var randomCurrentBest = Money(math.MaxInt32)

func (e RandomEngine) Name() string {
	return fmt.Sprintf("%s(%d)", "RndEngine", e.seed)
}

func (e RandomEngine) Solve(comm comm, p Problem) {
	//defer profile.Start(/*profile.MemProfile*/).Stop()
	rand.Seed(int64(e.seed) + time.Now().UTC().UnixNano())
	randomSolver(e.graph, comm, p.stats)
	//comm.done()
}

func randomSolver(graph Graph, comm comm, stats FlightStatistics) {
	solution := make([]Flight, 0, graph.size)
	var price Money
	var city City
	var toGo Day
	for {
		solution = solution[:0]
		visited := make([]City, 0, MAX_CITIES)
		city = City(0)
		price = Money(0)
		toGo = Day(graph.size)
		for d := 0; d < graph.size; d++ {
			//solution, city, price = randomFly(graph, solution, visited, d, city, price)
			flight, r := randomFlight(graph, visited, Day(d), toGo, city, stats)
			if !r {
				break
			}
			price += flight.Cost
			if price >= randomCurrentBest {
				break
			}
			city = flight.To
			visited = append(visited, city)
			solution = append(solution, flight)
			toGo--
		}
		if len(solution) == graph.size /*&& price < randomCurrentBest*/ {
			randomCurrentBest = price
			comm.sendSolution(NewSolution(solution))
		}
		RandomEngineResultsCounter++
	}
}

/*
func randomFly(graph Graph, solution []Flight, visited []City, day Day, city City, price Cost) []Flight, City, Cost {
	flightCnt = len(graph.data[city][day])
	if flightCnt == 0 {
		return nil
	}
	flight := graph.data[city][day][rand.Intn(flightCnt)]
	return append(solution, flight), flight.To, price + flight.Cost
}*/
func randomFlight(graph Graph, visited []City, day, toGo Day, city City, stats FlightStatistics) (Flight, bool) {
	possible_flights := make([]Flight, 0, MAX_CITIES)
	//progress := 1.0 - (float32(toGo)/float32(graph.size))
	for _, f := range graph.data[city][day] {
		if contains(visited, f.To) {
			continue
		}
		s := stats.ByDest[city][f.To]
		discount := s.AvgPrice - float32(f.Cost)
		discount_rate := discount / float32(f.Cost)
		//if f.Cost > 0 && discount_rate < -0.6 + float32(math.Abs(float64(0.5-progress))) {
		if f.Cost > 100 && discount_rate < -0.1 {
			// no discount, no deal, bro
			continue
		}
		possible_flights = append(possible_flights, *f)
	}
	flightCnt := len(possible_flights)

	if flightCnt == 0 {
		return Flight{0, 0, 0, 0, 0}, false
	}
	flight := possible_flights[rand.Intn(flightCnt)]
	return flight, true
}
