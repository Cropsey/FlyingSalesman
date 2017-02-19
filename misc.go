package fsp

func cost(p Problem, s Solution) int {
	sum := 0
	for _, i := range s {
		sum += p.flights[i].cost
	}
	return sum
}

// is solution correct? if not, why?
func correct(p Problem, s Solution) (bool,string) {
	day := 0
	for _, i := range s {
		if day > p.flights[i].day {
			return false, "timing"
		}
		day = p.flights[i].day + 1
	}
	return true, ""
}
