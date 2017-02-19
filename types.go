package fsp

type Flight struct {
	from string
	to   string
	day  int
	cost int
}

type Problem struct {
	flights []Flight
	stops   []string
}

// flight indices list
type Solution []int

type FspEngine interface {
	Solve(Problem) Solution
}
