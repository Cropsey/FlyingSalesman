package fsp

import (
	"math"
	"sort"
)

var currentBest = Money(math.MaxInt32)

type Greedy struct {
	graph Graph
}

func (d Greedy) Name() string {
	return "Greedy"
}

func (d Greedy) Solve(comm comm, problem Problem) {
	dst := d.graph.fromDaySortedCost[0][0]
	for _, f := range dst {
		flights := make([]Flight, 0, problem.n)
		visited := make(map[City]bool)
		partial := partial{visited, flights, problem.n, 0}
		partial.fly(f)
		dfs(comm, &d.graph, &partial)
		partial.backtrack()
	}
	comm.done()
}

type Bottleneck struct {
	graph Graph
}

func (d Bottleneck) Name() string {
	return "Bottleneck"
}

func (d Bottleneck) Solve(comm comm, problem Problem) {
	flights := make([]Flight, 0, problem.n)
	visited := make(map[City]bool)
	partial := partial{visited, flights, problem.n, 0}
	for _, b := range findBottlenecks(&d.graph, problem) {
		for _, f := range b.flights {
			partial.fly(f)
			dfs(comm, &d.graph, &partial)
			partial.backtrack()
		}
	}
	comm.done()
}

type fcs struct {
	from    City
	to      City
	flights []Flight
}

type byCount []fcs

func (f byCount) Len() int {
	return len(f)
}
func (f byCount) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byCount) Less(i, j int) bool {
	return len(f[i].flights) < len(f[j].flights)
}

func findBottlenecks(g *Graph, p Problem) []fcs {

	fromToF := make([][]fcs, g.size)
	for from := range fromToF {
		fromToF[from] = make([]fcs, g.size)
		for to := range fromToF[from] {
			fromToF[from][to] = fcs{City(from), City(to), make([]Flight, 0, g.size)}
		}
	}
	for _, f := range p.flights {
		fromToF[f.From][f.To].flights = append(fromToF[f.From][f.To].flights, f)
	}

	b := make([]fcs, 0, MAX_FLIGHTS)
	for from := range fromToF {
		for _, stats := range fromToF[from] {
			b = append(b, stats)
		}
	}
	sort.Sort(byCount(b))
	return b
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

func dfs(comm comm, graph *Graph, partial *partial) {
	if partial.cost > currentBest {
		return
	}
	if partial.roundtrip() {
		sf := make([]Flight, partial.n)
		copy(sf, partial.flights)
		sort.Sort(ByDay(sf))
		currentBest = comm.sendSolution(NewSolution(sf))
	}

	lf := partial.lastFlight()
	if partial.hasVisited(lf.To) {
		return
	}

	dst := graph.fromDaySortedCost[lf.To][int(lf.Day+1)%graph.size]
	for _, f := range dst {
		partial.fly(f)
		dfs(comm, graph, partial)
		partial.backtrack()
	}
}
