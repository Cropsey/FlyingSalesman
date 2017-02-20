package main

import "fmt"
import "fsp"

var engines = []fsp.FspEngine {
	fsp.One_places{},
}

func main() {
	done := make(chan struct{})
	defer close(done)
	out := make([]<-chan fsp.Solution, len(engines))
	problem := fsp.ExampleProblem
	for i, e := range engines {
		out[i] = e.Solve(done, problem)
	}
	for i, _ := range engines {
		s := <- out[i]
		fmt.Printf("%v: %v\n", fsp.Cost(problem, s), s)
	}
}
