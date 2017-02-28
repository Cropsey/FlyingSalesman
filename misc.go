package fsp

const MAX_CITIES int = 200
const MAX_FLIGHTS int = 27000000

func Cost(flights []Flight) Money {
	var sum Money
	for _, f := range flights {
		sum += f.cost
	}
	return sum
}

// is solution correct? if not, why?
func correct(p Problem, s Solution) (bool, string) {
	var day Day
	for _, f := range s.flights {
		if day > f.day {
			return false, "timing"
		}
		day = f.day + 1
	}
	return true, ""
}

func stops(p Problem) []City {
	m := make(map[City]bool)
	for _, f := range p.flights {
		m[f.from] = true
		m[f.to] = true
	}
	stops := make([]City, 0)
	for c, _ := range m {
		stops = append(stops, c)
	}
	return stops
}

func equal(a, b Flight) bool {
	return a.from == b.from && a.to == b.to && a.day == b.day && a.cost == b.cost
}
