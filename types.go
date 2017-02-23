package fsp

type City string
type Money int
type Day int

type Flight struct {
	from City
	to   City
	day  Day
	cost Money
}

type Problem struct {
	flights []Flight
    start   City
}

type Solution struct {
    flights   []Flight
    totalCost Money
}

/*type FspEngine interface {
	Solve(done <-chan struct{}, problem Problem) <-chan Solution
}*/
