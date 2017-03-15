package fsp

import "sort"

type Graph struct {
	data              [][][]*Flight
	fromDaySortedCost [][][]*Flight
	dayFromData       [][][]*Flight
	fromDayTo         [][][]*Flight
	toDayData         [][][]*Flight
	source            City
	size              int
}

func NewGraph(problem Problem) Graph {
	graph := new(Graph)
	graph.source = problem.start
	graph.size = problem.n
	filter(problem, graph)
	return *graph
}

type byCost []*Flight

func (g Graph) get(from City, day Day, to City) *Flight {
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

func set(slice [][][]*Flight, from City, day Day, flight *Flight) {
	if slice[from] == nil {
		slice[from] = make([][]*Flight, MAX_CITIES)
	}
	if slice[from][day] == nil {
		slice[from][day] = make([]*Flight, 0, MAX_CITIES)
	}
	slice[from][day] = append(slice[from][day], flight)
}
func setcc(slice [][][]*Flight, c1 City, day Day, c2 City, flight Flight) {
	if slice[c1] == nil {
		slice[c1] = make([][]*Flight, MAX_CITIES)
	}
	if slice[c1][day] == nil {
		slice[c1][day] = make([]*Flight, MAX_CITIES)
	}
	slice[c1][day][c2] = &flight
}

func setDayCity(slice [][][]*Flight, day Day, city City, flight *Flight) {
	if slice[day] == nil {
		slice[day] = make([][]*Flight, MAX_CITIES)
	}
	if slice[day][city] == nil {
		slice[day][city] = make([]*Flight, 0, MAX_CITIES)
	}
	slice[day][city] = append(slice[day][city], flight)
	//printInfo("appending", flight, "to [", day, city, "]")
}

func filter(p Problem, graph *Graph) {
	g := make([][][]*Flight, MAX_CITIES)
	fdsc := make([][][]*Flight, MAX_CITIES)
	dtf := make([][][]*Flight, MAX_CITIES)
	fdt := make([][][]*Flight, MAX_CITIES)
	tdf := make([][][]*Flight, MAX_CITIES)
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
		set(g, p.flights[i].From, p.flights[i].Day, &p.flights[i])
		set(fdsc, p.flights[i].From, p.flights[i].Day, &p.flights[i])
		setDayCity(dtf, p.flights[i].Day, p.flights[i].From, &p.flights[i])
		setcc(fdt, p.flights[i].From, p.flights[i].Day, p.flights[i].To, p.flights[i])
		set(tdf, p.flights[i].To, p.flights[i].Day, &p.flights[i])
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
