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

func sum(flights []Flight) Money {
	var sum Money
	for _, f := range flights {
		sum += f.cost
	}
	return sum
}

func dfs(graph Graph, lastFlight Flight, cost Money, visited map[City]bool) (Money, []Flight, error) {
	visited[lastFlight.from] = true
	//TODO optimize, len() is maybe O(n)
	if len(visited) == len(graph.data) {
		flights := make([]Flight, 0, len(visited))
		flights = append(flights, lastFlight)
		return cost, flights, nil
	}

	if visited[lastFlight.to] {
		return 0, nil, AlreadyVisited{}
	}

	isFirst := true
	var bestCost Money
	var bestFlights []Flight
	var bestFlight Flight
	var bestError error
	bestError = NoPath{}

	for _, f := range graph.data[lastFlight.to][lastFlight.day+1] {

		bc, bf, err := dfs(graph, f, cost+f.cost, visited)
		if err == nil {
			if isFirst == true {
				isFirst, bestFlight, bestCost, bestFlights, bestError = false, f, bc, bf, err
			} else {
				if bc < bestCost {
					bestFlight, bestCost, bestFlights, bestError = f, bc, bf, err
				}
			}
		}
	}
	delete(visited, lastFlight.from)
	if bestError == nil {
		return bestCost, append(bestFlights, bestFlight), nil
	}
	return 0, nil, bestError
}

func DFS(graph Graph) (Solution, error) {
	visited := make(map[City]bool)
	initialCost := Money(0)

	isFirst := true
	var bestCost Money
	var bestFlights []Flight
	var bestError error
	bestError = NoPath{}

	for _, f := range graph.data[graph.source][0] {
		bc, bf, err := dfs(graph, f, initialCost, visited)
		if err == nil {
			if isFirst == true {
				isFirst, bestCost, bestFlights, bestError = false, bc, bf, err
			} else {
				if bc < bestCost {
					bestCost, bestFlights = bc, bf
				}
			}
		}
	}
	return Solution{bestFlights, bestCost}, bestError
}

func NoSearch(graph Graph) Solution {
	flights := graph.Filtered()
	return Solution{flights, sum(flights)}
}
