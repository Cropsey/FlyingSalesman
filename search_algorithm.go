package fsp

type NoPath struct {}

func (e NoPath) Error() string {
	return "No path"
}

type AlreadyVisited struct {}

func (e AlreadyVisited) Error() string {
	return "Already visited"
}

func dfs(graph Graph, lf Flight, lastCity int, cost Money, visited map[City]bool) (Money, []Flight, error) {
	visited[lf.from] = true
	defer delete(visited, lf.from)
	if len(visited) == graph.size {
		flights := make([]Flight, 0, len(visited))
		flights = append(flights, lf)
		return cost, flights, nil
	}

	if visited[lf.to] {
		return 0, nil, AlreadyVisited{}
	}

	isFirst := true
	var bestCost Money
	var bestFlights []Flight
	var bestError error
	bestError = NoPath{}

	for dst, f := range graph.data[lastCity][lf.day+1] {
		if f == nil {
			continue
		}
		bc, bf, err := dfs(graph, *f, dst, cost+f.cost, visited)
		if err == nil {
			if isFirst == true {
				isFirst, bestCost, bestFlights, bestError = false, bc, bf, err
			} else {
				if bc < bestCost {
					bestCost, bestFlights, bestError = bc, bf, err
				}
			}
		}
	}
	if bestError == nil {
		return bestCost, append(bestFlights, lf), nil
	}
	return 0, nil, bestError
}

func DFS(task *taskData) (Solution, error) {
	problem := task.problem
	graph := task.graph
	visited := make(map[City]bool)

	isFirst := true
	var bestCost Money
	var bestFlights []Flight
	var bestError error
	bestError = NoPath{}

	for dst, f := range graph.data[0][0] {
		if f == nil {
			continue
		}
		bc, bf, err := dfs(graph, *f, dst, f.cost, visited)
		if err == nil {
			if isFirst == true {
				isFirst, bestCost, bestFlights, bestError = false, bc, bf, err
			} else {
				if bc < bestCost {
					bestCost, bestFlights, bestError = bc, bf, err
				}
			}
		}
	}
	return NewSolution(bestFlights, problem.cities), bestError
}

type partial struct {
    visited map[City]bool
    flights []Flight
    size int
}

func (p *partial) roundtrip() bool {
    return len(p.visited) == p.size
}

func (p *partial) hasVisited(c City) bool {
    return p.visited[c]
}

func (p *partial) fly(f *Flight) {
    p.visited[f.from] = true
    p.flights = append(p.flights, *f)
}

func (p *partial) lastFlight() Flight {
    return p.flights[len(p.flights)-1]
}

func (p *partial) backtrack() {
    f := p.flights[len(p.flights)-1]
    delete(p.visited, f.from)
    p.flights = p.flights[0:len(p.flights)-1]
}

func sendResult(comm comm, buffer *result, partial *partial) {
    comm.isFree()
    for i:=0; i<len(buffer.flights); i++ {
        buffer.flights[i] = partial.flights[i]
    }
    buffer.cost = Cost(buffer.flights)
    comm.resultReady()
}

func dfsEngine(comm comm, buffer *result, graph Graph, partial *partial)  {
	if partial.roundtrip() {
        sendResult(comm, buffer, partial)
	}

    lf := partial.lastFlight()
	if partial.hasVisited(lf.to) {
		return
	}

	for _, f := range graph.data[lf.to][lf.day+1] {
		if f == nil {
			continue
		}
        partial.fly(f)
		dfsEngine(comm, buffer, graph, partial)
        partial.backtrack()
	}
}

type DFSEngine struct {}

func (d DFSEngine) run(comm comm, buffer *result, task *taskData) {
    f := make([]Flight, 0, len(task.problem.cities))
	v := make(map[City]bool)
    partial := partial{v, f, len(task.problem.cities)}
	for _, f := range task.graph.data[0][0] {
		if f == nil {
			continue
		}
        partial.fly(f)
        dfsEngine(comm, buffer, task.graph, &partial)
        partial.backtrack()
    }
}
