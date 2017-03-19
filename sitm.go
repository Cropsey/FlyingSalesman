package fsp

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	//"github.com/pkg/profile"
)

// Reverse node heuristics and DFS
type Sitm struct {
	graph Graph
	skip  int
}

var SitmResultsCounter uint32
var SitmBranchCounter []uint32
var sitmCurrentBest = Money(math.MaxInt32)

// engine parms
var sitmMaxBranches = 0 // default value set in Solve() to graph.size/2
var sitmDiscountWeight = float32(0.6)
var sitmMinDiscount = float32(-0.3)
var sitmDiscountThreshold = Money(650)

func sitmLoadEnvParams() {
	var env string
	env = os.Getenv("SITM_MAX_BRANCHES")
	if len(env) > 0 {
		sitmMaxBranches, _ = strconv.Atoi(env)
	}
	env = os.Getenv("SITM_DISC_W")
	if len(env) > 0 {
		tmp, _ := strconv.ParseFloat(env, 32)
		sitmDiscountWeight = float32(tmp)
	}
	env = os.Getenv("SITM_MIN_DISC")
	if len(env) > 0 {
		tmp, _ := strconv.ParseFloat(env, 32)
		sitmMinDiscount = float32(tmp)
	}
	env = os.Getenv("SITM_DISC_THRESH")
	if len(env) > 0 {
		tmp, _ := strconv.Atoi(env)
		sitmDiscountThreshold = Money(tmp)
	}
}

func (e Sitm) Name() string {
	return fmt.Sprintf("%s(%d)", "Sitm", e.skip)
}

func (e Sitm) Solve(comm comm, p Problem) {
	//defer profile.Start(/*profile.MemProfile*/).Stop()
	sitmLoadEnvParams()
	if sitmMaxBranches == 0 {
		sitmMaxBranches = e.graph.size / 2
		if p.n >= 50 {
			sitmMaxBranches = 2
		}
	}
	SitmBranchCounter = make([]uint32, e.graph.size+1)
	sitmSolver(e.graph, p.stats, comm, e.skip)
	//comm.done()
}

type evaluatedCity struct {
	city City
	//value Money
	value float32
}

type evaluatedCityByValue []evaluatedCity

func (f evaluatedCityByValue) Len() int {
	return len(f)
}

func (f evaluatedCityByValue) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f evaluatedCityByValue) Less(i, j int) bool {
	return f[i].value < f[j].value
}

func sitmSolver(graph Graph, stats FlightStatistics, comm comm, skip int) /*[]Flight*/ {

	printInfo("starting sitm solver", skip)
	visited := make([]City, 0, MAX_CITIES)
	solution := make([]Flight, 0, graph.size)
	//home := City(0)
	day := Day(graph.size / 2)
	// find cheapest city in the middle
	//cheapestTotal := Money(math.MaxInt32)
	bestDisc := float32(-math.MaxFloat32)
	//cheapestToCity := City(0)
	//cheapestC := City(0)
	evaluatedCities := make([]evaluatedCity, 0, graph.size)
	for i := 1; i < graph.size; i++ {
		// forward
		//cheapestF := Money(math.MaxInt32)
		bestDiscF := float32(-math.MaxFloat32)
		for _, f := range graph.data[i][day] {
			s := stats.ByDest[f.From][f.To]
			discount := s.AvgPrice - float32(f.Cost)
			/*if cheapestF > f.Cost {
				cheapestF = f.Cost
				cheapestToCity = f.To
			}*/
			if bestDiscF < discount {
				bestDiscF = discount
				//cheapestToCity = f.To
			}
		}
		// backward
		//cheapestB := Money(math.MaxInt32)
		bestDiscB := float32(-math.MaxFloat32)
		for _, f := range graph.toDayData[i][day-1] {
			s := stats.ByDest[f.From][f.To]
			discount := s.AvgPrice - float32(f.Cost)
			// ignore flight back in price
			/*if cheapestToCity != f.From && cheapestB > f.Cost {
				cheapestB = f.Cost
			}*/
			if bestDiscB < discount {
				bestDiscB = discount
				//cheapestToCity = f.To
			}
		}
		/*
			if cheapestTotal > (cheapestF + cheapestB) {
				cheapestTotal = cheapestF + cheapestB
				//cheapestC = City(i)
			}*/
		if bestDisc < (bestDiscF + bestDiscB) {
			bestDisc = bestDiscF + bestDiscB
		}
		//evaluatedCities = append(evaluatedCities, evaluatedCity{City(i), cheapestTotal})
		evaluatedCities = append(evaluatedCities, evaluatedCity{City(i), bestDisc})
	}
	sort.Sort(sort.Reverse(evaluatedCityByValue(evaluatedCities)))
	//printInfo("Cheapest city in the middle:", cheapestC, cheapestTotal)
	for i, city := range evaluatedCities {
		if skip > 0 {
			skip--
			continue
		}
		printInfo("City in the middle ", city, i)
		price := Money(0)
		sitmIterate(true, solution, day, day-1, city.city, city.city,
			append(visited, city.city), graph, stats, price, comm, skip)
	}
}

func sitmInsertSortedFlight(slice []EvaluatedFlight, node EvaluatedFlight) []EvaluatedFlight {
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

func sitmIterate(forward bool, partial []Flight, dayF, dayB Day, cityF, cityB City,
	visited []City, graph Graph, stats FlightStatistics, price Money, comm comm, skip int) {

	if price >= sitmCurrentBest {
		// we have already got worse than best result, give it up, bro
		SitmResultsCounter++
		return
	}
	if len(partial) == graph.size {
		SitmResultsCounter++
		sitmCurrentBest = comm.sendSolution(NewSolution(partial))
		return
	}
	var currentDeal float32
	possibleFlights := make([]EvaluatedFlight, 0, MAX_CITIES)
	if forward {
		//printInfo("forward day", dayF, "at", cityF)
		SitmBranchCounter[dayF] += 1
		for _, f := range graph.fromDaySortedCost[cityF][dayF] {
			if contains(visited, f.To) {
				continue
			}
			/*
				s := stats.ByDest[cityF][f.To]
				discount := s.AvgPrice - float32(f.Cost)
				//discount_rate := discount / float32(f.Cost)*/
			currentDeal = float32(f.Cost) //- 0.6*discount
			possibleFlights = sitmInsertSortedFlight(possibleFlights, EvaluatedFlight{*f, currentDeal})
		}
		dayF++
	} else { // backward
		//printInfo("backward day", dayB, "at", cityB)
		SitmBranchCounter[dayB] += 1
		for _, f := range graph.toDayData[cityB][dayB] {
			if contains(visited, f.From) {
				continue
			}
			currentDeal = float32(f.Cost)
			possibleFlights = sitmInsertSortedFlight(possibleFlights, EvaluatedFlight{*f, currentDeal})
		}
		dayB--
	}
	//printInfo(possibleFlights)
	if len(possibleFlights) > sitmMaxBranches {
		possibleFlights = possibleFlights[:sitmMaxBranches]
	}

	for _, f := range possibleFlights {
		if forward {
			if f.flight.To != City(0) {
				visited = append(visited, f.flight.To)
			}
			sitmIterate(
				!forward, // cycle forward and backward
				append(partial, f.flight),
				dayF, dayB,
				f.flight.To,
				cityB,
				visited,
				graph, stats,
				price+f.flight.Cost,
				comm, skip)
		} else { // backward
			if f.flight.From != City(0) {
				visited = append(visited, f.flight.From)
			}
			sitmIterate(
				!forward, // cycle forward and backward
				append(partial, f.flight),
				dayF, dayB,
				cityF,
				f.flight.From,
				visited,
				graph, stats,
				price+f.flight.Cost,
				comm, skip)
		}

	}
	//printInfo("Sitm: no more possible flights, yay!", SitmResultsCounter)

	return
}
