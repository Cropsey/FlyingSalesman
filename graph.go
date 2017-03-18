package fsp

import "sort"

type Graph struct {
	data              [][][]FlightIndex
	fromDaySortedCost [][][]FlightIndex
	dayFromData       [][][]FlightIndex
	fromDayTo         [][][]FlightIndex
	toDayData         [][][]FlightIndex
	source            City
	size              int
	problem		  *Problem
}

func NewGraph(problem Problem) Graph {
	graph := new(Graph)
	graph.source = problem.start
	graph.size = problem.n
	filter(problem, graph)
	graph.problem = &problem
	return *graph
}

type byCost []FlightIndex

func (g Graph) get(from City, day Day, to City) FlightIndex {
	if g.fromDayTo[from] == nil {
		return nil
	}
	if g.fromDayTo[from][day] == nil {
		return nil
	}
	return g.fromDayTo[from][day][to]
}

//type byCost []Flight

func (f byCost) Len() int {
	return len(f)
}
func (f byCost) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byCost) Less(i, j int) bool {
	return f[i].Cost < f[j].Cost
}

func set(slice [][][]FlightIndex, from City, day Day, index FlightIndex) {
	if slice[from] == nil {
		slice[from] = make([][]FlightIndex, MAX_CITIES)
	}
	if slice[from][day] == nil {
		slice[from][day] = make([]FlightIndex, 0, MAX_CITIES)
	}
	slice[from][day] = append(slice[from][day], index)
}
func setcc(slice [][][]FlightIndex, c1 City, day Day, c2 City, index FlightIndex) {
	if slice[c1] == nil {
		slice[c1] = make([][]FlightIndex, MAX_CITIES)
	}
	if slice[c1][day] == nil {
		slice[c1][day] = make([]FlightIndex, MAX_CITIES)
	}
	slice[c1][day][c2] = index
}

func setDayCity(slice [][][]FlightIndex, day Day, city City, flight FlightIndex) {
	if slice[day] == nil {
		slice[day] = make([][]FlightIndex, MAX_CITIES)
	}
	if slice[day][city] == nil {
		slice[day][city] = make([]FlightIndex, 0, MAX_CITIES)
	}
	slice[day][city] = append(slice[day][city], flight)
	//printInfo("appending", flight, "to [", day, city, "]")
}

func filter(p Problem, graph *Graph) {
	g := make([][][]FlightIndex, MAX_CITIES)
	fdsc := make([][][]FlightIndex, MAX_CITIES)
	dtf := make([][][]FlightIndex, MAX_CITIES)
	fdt := make([][][]FlightIndex, MAX_CITIES)
	tdf := make([][][]FlightIndex, MAX_CITIES)
	lastDay := Day(graph.size - 1)
	for i, _ := range p.flights {
		if p.flights[i].To == 0 && p.flights[i].Day != lastDay {
			// no need to append paths to home city before last day
			continue
		}
		if p.flights[i].To != 0 && p.flights[i].Day == lastDay {
			// no need to append paths to another city on last day
			continue
		}
		set(g, p.flights[i].From, p.flights[i].Day, i)
		set(fdsc, p.flights[i].From, p.flights[i].Day, i)
		setDayCity(dtf, p.flights[i].Day, p.flights[i].From, i)
		setcc(fdt, p.flights[i].From, p.flights[i].Day, p.flights[i].To, i)
		set(tdf, p.flights[i].To, p.flights[i].Day, i)
	}
	for _, dayList := range fdsc {
		for _, flightList := range dayList {
			sort.Sort(byCost(flightList))
		}
	}
	graph.data = g
	graph.fromDaySortedCost = fdsc
	graph.dayFromData = dtf

	/*
		for i, x := range g {
			for j, y := range x {
				for k, z := range y {
					printInfo("[", i, j, k, "]:", *z)
					}
				}
			}
	*/
	graph.fromDayTo = fdt
	graph.toDayData = tdf
}
