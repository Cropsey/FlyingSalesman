package main

import (
	"bufio"
	"fmt"
	"github.com/Cropsey/fsp"
	//	"github.com/pkg/profile"
	"os"
	"time"
	//	"runtime/debug"
	"strconv"
	"strings"
)

type lookup struct {
	cityToIndex map[string]fsp.City
	indexToCity []string
}

func getIndex(city string, l *lookup) fsp.City {
	ci, found := l.cityToIndex[city]
	if found {
		return ci
	}
	ci = fsp.City(len(l.cityToIndex))
	l.cityToIndex[city] = ci
	l.indexToCity = append(l.indexToCity, city)
	return ci
}

func readInput() fsp.Problem {
	lookup := &lookup{make(map[string]fsp.City), make([]string, 0, fsp.MAX_CITIES)}
	flights := make([]fsp.Flight, 0, fsp.MAX_FLIGHTS)

	var src string
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() {
		src = stdin.Text()
		getIndex(src, lookup)
	}
	l := make([]string, 4)
	var i int
	var from,to fsp.City
	var day fsp.Day
	var cost fsp.Money
	for stdin.Scan() {
		customSplit(stdin.Text(), l)
		i, _ = strconv.Atoi(l[2])
		day = fsp.Day(i)
		i, _ = strconv.Atoi(l[3])
		cost = fsp.Money(i)
		from = getIndex(l[0], lookup)
		to = getIndex(l[1], lookup)
		flights = append(flights, fsp.Flight{from, to, day, cost})
	}
	p := fsp.NewProblem(flights, len(lookup.indexToCity))
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
	//debug.SetGCPercent(-1)
	//defer profile.Start(profile.MemProfile).Stop()
	timeout := time.After(29 * time.Second)
	problem := readInput()
	solution, err := problem.Solve(timeout)
	if err == nil {
		fmt.Print(solution)
	} else {
		fmt.Println(err)
	}
}
