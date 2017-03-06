package fsp

import (
	"math"
	"time"
)

var best Solution

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
	searchedAll chan<- bool
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
	return r.totalCost
}

func (c bufferComm) done() {
	c.searchedAll <- true
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
	graph := NewGraph(p)
	return []Engine{
		DFSEngine{graph, true},
		DFSEngine{graph, false},
		Mitm{},
	}
}

func saveBest(b *Solution, r Solution) {
	if b.totalCost > r.totalCost {
		for i, f := range r.flights {
			b.flights[i] = f
		}
		b.totalCost = r.totalCost
	}
}

func kickTheEngines(problem Problem, timeout <-chan time.Time) (Solution, error) {
	nCities := problem.n
	engines := initEngines(problem)

	//query/response what is current best
	bestResponse := initBestChannels(len(engines))
	bestQuery := make(chan int)

	//signalize goroutine they can write to their buffer
	bufferFree := initBufferChannels(len(engines))
	buffer := initBuffer(nCities, len(engines))
	best = Solution{make([]Flight, nCities), math.MaxInt32}

	//goroutine with id signals its buffer is ready
	bufferReady := make(chan int, len(engines))

	//goroutine signals it has searched the entire state space, we can finish
	done := make(chan bool)

	for i, e := range engines {
		go e.Solve(&bufferComm{&buffer[i], bufferFree[i], bufferReady,
			bestQuery, bestResponse[i], done, i}, problem)
		bufferFree[i] <- true
	}
	for {
		select {
		case i := <-bufferReady:
			saveBest(&best, buffer[i])
			bufferFree[i] <- true
		case i := <-bestQuery:
			bestResponse[i] <- best.totalCost
		case <-done:
			return best, nil
		case <-timeout:
			return best, nil
		}
	}
}
