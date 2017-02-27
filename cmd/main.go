package main

import (
    "fmt"
    "fsp"
    "bufio"
    "os"
    "strings"
    "strconv"
)

/*var engines = []fsp.Engine {
	fsp.One_places{},
}*/

func readInput() fsp.Graph {
    var graph fsp.Graph
    stdin := bufio.NewScanner(os.Stdin)
    if stdin.Scan() {
        src := stdin.Text()
        graph = fsp.NewGraph(src)
    }
    for stdin.Scan() {
        l := strings.Split(stdin.Text(), " ")
        day, _ := strconv.Atoi(l[2])
        cost, _ := strconv.Atoi(l[3])
        flight := fsp.NewFlight(l[0], l[1], day, cost)
        graph.AddFlight(flight)
    }
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
    var err error

    if len(graph.Filtered()) < 10000 {
        solution, err = fsp.DFS(graph)
    } else {
        solution = kickTheEngines(graph)
    }
    if err == nil {
        fmt.Print(solution)
    } else {
        fmt.Println(err)
    }
}
