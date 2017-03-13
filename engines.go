package fsp

import (
	"math"
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

type bufferComm struct {
	buffer      *Solution
	bufferFree  <-chan bool
	bufferReady chan<- int
	queryBest   chan<- int
	receiveBest <-chan Money
	searchedAll chan<- int
	id          int
}

func (c *bufferComm) sendSolution(r Solution) Money {
	c.queryBest <- c.id
	bestCost := <-c.receiveBest
	if bestCost < r.totalCost {
		return bestCost
	}

	<-c.bufferFree
	for i := 0; i < len(r.flights); i++ {
		c.buffer.flights[i] = r.flights[i]
	}

	c.buffer.totalCost = r.totalCost
	c.bufferReady <- c.id
	//printInfo("New solution found with price", r.totalCost, "by", c.id, engines[c.id].Name() )
	return r.totalCost
}

func (c bufferComm) done() {
	c.searchedAll <- c.id
}

func initBuffer(size, engines int) []Solution {
	b := make([]Solution, engines)
	for i, _ := range b {
		b[i] = Solution{make([]Flight, size), 0}
	}
	return b
}

func initBufferChannels(engines int) []chan bool {
	bufferFree := make([]chan bool, engines)
	for i := 0; i < engines; i++ {
		bufferFree[i] = make(chan bool, 1)
	}
	return bufferFree
}

func initBestChannels(engines int) []chan Money {
	ch := make([]chan Money, engines)
	for i := 0; i < engines; i++ {
		ch[i] = make(chan Money, 1)
	}
	return ch
}

func initEngines(p Problem) []Engine {
	graph = NewGraph(p)
	printInfo("Graph ready")
	return []Engine{
		NewBottleneck(graph),
		Dcfs{graph, 0}, // single instance runs from start
		Dcfs{graph, 1}, // additional instances can start with n-th branch in 1st level
		//Dcfs{graph, 2},
		//Dcfs{graph, 3},
		//Mitm{},
		//Bhdfs{graph, 0},
		//NewGreedy(graph),
		RandomEngine{graph, 0},
	}
}

func noBullshit(b Solution, engine string) bool {
	visited := make(map[City]bool)
	prevFlight := b.flights[0]
	for _, flight := range b.flights[1:] {
		var flightFound bool
		for _, graphFlight := range graph.data[flight.From][flight.Day] {
			if *graphFlight == flight {
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

func saveBest(b *Solution, r Solution, engine string) {
	if b.totalCost > r.totalCost && noBullshit(r, engine) {
		for i, f := range r.flights {
			b.flights[i] = f
		}
		b.totalCost = r.totalCost
		printInfo("New best solution found by", engine, "with price", b.totalCost)
	}
}

func kickTheEngines(problem Problem, timeout <-chan time.Time) (Solution, error) {
	nCities := problem.n
	engines = initEngines(problem)

	//query/response what is current best
	bestResponse := initBestChannels(len(engines))
	bestQuery := make(chan int)

	//signalize goroutine they can write to their buffer
	bufferFree := initBufferChannels(len(engines))
	buffer := initBuffer(nCities, len(engines))
	best := Solution{make([]Flight, nCities), math.MaxInt32}

	//goroutine with id signals its buffer is ready
	bufferReady := make(chan int, len(engines))

	//goroutine signals it has searched the entire state space, we can finish
	done := make(chan int)

	for i, e := range engines {
		go e.Solve(&bufferComm{&buffer[i], bufferFree[i], bufferReady,
			bestQuery, bestResponse[i], done, i}, problem)
		bufferFree[i] <- true
	}
	for {
		select {
		case i := <-bufferReady:
			saveBest(&best, buffer[i], engines[i].Name())
			bufferFree[i] <- true
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
