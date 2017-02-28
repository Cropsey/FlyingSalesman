package main

import (
	"bufio"
	"fmt"
	"fsp"
	"os"
	"strconv"
	"strings"
	//"github.com/pkg/profile"
	"runtime/debug"
)

/*var engines = []fsp.Engine {
	fsp.One_places{},
}*/

func readInput() fsp.Problem {
	var src string
	flights := make([]fsp.Flight, 0, fsp.MAX_FLIGHTS)
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() {
		src = stdin.Text()
	}
	for stdin.Scan() {
		// l := strings.Split(stdin.Text(), " ")
		l := customSplit(stdin.Text())
		day, _ := strconv.Atoi(l[2])
		cost, _ := strconv.Atoi(l[3])
		flight := fsp.NewFlight(l[0], l[1], day, cost)
		flights = append(flights, flight)
	}
	p := fsp.NewProblem(src, flights)
	return p
}

func customSplit(s string) []string {
        /* Splits lines of input into 4 parts
           strictly expects format "{3}[A-Z] {3}[A-Z] \d \d"
           WARNING: no checks are done at all */
        r := make([]string, 4)
        r[0] = s[:3]
        r[1] = s[4:7]
        pos2 := strings.LastIndexByte(s, ' ')
        r[2] = s[8:pos2]
        r[3] = s[pos2+1:]
        return r
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
	debug.SetGCPercent(-1)
	//defer profile.Start().Stop()
	problem := readInput()
	graph := fsp.NewGraph(problem)
	var solution fsp.Solution
	var err error

	if len(graph.Filtered()) < 100 {
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
