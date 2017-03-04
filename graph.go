package fsp

import "fmt"

type Graph struct {
	data   [][][]*Flight
	source City
	size   int
}

func NewGraph(problem Problem) Graph {
	graph := new(Graph)
	graph.source = problem.start
	graph.size = len(problem.cities)
	filter(problem, graph)
	return *graph
}

func (g Graph) String() string {
	var s string
	for _, dayList := range g.data {
		for _, dstList := range dayList {
			for _, f := range dstList {
				if f != nil {
					s = fmt.Sprintf("%s%d->%d %d %d\n", s, f.From, f.To, f.Day, f.Cost)
				}
			}
		}
	}
	return s
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
		if f.Cost > flight.Cost {
			slice[from][day][to] = &flight
		}
	} else {
		slice[from][day][to] = &flight
	}
}

func filter(p Problem, graph *Graph) {
	g := make([][][]*Flight, MAX_CITIES)
	for _, f := range p.flights {
		set(g, f.From, f.To, f.Day, f)
	}
	graph.data = g
}
