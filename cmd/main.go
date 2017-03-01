package main

import (
	"bufio"
	"fmt"
	"fsp"
	"github.com/pkg/profile"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

/*var engines = []fsp.Engine {
	fsp.One_places{},
}*/

type lookup struct {
	cityToIndex map[string]uint32
	indexToCity []string
}

func getIndex(city string, l *lookup) uint32 {
	ci, found := l.cityToIndex[city]
	if found {
		return ci
	}
	ci = uint32(len(l.cityToIndex))
	l.cityToIndex[city] = ci
	l.indexToCity = append(l.indexToCity, city)
	return ci
}
func readInput() fsp.Problem {
	lookup := &lookup{make(map[string]uint32), make([]string, 0, fsp.MAX_CITIES)}
	flights := make([]fsp.Flight, 0, fsp.MAX_FLIGHTS)

	var src string
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() {
		src = stdin.Text()
		getIndex(src, lookup)
	}
	for stdin.Scan() {
		l := customSplit(stdin.Text())
		day, _ := strconv.Atoi(l[2])
		cost, _ := strconv.Atoi(l[3])
		from := getIndex(l[0], lookup)
		to := getIndex(l[1], lookup)
		flight := fsp.NewFlight(from, to, uint16(day), cost)
		flights = append(flights, flight)
	}
	p := fsp.NewProblem(flights, lookup.indexToCity)
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
	defer profile.Start().Stop()
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
