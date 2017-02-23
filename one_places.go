package fsp

// engine that tries to find at least one solution,
// not considering time constraints
type One_places struct{}

func (e One_places) Solve(done <-chan struct{}, p Problem) <-chan Solution {
	result := make(chan Solution)
	go func() {
		result <- solveTODO2(p)
	}()
	return result
}

func solveTODO2(p Problem) Solution {
	stops := stops(p)
	flights := p.flights
	if len(stops) < 2 {
		return Solution{}
	}
	// stops = { brq, lon, xxx }
	// visited = { brq }
	visited := make([]City, 1, len(stops))
	visited[0] = stops[0]
	// to_visit = { lon, xxx, brq }
	to_visit := append(stops[1:], stops[0])
	partial := make([]Flight, 0, len(stops))
	solution := solveTODO(partial, visited, to_visit, flights)
	return NewSolution(solution)
}

func indexOf(haystack []City, needle City) int {
	for i, item := range haystack {
		if item == needle {
			return i
		}
	}
	return -1
}

func solveTODO(partial []Flight, visited, to_visit []City, flights []Flight) []Flight {
	if len(to_visit) == 0 {
		return partial
	}
	for _, f := range flights {
		if f.from == visited[len(visited)-1] {
			if si := indexOf(to_visit, f.to); si != -1 {
				solution := solveTODO(append(partial, f),
					append(visited, f.to),
					append(to_visit[:si], to_visit[si+1:]...),
					flights)
				if len(solution) != 0 {
					// soluton found, yaaaay!
					return solution
				} else {
					// dead end, let's continue the loop
					partial = partial[0 : len(partial)-1]
					visited = visited[0 : len(visited)-1]
					to_visit = append(to_visit, f.to)
				}
			}
		}
	}
	// no solution
	return []Flight{}
}
