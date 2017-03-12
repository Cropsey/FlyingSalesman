package fsp

import (
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

func initStart(g Graph, problem Problem) Flight {
	var bestDiscount float32
	var bestFlight Flight
	for _, fromList := range g.data {
		for _, flights := range fromList {
			for _, f := range flights {
				stat := problem.stats.ByDest[f.From][f.To]
				discount := stat.AvgPrice - float32(f.Cost)
				if discount > bestDiscount {
					bestDiscount = discount
					bestFlight = f
				}
			}
		}
	}
	return bestFlight
}

func (d Greedy) Solve(comm comm, problem Problem) {
	flights := make([]Flight, 0, problem.n)
	visited := make(map[City]bool)
	partial := partial{visited, flights, problem.n, 0}

	f := initStart(d.graph, problem)
	printInfo("Greedy start", f)
	partial.fly(f)
	d.dfs(comm, &partial)
	partial.backtrack()

	dst := d.graph.fromDaySortedCost[0][0]
	for _, f := range dst {
		partial.fly(f)
		d.dfs(comm, &partial)
		partial.backtrack()
	}
	comm.done()
}

type partial struct {
	visited map[City]bool
	flights []Flight
	n       int
	cost    Money
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

func (p *partial) fly(f Flight) {
	p.visited[f.From] = true
	p.flights = append(p.flights, f)
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

func (d *Greedy) dfs(comm comm, partial *partial) {
	if partial.cost > d.currentBest {
		return
	}
	if partial.roundtrip() {
		sf := make([]Flight, partial.n)
		copy(sf, partial.flights)
		sort.Sort(ByDay(sf))
		d.currentBest = comm.sendSolution(NewSolution(sf))
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
