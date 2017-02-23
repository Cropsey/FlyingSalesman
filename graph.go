package fsp

import "fmt"

type Graph struct {
	data   map[City]map[Day]map[City]Flight
	source City
}

func NewGraph(source string) Graph {
	graph := new(Graph)
	graph.source = City(source)
	graph.data = make(map[City]map[Day]map[City]Flight)
	return *graph
}

func (g Graph) String() string {
	var s string
	for src, dayList := range g.data {
		for day, dstMap := range dayList {
			for dst, f := range dstMap {
				s = fmt.Sprintf("%s%s->%s %d %d\n", s, src, dst, day, f.cost)
			}
		}
	}
	return s
}

func (g *Graph) AddFlight(e Flight) {
	if g.data[e.from] == nil {
		g.data[e.from] = make(map[Day]map[City]Flight)
	}
	if g.data[e.from][e.day] == nil {
		g.data[e.from][e.day] = make(map[City]Flight)
	}
	current, present := g.data[e.from][e.day][e.to]
	if present {
		if current.cost > e.cost {
			g.data[e.from][e.day][e.to] = e
		}
	} else {
		g.data[e.from][e.day][e.to] = e
	}
}

func (g Graph) Filtered() []Flight {
	filtered := make([]Flight, 0)
	for _, dayList := range g.data {
		for _, dstMap := range dayList {
			for _, flight := range dstMap {
				filtered = append(filtered, flight)
			}
		}
	}
	return filtered
}
