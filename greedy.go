package fsp

import (
	"fmt"
	"math"
	"sort"
)

type Greedy struct {
	graph       Graph
	currentBest Money
}

func (d Greedy) Name() string {
	return "Greedy"
}

func NewGreedy(g Graph) Greedy {
	return Greedy{graph, Money(math.MaxInt32)}
}

func (d Greedy) Solve(comm comm, problem Problem) {
	if problem.n <= 10 {
		flights := make([]*Flight, 0, problem.n)
		visited := make(map[City]bool)
		partial := partial{visited, flights, problem.n, 0}

		dst := d.graph.fromDaySortedCost[0][0]
		for _, f := range dst {
			partial.fly(f)
			d.dfs(comm, &partial)
			partial.backtrack()
		}
		comm.done()
	} else {
		printInfo("Greedy not running")
	}
}

type partial struct {
	visited map[City]bool
	flights []*Flight
	n       int
	cost    Money
}

func (p *partial) solution() []Flight {
	flights := make([]Flight, len(p.flights))
	for i, f := range p.flights {
		flights[i] = *f
	}
	sort.Sort(ByDay(flights))
	return flights
}

func (p partial) String() string {
	var str string
	for _, f := range p.flights {
		str = fmt.Sprintf("%s, %d->%d %d %f", str, f.From, f.To, f.Cost, f.Penalty)
	}
	return str
}

func (p *partial) roundtrip() bool {
	ff := p.flights[0]
	lf := p.lastFlight()
	isHome := lf.To == ff.From
	return len(p.visited) == p.n && isHome
}

func (p *partial) hasVisited(c City) bool {
	return p.visited[c]
}

func (p *partial) fly(f *Flight) {
	p.visited[f.From] = true
	p.flights = append(p.flights, f)
	p.cost += f.Cost
}

func (p *partial) lastFlight() *Flight {
	return p.flights[len(p.flights)-1]
}

func (p *partial) backtrack() {
	f := p.flights[len(p.flights)-1]
	delete(p.visited, f.From)
	p.flights = p.flights[0 : len(p.flights)-1]
	p.cost -= f.Cost
}

func (d *Greedy) dfs(comm comm, partial *partial) {
	if partial.cost > d.currentBest {
		return
	}
	if partial.roundtrip() {
		d.currentBest = comm.sendSolution(NewSolution(partial.solution()))
	}

	lf := partial.lastFlight()
	if partial.hasVisited(lf.To) {
		return
	}

	dst := d.graph.fromDaySortedCost[lf.To][int(lf.Day+1)%d.graph.size]
	for _, f := range dst {
		partial.fly(f)
		d.dfs(comm, partial)
		partial.backtrack()
	}
}
