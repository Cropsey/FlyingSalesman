package fsp

type NoPath struct{}

func (e NoPath) Error() string {
	return "No path"
}

type AlreadyVisited struct{}

func (e AlreadyVisited) Error() string {
	return "Already visited"
}

type DFSEngine struct{}

func (d DFSEngine) run(comm comm, buffer *result, task *taskData) {
	f := make([]Flight, 0, task.problem.n)
	v := make(map[City]bool)
	partial := partial{v, f, task.problem.n}
	for _, f := range task.graph.data[0][0] {
		if f == nil {
			continue
		}
		partial.fly(f)
		dfsEngine(comm, buffer, task.graph, &partial)
		partial.backtrack()
	}
	comm.done()
}

type partial struct {
	visited map[City]bool
	flights []Flight
	size    int
}

func (p *partial) roundtrip() bool {
	lf := p.lastFlight()
	isHome := lf.To == 0
	return len(p.visited) == p.size && isHome
}

func (p *partial) hasVisited(c City) bool {
	return p.visited[c]
}

func (p *partial) fly(f *Flight) {
	p.visited[f.From] = true
	p.flights = append(p.flights, *f)
}

func (p *partial) lastFlight() Flight {
	return p.flights[len(p.flights)-1]
}

func (p *partial) backtrack() {
	f := p.flights[len(p.flights)-1]
	delete(p.visited, f.From)
	p.flights = p.flights[0 : len(p.flights)-1]
}

func sendResult(comm comm, buffer *result, partial *partial) {
	comm.isFree()
	for i := 0; i < len(buffer.flights); i++ {
		buffer.flights[i] = partial.flights[i]
	}
	buffer.cost = Cost(buffer.flights)
	comm.resultReady()
}

func dfsEngine(comm comm, buffer *result, graph Graph, partial *partial) {
	if partial.roundtrip() {
		sendResult(comm, buffer, partial)
	}

	lf := partial.lastFlight()
	if partial.hasVisited(lf.To) {
		return
	}

	for _, f := range graph.data[lf.To][lf.Day+1] {
		if f == nil {
			continue
		}
		partial.fly(f)
		dfsEngine(comm, buffer, graph, partial)
		partial.backtrack()
	}
}
