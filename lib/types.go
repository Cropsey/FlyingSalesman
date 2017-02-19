package fsp

type Flight struct {
	from string
	to   string
	day  int
	cost int
}

type Problem struct {
	flights *[]Flight
	stops   *[]string
}

type Solution []int

type Fsp interface {
	Solve(Problem) Solution
}
