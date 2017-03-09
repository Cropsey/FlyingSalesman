package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/Cropsey/fsp"
	//	"github.com/pkg/profile"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var argVerbose *bool
var argTimeout *int
var argStats *bool

func printInfo(args ...interface{}) {
	if *argVerbose {
		fmt.Fprintln(os.Stderr, args...)
	}
}

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
	stats[from][to].AvgPrice = (stats[from][to].AvgPrice*float32(stats[from][to].FlightCount) +
		float32(cost)) / float32(stats[from][to].FlightCount+1)
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
		printInfo("Signal received ", sig)

	}

}

func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	go sigHandler()
	start_time := time.Now()
	argTimeout = flag.Int("t", 29, "Maximal time in seconds to run")
	argVerbose = flag.Bool("v", false, "Be verbose and print some info to stderr")
	argStats = flag.Bool("s", false, "Just read input and print some statistics")
	flag.Parse()
	fsp.BeVerbose = *argVerbose
	fsp.StartTime = start_time

	timeout := time.After(time.Duration(*argTimeout) * time.Second)
	problem, lookup := readInput()
	//printLookup(lookup)
	printInfo("Input read ", problem.FlightsCnt(), " flights, after", time.Since(start_time))
	if *argStats {
		printFlightStatistics(lookup, problem)
		return
	}
	solution, err := problem.Solve(timeout)
	if err == nil {
		fmt.Print(printSolution(solution, lookup))
		if *argVerbose {
			fmt.Fprint(os.Stderr, printVerboseSolution(solution, lookup, problem))
		}
	} else {
		fmt.Println(err)
	}
	printInfo("Problem solved after", time.Since(start_time), "with total cost", solution.GetTotalCost())
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

func printVerboseSolution(s fsp.Solution, m []string, p fsp.Problem) string {
	var buffer bytes.Buffer
	buffer.WriteString(s.GetTotalCost().String())
	buffer.WriteString("\n")
	for _, f := range s.GetFlights() {
		from := m[f.From]
		to := m[f.To]
		avg := p.FlightStats()[f.From][f.To].AvgPrice
		perc := float32(f.Cost) / avg * 100.0
		flight := fmt.Sprintf("%s %s %3d %4d [%7.3f%% of avg %7.2f]\n", from, to, f.Day, f.Cost, perc, avg)
		buffer.WriteString(flight)
	}
	return buffer.String()
}

func printFlightStatistics(m []string, p fsp.Problem) {
	for i, r := range p.FlightStats() {
		if i >= p.CitiesCnt() {
			break
		}
		var dests uint16
		var destsDays uint16
		var sum, cheapestCost, mostExpCost float32
		var cheapestDest, mostExpDest fsp.City
		cheapestCost, mostExpCost = math.MaxInt32, 0
		for j, s := range r {
			if s.AvgPrice != 0.0 {
				dests++
				destsDays += s.FlightCount
				sum += s.AvgPrice
				if s.AvgPrice < cheapestCost {
					cheapestCost, cheapestDest = s.AvgPrice, fsp.City(j)
				}
				if s.AvgPrice > mostExpCost {
					mostExpCost, mostExpDest = s.AvgPrice, fsp.City(j)
				}
			}
		}
		avg := sum / float32(dests)
		fmt.Printf("%s: destinations: %3d(%4d), cheap: %s(%7.2f), expensive: %s(%7.2f), avg: %7.2f\n",
			m[i], dests, destsDays, m[cheapestDest], cheapestCost, m[mostExpDest], mostExpCost, avg)
	}
}

func printLookup(m []string) {
	for i, s := range m {
		fmt.Fprintln(os.Stderr, i, "->", s)
	}
}
