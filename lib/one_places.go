package fsp

// engine that tries to find at least one solution,
// not considering time constraints
type one_places struct{}

func (e one_places) Solve(p Problem) Solution {
	stops := p.stops
	flights := p.flights
	if len(stops) < 2 {
		return []int{}
	}
	// stops = { brq, lon, xxx }
	// visited = { brq }
	visited := make([]string, 1, len(stops))
	visited[0] = stops[0]
	// to_visit = { lon, xxx, brq }
	to_visit := append(stops[1:], stops[0])
	partial := make([]int, 0, len(stops))
	solution := solve(partial, visited, to_visit, flights)
	return solution
}

func indexOf(haystack []string, needle string) int {
	for i, item := range haystack {
		if item == needle {
			return i
		}
	}
	return -1
}

func solve(partial []int, visited,to_visit []string, flights []Flight) []int {
	if len(to_visit) == 0 {
		return partial
	}
	for i, f := range flights {
		if f.from == visited[len(visited)-1] {
			if si := indexOf(to_visit, f.to); si != -1 {
					solution := solve(append(partial, i),
					append(visited, f.to),
					append(to_visit[:si], to_visit[si+1:]...),
					flights)
				if len(solution) != 0 {
					return solution
				} else {
					return solve(partial[0:len(partial)-1],
						visited[0:len(visited)-1],
						append(to_visit, f.to),
						flights)
				}
			}
		}
	}
	return []int{}
}
