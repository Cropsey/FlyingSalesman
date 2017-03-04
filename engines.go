package fsp

import "math"

type engine interface {
	run(comm comm, buffer *result, task *taskData)
}

type comm struct {
    bufferFree <-chan bool
    bufferReady chan<- int
    id int
}
func (c comm) isFree() {
    <-c.bufferFree
}
func (c comm) resultReady() {
    c.bufferReady <- c.id
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

func initChannels(engines int) []chan bool {
	bufferFree := make([]chan bool, engines)
    for i:=0; i<engines; i++ {
        bufferFree[i] = make(chan bool, 1)
    }
    return bufferFree
}

func initEngines() []engine {
	return []engine{DFSEngine{}}
}

func saveBest(b *result, r result) {
    if b.cost > r.cost {
        for i,f := range r.flights {
            b.flights[i] = f
        }
        b.cost = r.cost
    }
}

func KickTheEngines(task *taskData) (Solution, error) {
	cities := task.problem.cities
	engines := initEngines()

	//signalize goroutine they can write to their buffer
	bufferFree := initChannels(len(engines))
	buffer := initBuffer(len(cities), len(engines))
	best := result{math.MaxInt32, make([]Flight, len(cities))}

	//goroutine with id signals its buffer is ready
	bufferReady := make(chan int, len(engines))

	for i, e := range engines {
		go e.run(comm{bufferFree[i], bufferReady, i}, &buffer[i], task)
        bufferFree[i] <- true
	}
	for {
		select {
		case i := <-bufferReady:
            saveBest(&best, buffer[i])
            bufferFree[i] <- true
		case <-task.timeout:
			return Solution{best.flights, best.cost, cities}, nil
		}
	}
}
