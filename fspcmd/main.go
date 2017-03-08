package main

import (
	"bufio"
	"bytes"
	"fmt"
	"flag"
	"github.com/Cropsey/fsp"
	//	"github.com/pkg/profile"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
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

func readInput() (fsp.Problem, []string) {
	lookup := &lookup{make(map[string]fsp.City), make([]string, 0, fsp.MAX_CITIES)}
	flights := make([]fsp.Flight, 0, fsp.MAX_FLIGHTS)
	stats := make([][]fsp.FlightStats, fsp.MAX_CITIES)
	for s := range stats {
		stats[s] = make([]fsp.FlightStats, fsp.MAX_CITIES)
	}

	var src string
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() {
		src = stdin.Text()
		getIndex(src, lookup)
	}
	l := make([]string, 4)
	var i int
	var from, to fsp.City
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
		updateStats(stats, from, to, cost)
		if from == fsp.City(0) && day != 0 {
			// ignore any flight from src city not on the first day
			// fmt.Fprintln(os.Stderr, "Dropping flight", l)
			continue
		}
		if day == 0 && from != fsp.City(0) {
			// also flights originating in different than home city are wasteful
			// fmt.Fprintln(os.Stderr, "Dropping flight", l)
			continue
		}
		flights = append(flights, fsp.Flight{from, to, day, cost})
	}
	p := fsp.NewProblem(flights, len(lookup.indexToCity), stats)
	return p, lookup.indexToCity
}

func updateStats(stats [][]fsp.FlightStats, from, to fsp.City, cost fsp.Money) {
	if stats[from][to].BestPrice == fsp.Money(0) || stats[from][to].BestPrice > cost {
		stats[from][to].BestPrice = cost
	}
	stats[from][to].AvgPrice = (stats[from][to].AvgPrice*float32(stats[from][to].FlightCount) + float32(cost)) / float32(stats[from][to].FlightCount+1)
	stats[from][to].FlightCount += 1

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

func sigHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	//fmt.Fprintln(os.Stderr, "Signal handler running")
	for sig := range sigs {
		fmt.Fprintln(os.Stderr, "Signal received ", sig)

	}

}

func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	go sigHandler()
	start_time := time.Now()
	argTimeout := flag.Int("t", 29, "Maximal time to run")
	flag.Parse()

	timeout := time.After(time.Duration(*argTimeout) * time.Second)
	problem, lookup := readInput()
	//printLookup(lookup)
	fmt.Fprintln(os.Stderr, "Input read ", problem.FlightsCnt(), " flights, after", time.Since(start_time))
	solution, err := problem.Solve(timeout)
	if err == nil {
		fmt.Print(printSolution(solution, lookup))
	} else {
		fmt.Println(err)
	}
	fmt.Fprintln(os.Stderr, "Problem solved after", time.Since(start_time))
}

func printSolution(s fsp.Solution, m []string) string {
	var buffer bytes.Buffer
	buffer.WriteString(s.GetTotalCost().String())
	buffer.WriteString("\n")
	for _, f := range s.GetFlights() {
		from := m[f.From]
		to := m[f.To]
		flight := fmt.Sprintf("%s %s %d %d\n", from, to, f.Day, f.Cost)
		buffer.WriteString(flight)
	}
	return buffer.String()
}

func printLookup(m []string) {
	for i, s := range m {
		fmt.Fprintln(os.Stderr, i, "->", s)
	}
}
