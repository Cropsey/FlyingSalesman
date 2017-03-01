package fsp

import "fmt"

type Graph struct {
	data             [][][]*Flight
	hasNegativeEdges bool
	source           City
	size             int
	filtered         []Flight
	problem          Problem
}

func NewGraph(problem Problem) Graph {
	graph := new(Graph)
	graph.source = problem.start
	graph.size = len(problem.cities)
	graph.problem = problem
	filter(problem, graph)
	setFiltered(graph)
	return *graph
}

func (g Graph) String() string {
	var s string
	for _, dayList := range g.data {
		for _, dstList := range dayList {
			for _, f := range dstList {
				if f != nil {
					s = fmt.Sprintf("%s%d->%d %d %d\n", s, f.from, f.to, f.day, f.cost)
				}
			}
		}
	}
	return s
}

func (g Graph) Filtered() []Flight {
	return g.filtered
}

func setFiltered(g *Graph) {
	filtered := make([]Flight, 0, MAX_FLIGHTS)
	for _, dayList := range g.data {
		for _, dstList := range dayList {
			for _, f := range dstList {
				if f != nil {
					filtered = append(filtered, *f)
				}
			}
		}
	}
	g.filtered = filtered
}

func set(slice [][][]*Flight, from, to City, day Day, flight Flight) {
	if slice[from] == nil {
		slice[from] = make([][]*Flight, MAX_CITIES)
	}
	if slice[from][day] == nil {
		slice[from][day] = make([]*Flight, MAX_CITIES)
	}
	f := slice[from][day][to]
	if f != nil {
		if f.cost > flight.cost {
			slice[from][day][to] = &flight
		}
	} else {
		slice[from][day][to] = &flight
	}
}

func filter(p Problem, graph *Graph) {
	g := make([][][]*Flight, MAX_CITIES)
	hasNegativeEdges := false
	for _, f := range p.flights {
		set(g, f.from, f.to, f.day, f)
		if f.cost < 0 {
			hasNegativeEdges = true
		}
	}
	graph.data = g
	graph.hasNegativeEdges = hasNegativeEdges
}
