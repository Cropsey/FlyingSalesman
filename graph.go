package fsp

import "fmt"

type Graph struct {
	data             [][][]*Flight
	hasNegativeEdges bool
	source           City
	cityToIndex      map[City]int
	indexToCity      []City
	filtered         []Flight
}

func NewGraph(problem Problem) Graph {
	graph := new(Graph)
	graph.source = problem.start
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
					s = fmt.Sprintf("%s%s->%s %d %d\n", s, f.from, f.to, f.day, f.cost)
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

func getIndex(city City, cityToIndex map[City]int, indexToCity []City) int {
	ci, found := cityToIndex[city]
	if found {
		return ci
	}
	ci = len(cityToIndex)
	cityToIndex[city] = ci
	indexToCity = append(indexToCity, city)
	return ci
}

func set(slice [][][]*Flight, from, to int, day Day, flight Flight) {
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
	cityToIndex := make(map[City]int)
	indexToCity := make([]City, 0, MAX_CITIES)
	getIndex(p.start, cityToIndex, indexToCity)
	g := make([][][]*Flight, MAX_CITIES)
	hasNegativeEdges := false
	for _, f := range p.flights {
		cif := getIndex(f.from, cityToIndex, indexToCity)
		cit := getIndex(f.to, cityToIndex, indexToCity)
		set(g, cif, cit, f.day, f)
		if f.cost < 0 {
			hasNegativeEdges = true
		}
	}
	graph.data = g
	graph.hasNegativeEdges = hasNegativeEdges
	graph.cityToIndex = cityToIndex
	graph.indexToCity = indexToCity
}
