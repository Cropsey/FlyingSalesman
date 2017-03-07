package fsp

import "fmt"

type Graph struct {
	data   [][][]Flight
	source City
	size   int
}

func NewGraph(problem Problem) Graph {
	graph := new(Graph)
	graph.source = problem.start
	graph.size = problem.n
	filter(problem, graph)
	return *graph
}

func (g Graph) String() string {
	var s string
	for _, dayList := range g.data {
		for _, dstList := range dayList {
			for _, f := range dstList {
				s = fmt.Sprintf("%s%d->%d %d %d\n", s, f.From, f.To, f.Day, f.Cost)
			}
		}
	}
	return s
}

func set(slice [][][]Flight, from City, day Day, flight Flight) {
	if slice[from] == nil {
		slice[from] = make([][]Flight, MAX_CITIES)
	}
	if slice[from][day] == nil {
		slice[from][day] = make([]Flight, 0, MAX_CITIES)
	}
	slice[from][day] = append(slice[from][day], flight)
}

func filter(p Problem, graph *Graph) {
	g := make([][][]Flight, MAX_CITIES)
	for _, f := range p.flights {
		set(g, f.From, f.Day, f)
	}
	graph.data = g
}
