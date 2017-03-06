package fsp

import "math"

var currentBest = Money(math.MaxInt32)

type NoPath struct{}

func (e NoPath) Error() string {
	return "No path"
}

type AlreadyVisited struct{}

func (e AlreadyVisited) Error() string {
	return "Already visited"
}

type DFSEngine struct {
	graph   Graph
	reverse bool
}

func (d DFSEngine) Name() string {
	return "DFSEngine"
}

func (d DFSEngine) Solve(comm comm, problem Problem) {
	f := make([]Flight, 0, problem.n)
	v := make(map[City]bool)
	partial := partial{v, f, problem.n, 0}
	dst := d.graph.data[0][0]
	for i := 0; i < len(dst); i++ {
		var f *Flight
		if d.reverse {
			f = dst[len(dst)-1-i]
		} else {
			f = dst[i]
		}
		if f == nil {
			continue
		}
		partial.fly(f)
		d.dfsEngine(comm, &partial)
		partial.backtrack()
	}
	comm.done()
}

type partial struct {
	visited map[City]bool
	flights []Flight
	size    int
	cost    Money
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
	p.cost += f.Cost
}

func (p *partial) lastFlight() Flight {
	return p.flights[len(p.flights)-1]
}

func (p *partial) backtrack() {
	f := p.flights[len(p.flights)-1]
	delete(p.visited, f.From)
	p.flights = p.flights[0 : len(p.flights)-1]
	p.cost -= f.Cost
}

func (d DFSEngine) dfsEngine(comm comm, partial *partial) {
	if partial.cost > currentBest {
		return
	}
	if partial.roundtrip() {
		currentBest = comm.sendSolution(Solution{partial.flights, partial.cost})
	}

	lf := partial.lastFlight()
	if partial.hasVisited(lf.To) {
		return
	}

	dst := d.graph.data[lf.To][lf.Day+1]
	for i := 0; i < len(dst); i++ {
		f := dst[i]
		if f == nil {
			continue
		}
		partial.fly(f)
		d.dfsEngine(comm, partial)
		partial.backtrack()
	}
}
