package fsp

func Cost(p Problem, s Solution) int {
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

var ExampleProblem = Problem{
		[]Flight{
			{"aaa", "bbb", 1, 0},
			{"aaa", "bbb", 2, 0},
			{"aaa", "bbb", 3, 0},
			{"lon", "xxx", 2, 0}, // 3
			{"bbb", "ccc", 1, 0},
			{"bbb", "ccc", 2, 0},
			{"bbb", "ccc", 3, 0},
			{"xxx", "brq", 3, 0}, // 7
			{"ccc", "aaa", 1, 0},
			{"ccc", "aaa", 2, 0},
			{"ccc", "aaa", 3, 0},
			{"brq", "lon", 1, 0}, // 11
		},
		[]string{"brq", "lon", "xxx"},
	}
