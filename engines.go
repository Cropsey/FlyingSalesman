package fsp

import (
	"math"
	"sort"
	"time"
)

var engines []Engine
var graph Graph

type Engine interface {
	Name() string
	Solve(comm comm, problem Problem)
}

type comm interface {
	sendSolution(r Solution) Money
	done()
}

type update struct {
	s Solution
	i int
}

type solutionComm struct {
	solutionReady chan<- update
	queryBest     chan<- int
	receiveBest   <-chan Money
	searchedAll   chan<- int
	id            int
}

func (c *solutionComm) sendSolution(r Solution) Money {
	c.queryBest <- c.id
	bestCost := <-c.receiveBest
	if bestCost < r.totalCost {
		return bestCost
	}

	solution := make([]Flight, len(r.flights))
	copy(solution, r.flights)
	sort.Sort(ByDay(solution))

	c.solutionReady <- update{NewSolution(solution), c.id}
	//printInfo("New solution found with price", r.totalCost, "by", c.id, engines[c.id].Name() )
	return r.totalCost
}

func (c solutionComm) done() {
	c.searchedAll <- c.id
}

func initBestChannels(engines int) []chan Money {
	ch := make([]chan Money, engines)
	for i := 0; i < engines; i++ {
		ch[i] = make(chan Money, 1)
	}
	return ch
}

func initEngines(p Problem) ([]Engine, Polisher) {
	graph = NewGraph(p)
	printInfo("Graph ready")
	polisher := NewPolisher(graph)
	return []Engine{
		NewBottleneck(graph),
		Dcfs{graph, 0}, // single instance runs from start
		Dcfs{graph, 1}, // additional instances can start with n-th branch in 1st level
		Dcfs{graph, 2},
		//Dcfs{graph, 3},
		Mitm{},
		Bhdfs{graph, 0},
		Bhdfs{graph, 1}, // we should avoid running evaluation phase of Bhdfs more than once
		//Bhdfs{graph, 2},
		//NewGreedy(graph),
		//RandomEngine{graph, 0},
		Sitm{graph, 0},
		NewGreedyRounds(graph),
		polisher,
	}, polisher
}

func sameFlight(f1, f2 Flight) bool {
	//ignore heuristics part in comparison as it can change during processing
	if f1.From == f2.From && f1.To == f2.To && f1.Day == f2.Day && f1.Cost == f2.Cost {
		return true
	}
	return false
}

func noBullshit(b Solution, engine string) bool {
	visited := make(map[City]bool)
	prevFlight := b.flights[0]
	for _, flight := range b.flights[1:] {
		var flightFound bool
		for _, graphFlight := range graph.data[flight.From][flight.Day] {
			//if *graphFlight == flight {
			if sameFlight(*graphFlight, flight) {
				flightFound = true
				break
			}
		}
		if !flightFound {
			printInfo("!!!", engine, "tried to bullshit sending fake flight", flight)
			return false
		}
		if visited[flight.To] {
			printInfo("!!!", engine, "tried to bullshit visiting city", flight.To, "twice")
			return false
		}
		if prevFlight.To != flight.From {
			printInfo("!!!", engine, "tried to bullshit with not connecting flights", prevFlight, flight)
			return false
		}
		visited[flight.To] = true
		prevFlight = flight
	}
	return true
}

func saveBest(b *Solution, r Solution, engine string) bool {
	if b.totalCost > r.totalCost && noBullshit(r, engine) {
		for i, f := range r.flights {
			b.flights[i] = f
		}
		b.totalCost = r.totalCost
		printInfo("New best solution found by", engine, "with price", b.totalCost)
		return true
	}
	return false
}

func kickTheEngines(problem Problem, timeout <-chan time.Time) (Solution, error) {
	nCities := problem.n
	engines, polisher := initEngines(problem)

	//query/response what is current best
	bestResponse := initBestChannels(len(engines))
	bestQuery := make(chan int)

	//signalize goroutine they can write to their buffer
	sol := make(chan update, len(engines))
	best := Solution{make([]Flight, nCities), math.MaxInt32}

	//goroutine signals it has searched the entire state space, we can finish
	done := make(chan int)

	for i, e := range engines {
		go e.Solve(&solutionComm{sol, bestQuery, bestResponse[i], done, i}, problem)
	}
	for {
		select {
		case u := <-sol:
			isBest := saveBest(&best, u.s, engines[u.i].Name())
			if isBest {
				polisher.try(u)
			}
		case i := <-bestQuery:
			bestResponse[i] <- best.totalCost
		case i := <-done:
			printInfo("Fearles engine", engines[i].Name(), "thinks it's done, let's see")
			return best, nil
		case <-timeout:
			printInfo("Out of time!")
			return best, nil
		}
	}
}
