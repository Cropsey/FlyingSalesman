package main

import "fmt"
import "fsp"

/*var engines = []fsp.FspEngine {
	fsp.One_places{},
}*/

func readInput() fsp.Graph {
    graph := fsp.NewGraph()
    //for each line from STDIN
    //  flight := parse(line)
    //  graph.AddFlight(flight)
    return graph 
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
    graph := readInput()
    var solution fsp.Solution
    if len(graph.Filtered()) < 50 {
        solution = fsp.BellmanFord(graph)
    } else {
        solution = kickTheEngines(graph)
    }
    fmt.Println(solution)
}
