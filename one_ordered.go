package fsp

// engine that tries to find at least one solution in order specified
// (simplification of the problem)
type One_ordered struct{}

func (e One_ordered) Solve(done <-chan struct{}, p Problem) <-chan Solution {
	result := make(chan Solution)
	go func() {
		result <- solve(p)
	}()
	return result
}

func solve(p Problem) Solution {
	stops := stops(p)
	flights := p.flights
	if len(stops) == 0 {
		return Solution{}
	}
	solution := []Flight{}
	for si, current := range stops {
		next := stops[(si+1)%len(stops)]
		found := false
		for i := range flights {
			if flights[i].from == current && flights[i].to == next {
				solution = append(solution, flights[i])
				found = true
				break
			}
		}
		if !found {
			return Solution{}
		}
	}
	return NewSolution(solution)
}
