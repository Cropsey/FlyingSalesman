package fsp

import "fmt"

type Graph struct {
	problem  Problem
	data     map[City]map[Day]map[City][]int
    filtered []Flight
}

func NewGraph(problem Problem) Graph {
	graph := createMap(problem, problem.flights)
	graph.filterDuplicates()
    graph.populateFiltered()
	return graph
}

func (g Graph) Size() int {
	return len(g.filtered) 
}

func (g Graph) String() string {
	var s string
	for src, dayList := range g.data {
		for day, dstMap := range dayList {
			for dst, flightList := range dstMap {
				for _, i := range flightList {
					f := g.problem.flights[i]
					s = fmt.Sprintf("%s%s->%s %d %d\n", s, src, dst, day, f.cost)
				}
			}
		}
	}
	return s
}

func (g* Graph) populateFiltered() {
	for _, dayList := range g.data {
		for _, dstMap := range dayList {
			for _, flightList := range dstMap {
				for _, i := range flightList {
					f := g.problem.flights[i]
					g.filtered = append(g.filtered, f)
				}
			}
		}
	}
}


func emptyGraph(problem Problem) Graph {
	graph := new(Graph)
	graph.data = make(map[City]map[Day]map[City][]int)
    graph.filtered = make([]Flight, 0)
	graph.problem = problem
	return *graph
}


func (g Graph) addFlight(e Flight, i int) {
	if g.data[e.from] == nil {
		g.data[e.from] = make(map[Day]map[City][]int)
	}
	if g.data[e.from][e.day] == nil {
		g.data[e.from][e.day] = make(map[City][]int)
	}
	if g.data[e.from][e.day][e.to] == nil {
		g.data[e.from][e.day][e.to] = make([]int, 0)
	}
	g.data[e.from][e.day][e.to] = append(g.data[e.from][e.day][e.to], i)
}

func (g Graph) filter(flightList []int) []int {
	if len(flightList) == 0 {
		return flightList
	}
	bestI := 0
	bestCost := g.problem.flights[flightList[bestI]].cost
	for i, index := range flightList[1:] {
		f := g.problem.flights[index]
		cost := f.cost
		if cost < bestCost {
			bestCost = cost
			bestI = i + 1
		}
	}
	return flightList[bestI : bestI+1]
}

func (g Graph) filterDuplicates() {
	for src, dayMap := range g.data {
		for day, dstMap := range dayMap {
			for dst, flightList := range dstMap {
				g.data[src][day][dst] = g.filter(flightList)
			}
		}
	}
}

func createMap(problem Problem, flights []Flight) Graph {
	//go over all flights
	graph := emptyGraph(problem)
	for i, e := range flights {
        //create node where adjacent edges are stored by day
        graph.addFlight(e, i)
	}
	return graph
}
