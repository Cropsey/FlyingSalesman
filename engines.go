package fsp

import "math"

type engine interface {
	run(comm comm, buffer *result, task *taskData)
}

type comm struct {
	bufferFree  <-chan Money
	bufferReady chan<- int
	searchedAll chan<- bool
	id          int
}

func (c comm) isFree() {
	<-c.bufferFree
}
func (c comm) resultReady() {
	c.bufferReady <- c.id
}
func (c comm) done() {
	c.searchedAll <- true
}

type result struct {
	cost    Money
	flights []Flight
}

func initBuffer(size, engines int) []result {
	b := make([]result, engines)
	for i, _ := range b {
		b[i] = result{0, make([]Flight, size)}
	}
	return b
}

func initChannels(engines int) []chan Money {
	bufferFree := make([]chan Money, engines)
	for i := 0; i < engines; i++ {
		bufferFree[i] = make(chan Money, 1)
	}
	return bufferFree
}

func initEngines() []engine {
	return []engine{DFSEngine{}}
}

func saveBest(b *result, r result) {
	if b.cost > r.cost {
		for i, f := range r.flights {
			b.flights[i] = f
		}
		b.cost = r.cost
	}
}

func kickTheEngines(task *taskData) (Solution, error) {
	nCities := task.problem.n
	engines := initEngines()

	//signalize goroutine they can write to their buffer
	bufferFree := initChannels(len(engines))
	buffer := initBuffer(nCities, len(engines))
	best := result{math.MaxInt32, make([]Flight, nCities)}

	//goroutine with id signals its buffer is ready
	bufferReady := make(chan int, len(engines))

	//goroutine signals it has searched the entire state space, we can finish
	done := make(chan bool)

	for i, e := range engines {
		go e.run(comm{bufferFree[i], bufferReady, done, i}, &buffer[i], task)
		bufferFree[i] <- best.cost
	}
	for {
		select {
		case i := <-bufferReady:
			saveBest(&best, buffer[i])
			bufferFree[i] <- best.cost
		case <-done:
			return Solution{best.flights, best.cost}, nil
		case <-task.timeout:
			return Solution{best.flights, best.cost}, nil
		}
	}
}
