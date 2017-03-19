package fsp

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
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
	day     Day
	city    City
	total   Money
	visited []City
	fis     []FlightIndex
}

const EVAPORATE_P = 0.2 // percent to evaporate
const FEROM_C = 2.0
const PRICE_C = 1.0
const FEROMONE_WEIGHT = 0.9

var ANTS = 0
var ants []ant

var antCurrentBest = Money(math.MaxInt32)

func (e AntEngine) Name() string {
	return fmt.Sprintf("%s(%d)", "AntEngine", e.seed)
}

func (e AntEngine) Solve(comm comm, p Problem) {
	//defer profile.Start(/*profile.MemProfile*/).Stop()
	fmt.Fprintf(os.Stderr, "") // TODO anti error, remove
	rand.Seed(int64(e.seed) + time.Now().UTC().UnixNano())
	feromones = make([]float32, len(p.flights))
	antInit(p.n/2, p.n)
	antSolver(p, e.graph, comm)
	//comm.done()
}

func antInit(ant_n, problem_n int) {
	ANTS = ant_n
	ants = make([]ant, ant_n, ant_n)
	for ai := range ants {
		ants[ai].visited = make([]City, 0, problem_n)
		ants[ai].fis = make([]FlightIndex, 0, problem_n)
	}
}

func antSolver(problem Problem, graph Graph, comm comm) {
	//solution := make([]Flight, 0, graph.size)
	antsFinished := 0
	for {
		minTotal := ants[0].total
		minIndex := 0
		for ai := range ants {
			if ants[ai].total < minTotal {
				minTotal = ants[ai].total
				minIndex = ai
			}
		}
		ai := minIndex // the chosen one
		//printInfo("The chosen one", ai, ants[ai])
		//printInfo("Ant:", ai)
		fi, r := antFlight(problem, graph, ants[ai].visited, ants[ai].day, ants[ai].city)
		if !r {
			//printInfo("ant to die", ai, ants[ai].visited, "day", ants[ai].day, "city", ants[ai].city)
			die(ai)
			continue
		}
		//printInfo("FI:", fi)
		flight := problem.flights[fi]
		ants[ai].total += flight.Cost
		ants[ai].day++
		ants[ai].city = flight.To
		if ants[ai].city == 0 { // ant has completed the route
			ants[ai].day = 0
			evaporate(EVAPORATE_P)
			// place the feromones
			for _, fi := range ants[ai].fis {
				feromones[fi] += 1.0
			}
			ants[ai].visited = ants[ai].visited[:0]
			ants[ai].fis = ants[ai].fis[:0]
			antsFinished++
		} else {
			ants[ai].visited = append(ants[ai].visited, ants[ai].city)
			ants[ai].fis = append(ants[ai].fis, fi)
		}
		if antsFinished > ANTS {
			//printInfo("ants finished")
			antsFinished = 0
			//printInfo("Feromones:", feromones)
			followAnts(problem, graph, comm)
		}
	}
}

func evaporate(x float32) {
	mf := float32(0.0)
	remain := 1.0 - x
	for fi := range feromones {
		feromones[fi] *= remain
		if feromones[fi] > mf {
			mf = feromones[fi]
		}
	}
	//printInfo("Max feromone:", mf)
}

func followAnts(problem Problem, graph Graph, comm comm) {
	solution := make([]Flight, 0, graph.size)
	var price Money
	var city City
	for {
		solution = solution[:0]
		visited := make([]City, 0, MAX_CITIES)
		city = City(0)
		price = Money(0)
		for d := 0; d < graph.size; d++ {
			//printInfo("FA:")
			fi, r := antFlight(problem, graph, visited, Day(d), city)
			if !r {
				return
			}
			price += problem.flights[fi].Cost
			if price >= randomCurrentBest {
				break
			}
			city = problem.flights[fi].To
			visited = append(visited, city)
			solution = append(solution, problem.flights[fi])
		}
		if len(solution) == graph.size && price < antCurrentBest {
			antCurrentBest = price
			comm.sendSolution(NewSolution(solution))
			//printInfo("ant solution sent, price", price)
			/*
				printInfo("Stats:")
				dg := make([]struct{maxF float32; flights,f50 int}, problem.n)
				for _, dtfi := range graph.antsGraph {
					for d, tfi := range dtfi {
						for _, fi := range tfi {
							if dg[int(d)].maxF < feromones[fi] {
								dg[int(d)].maxF = feromones[fi]
							}
							dg[int(d)].flights += 1
						}
					}
				}
				for _, dtfi := range graph.antsGraph {
					for d, tfi := range dtfi {
						for _, fi := range tfi {
							if dg[int(d)].maxF/2.0 < feromones[fi] {
								dg[int(d)].f50 += 1
							}
						}
					}
				}
				for d:=0; d<problem.n; d++ {
					x := dg[int(d)]
					printInfo("day", d, "max", x.maxF, "flights", x.flights, "flights>50%", x.f50)
				}
			*/
			return
		}
	}
}

func die(ai int) {
	//printInfo("ant", ai, "dying")
	ants[ai].day = 0
	ants[ai].city = 0
	ants[ai].visited = ants[ai].visited[:0]
	ants[ai].fis = ants[ai].fis[:0]
	// keep current total cost for now; maybe add maximum flight cost or assign current worst running ant total
}

// ants don't fly

func antWeight(problem Problem, fi FlightIndex, flights int, avgCost float64, avgFeromones float32) float32 {
	// feromones influence
	price := problem.flights[fi].Cost
	rel_price := avgCost / float64(price) // 1.0 for average, 2.0 for 2x better than average
	//printInfo("xxx", avgFeromones, flights, feromones[fi], ANTS)
	rel_feromones := 1.0
	if avgFeromones > 0.0 {
		rel_feromones = float64(avgFeromones) * (1.0 - FEROMONE_WEIGHT)
		rel_feromones += float64(feromones[fi]/avgFeromones) * FEROMONE_WEIGHT
	}
	//fmt.Fprintf(os.Stderr, "rf avg %.2f cur %.2f res %.2f %v\n", avgFeromones, feromones[fi], rel_feromones, flights)
	f := math.Pow(rel_feromones, FEROM_C)
	// price influence
	p := math.Pow(rel_price, PRICE_C)
	var result float32 = float32(f * p)
	//fmt.Fprintf(os.Stderr, "f/p: %.4f * %.2f = %.4f, (feromones %.2f/%.2f, cost %v, fi %v)\n", f, p, result, feromones[fi], rel_feromones, price, fi)
	return result
}

// choose the flight ant will take
func antFlight(problem Problem, graph Graph, visited []City, day Day, city City) (FlightIndex, bool) {
	// first, find all possible flights and construct random distribution
	possible_flights := make([]FlightIndex, 0, MAX_CITIES)
	var maxCost Money = 0 // needed to normalize costs
	var sumCost Money = 0
	var sumFeromones float32 = 0.0
	//var mw float32 = 0.0
	for _, fi := range graph.antsGraph[city][day] {
		if contains(visited, problem.flights[fi].To) {
			continue
		}
		possible_flights = append(possible_flights, fi)
		cost := problem.flights[fi].Cost
		if cost > maxCost {
			maxCost = cost
		}
		sumCost += cost
		sumFeromones += feromones[fi]
	}
	//printInfo("mw:", mw)
	flightCnt := len(possible_flights)
	var avgCost = float64(sumCost) / float64(flightCnt)
	var avgFeromones = sumFeromones / float32(flightCnt)

	// second, return that ant is stuck if no flight possible
	if flightCnt == 0 {
		return 0, false
	}

	// third, compute weights, we do in in extra cycle beacause of normalization
	var fsum float32 = 0.0
	thres := make([]float32, 0, MAX_CITIES+1) // array of thresholds
	thres = append(thres, 0.0)                // easier logic later if we always start with 0.0
	for _, fi := range possible_flights {
		// compute weight of the flight
		// TODO scale according to average flight price
		w := antWeight(problem, fi, flightCnt, avgCost, avgFeromones)
		//if w > mw { mw = w }
		fsum += w
		thres = append(thres, fsum)
	}

	// fourth, choose flight randomly based on the distribution
	r := rand.Float32() * fsum
	result := flightCnt - 1
	for i, f := range thres {
		if r < f {
			result = i - 1
			break
		}
	}
	return possible_flights[result], true
}
