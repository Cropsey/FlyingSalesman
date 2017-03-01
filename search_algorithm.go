package fsp

type NoPath struct {
	msg string
}

func (e NoPath) Error() string {
	return "No path"
}

type AlreadyVisited struct {
	msg string
}

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

func DFS(graph Graph) (Solution, error) {
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
	return NewSolution(bestFlights, graph.problem), bestError
}

func NoSearch(graph Graph) Solution {
	flights := graph.Filtered()
	return Solution{flights, Cost(flights), graph.problem}
}
