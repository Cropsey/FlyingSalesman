package fsp

// engine that tries to find at least one solution in order specified
// (simplification of the problem)
/*type One_ordered struct{}

func (e One_ordered) Solve(done <-chan struct{}, p Problem) <-chan Solution {
	result := make(chan Solution)
	go func() {
		result <- solve(p)
	}()
	return result
}

func solve(p Problem) Solution {
	stops := p.stops
	flights := p.flights
	if len(stops) == 0 {
		return []int{}
	}
	solution := []int{}
	for si, current := range stops {
		next := stops[(si+1)%len(stops)]
		found := false
		for i := range flights {
			if flights[i].from == current && flights[i].to == next {
				solution = append(solution, i)
				found = true
				break
			}
		}
		if !found {
			return []int{}
		}
	}
	return solution
}*/
