package main

import (
	"bufio"
	"fmt"
	"fsp"
	//	"github.com/pkg/profile"
	"os"
    "time"
	//	"runtime/debug"
	"strconv"
	"strings"
)

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
	l := make([]string, 4)
	for stdin.Scan() {
		customSplit(stdin.Text(), l)
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

func customSplit(s string, r []string) {
	/* Splits lines of input into 4 parts
	   strictly expects format "{3}[A-Z] {3}[A-Z] \d \d"
	   WARNING: no checks are done at all */
	r[0] = s[:3]
	r[1] = s[4:7]
	pos2 := strings.LastIndexByte(s, ' ')
	r[2] = s[8:pos2]
	r[3] = s[pos2+1:]
}

func main() {
	//  debug.SetGCPercent(-1)
	//	defer profile.Start(profile.MemProfile).Stop()
	timeout := time.After(29 * time.Second)
	problem := readInput()
	solution, err := problem.Solve(timeout)
	if err == nil {
		fmt.Print(solution)
	} else {
		fmt.Println(err)
	}
}
