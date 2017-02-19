package fsp

// engine that tries to find at least one solution in order specified
// (simplification of the problem)
type one_ordered struct{}

func (e one_ordered) Solve(p Problem) Solution {
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
}
