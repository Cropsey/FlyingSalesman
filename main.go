package fsp

import "fmt"

var engines = []FspEngine {
	one_places{},
}

func main() {
	done := make(chan struct{})
	defer close(done)
	out := make([]<-chan Solution, len(engines))
	problem := Problem{}
	for i, e := range engines {
		out[i] = e.Solve(done, problem)
	}
	for i, _ := range engines {
		s := <- out[i]
		fmt.Printf("%v: %v\n", cost(problem, s), s)
	}
}
