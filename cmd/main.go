package main

import "fmt"
import "fsp"

/*var engines = []fsp.FspEngine {
	fsp.One_places{},
}*/

func readInput() fsp.Problem {
    return fsp.Problem{}
}

func kickTheEngines(graph fsp.Graph) fsp.Solution {
	/*done := make(chan struct{})
	defer close(done)
	out := make([]<-chan fsp.Solution, len(engines))
	problem := fsp.ExampleProblem
	for i, e := range engines {
		out[i] = e.Solve(done, problem)
	}
	for i, _ := range engines {
		s := <-out[i]
		fmt.Printf("%v: %v\n", fsp.Cost(problem, s), s)
	}*/
    return fsp.Solution{}
}

func main() {
    problem := readInput()
    graph := fsp.NewGraph(problem)
    var solution fsp.Solution
    if graph.Size() < 50 {
        solution = fsp.BellmanFord(graph)
    } else {
        solution = kickTheEngines(graph)
    }
    fmt.Println(solution)
}
