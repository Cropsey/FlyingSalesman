package fsp

// engine that tries to find at least one solution,
// using DFS, not considering time constraints
type One_places struct{}

func (e One_places) Solve(done <-chan struct{}, p Problem) <-chan Solution {
	result := make(chan Solution)
	go func() {
		stops := stops(p)
		flights := p.flights
		if len(stops) < 2 {
			result <- Solution{}
			return
		}
		// stops = { brq, lon, xxx }
		// visited = { brq }
		visited := make([]City, 1, len(stops))
		visited[0] = stops[0]
		// to_visit = { lon, xxx, brq }
		to_visit := append(stops[1:], stops[0])
		partial := make([]Flight, 0, len(stops))
		result <- NewSolution(one_dfs(partial, visited, to_visit, flights), p.cities)
	}()
	return result
}

func indexOf(haystack []City, needle City) int {
	for i, item := range haystack {
		if item == needle {
			return i
		}
	}
	return -1
}

func one_dfs(partial []Flight, visited, to_visit []City, flights []Flight) []Flight {
	if len(to_visit) == 0 {
		return partial
	}
	for _, f := range flights {
		if f.from == visited[len(visited)-1] {
			if si := indexOf(to_visit, f.to); si != -1 {
				solution := one_dfs(append(partial, f),
					append(visited, f.to),
					append(to_visit[:si], to_visit[si+1:]...),
					flights)
				if len(solution) != 0 {
					// soluton found, yaaaay!
					return solution
				}
			}
		}
	}
	// no solution
	return []Flight{}
}
