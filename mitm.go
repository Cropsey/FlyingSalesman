package fsp

import (
	"fmt"
	"github.com/emef/bitfield"
	"math"
)

type Mitm struct{} // meet in the middle

func (m Mitm) Name() string {
	return "MeetInTheMiddle"
}

func printTree(ft *flightTree) {
	for day, d1 := range *ft {
		fmt.Println("day", day)
		for f, d2 := range d1 {
			fmt.Println("  from", f)
			for _, t := range d2 {
				fmt.Println("    ", t.to, t.cost)
			}
		}
	}
}

func printMps(mps map[City]meetPlace) {
	for k, mp := range mps {
		fmt.Println("city", k)
		fmt.Println("  left", len(*mp.left))
		for _, hr := range *mp.left {
			fmt.Println("    ", hr.visited.String(), hr.route, hr.cost)
		}
		fmt.Println("  right", len(*mp.right))
		for _, hr := range *mp.right {
			fmt.Println("    ", hr.visited.String(), hr.route, hr.cost)
		}
	}
}

func (m Mitm) Solve(comm comm, problem Problem) {
	if problem.n < 2 {
		comm.sendSolution(Solution{})
		return
	}
	// processing Problem into two trees
	there, back := makeTwoTrees(problem)
	/*
		fmt.Println("there:")
		printTree(&there)
		fmt.Println("-----")
		fmt.Println("back:")
		printTree(&back)
	*/

	var mps meetPlaces = make(map[City]meetPlace)

	// we
	left := make(chan halfRoute)
	right := make(chan halfRoute)

	// run, Forrest!
	go startHalfDFS(left, problem, &there, true)
	go startHalfDFS(right, problem, &back, false)

	var found *[]City = nil
	var hr halfRoute
	var ok bool
	var bestCost Money = Money(math.MaxInt32)
	var solution Solution
	for {
		select {
		case hr, ok = <-left:
			if !ok {
				left = nil
			} else {
				found = mps.add(true, &hr)
			}
		case hr, ok = <-right:
			if !ok {
				right = nil
			} else {
				found = mps.add(false, &hr)
			}
		}
		if found != nil {
			solution = problem.route2solution(*found)
			if solution.totalCost < bestCost {
				bestCost = solution.totalCost
				comm.sendSolution(solution)
			}
			found = nil
			//fmt.Println("solution sent")
		}
		if left == nil && right == nil {
			comm.done()
			/*
				fmt.Println("mps:", mps)
				printMps(mps)
			*/
			break
		}
	}
}

type citySet struct {
	n    int
	data bitfield.BitField
}

func csInit(n int) (cs citySet) {
	cs.n = n
	cs.data = bitfield.New(n)
	return
}
func (cs citySet) add(c City) citySet {
	cs.data.Set(uint32(c))
	return cs
}
func (cs citySet) remove(c City) citySet {
	cs.data.Clear(uint32(c))
	return cs
}
func (cs citySet) test(c City) bool {
	return cs.data.Test(uint32(c))
}
func (cs citySet) copy() citySet {
	data := bitfield.New(cs.n)
	copy(data, cs.data)
	return citySet{cs.n, data}
}

//TODO this is terrible name, make something better
func (cs citySet) allVisited(other citySet, meetIndex int) bool {
	var bi uint32
	// we are starting from 1 deliberately, start city
	// should be visited in both
	for i := 1; i < other.n; i++ {
		bi = uint32(i)
		ob := other.data.Test(bi)
		cb := cs.data.Test(bi)
		if i == meetIndex {
			if !(cb && ob) { // and
				return false
			}
		} else {
			if !((cb || ob) && !(cb && ob)) { // xor
				return false
			}
		}
	}
	return true
}
func (cs citySet) full() bool { //naive, could be more efficient
	for i := 0; i < cs.n; i++ {
		if !cs.data.Test(uint32(i)) {
			return false
		}
	}
	return true
}
func (cs citySet) String() string {
	res := make([]byte, cs.n)
	for i := 0; i < cs.n; i++ {
		if cs.data.Test(uint32(i)) {
			res[i] = '1'
		} else {
			res[i] = '0'
		}
	}
	return string(res)
}

//IDEA
// could be probably optimized to
// map[Day]map[City][int]
// where those ints are indexes to Problem.Flights sorted by cost
type flightTree map[Day]map[City][]flightTo

type halfRoute struct {
	visited citySet
	route   []City
	cost    Money
}
type meetPlaces map[City]meetPlace

type meetPlace struct {
	left, right *[]halfRoute
}

// returns route if full route can be constructed, otherwise nil
func (mps meetPlaces) add(left bool, hr *halfRoute) *[]City {
	city := (*hr).route[len((*hr).route)-1]
	mp, present := mps[city]
	if !present {
		l := []halfRoute{}
		r := []halfRoute{}
		if left {
			l = append(l, *hr)
		} else {
			r = append(r, *hr)
		}
		mps[city] = meetPlace{&l, &r}
		return nil // no chance for matching meetPlace here
	}
	hrsCurrent := mp.left
	hrsOther := mp.right
	if !left {
		hrsCurrent = mp.right
		hrsOther = mp.left
	}
	*hrsCurrent = append(*hrsCurrent, *hr)

	bestCost := Money(math.MaxInt32)
	// TODO consider cost
	var found *halfRoute = nil
	for _, v := range *hrsOther {
		if v.visited.allVisited(hr.visited, int(city)) {
			if v.cost < bestCost {
				found = &v
				bestCost = v.cost
			}
		}
	}
	if found != nil {
		//fmt.Println("found:", found.visited.String(), found.route, found.cost)
		result := make([]City, 0, (*hr).visited.n)
		if left {
			result = append(result, hr.route...)
			for i := len((*found).route) - 2; i >= 1; i-- {
				result = append(result, (*found).route[i])
			}
		} else {
			result = append(result, (*found).route...)
			for i := len((*hr).route) - 2; i >= 1; i-- {
				result = append(result, ((*hr).route)[i])
			}
		}
		return &result
	}
	return nil
}

// wrapper around halfDFS
func startHalfDFS(output chan halfRoute, problem Problem, ft *flightTree, left bool) {
	defer close(output)

	visited := csInit(problem.n)
	visited.add(problem.start)
	halfDFS(output, []City{problem.start}, visited, 0, Day(len(*ft)), 0, ft, left)
}

func halfDFS(output chan halfRoute, partial []City, visited citySet, day, endDay Day, cost Money, ft *flightTree, left bool) {
	if day == endDay {
		// we have reached the meeting day
		//fmt.Println("route:", left, visited.String(), partial, cost)
		route := make([]City, len(partial))
		copy(route, partial)
		output <- halfRoute{visited.copy(), route, cost}
		return
	}
	lastVisited := partial[len(partial)-1]
	//TODO not looking at cost at all
	for _, fl := range (*ft)[day][lastVisited] {
		city := fl.to
		if !visited.test(city) {
			halfDFS(output, append(partial, city),
				visited.add(city),
				day+1, endDay, cost+fl.cost, ft, left)
			visited.remove(city)
		}
	}
	return
}

type flightTo struct {
	to   City
	cost Money
}

func addFlight(ft *flightTree, day Day, from, to City, cost Money, n int) {
	if (*ft) == nil {
		(*ft) = make(map[Day]map[City][]flightTo)
	}
	if (*ft)[day] == nil {
		(*ft)[day] = make(map[City][]flightTo)
	}
	if (*ft)[day][from] == nil {
		//(*ft)[day][from] = make(map[City]Money)
		(*ft)[day][from] = make([]flightTo, 0, n)
	}
	insertIndex := 0
	for _, v := range (*ft)[day][from] {
		if cost < v.cost {
			break
		}
		insertIndex++
	}
	(*ft)[day][from] = append((*ft)[day][from][:insertIndex],
		append([]flightTo{flightTo{to, cost}},
			(*ft)[day][from][insertIndex:]...)...)
	//(*ft)[day][from][to] = cost
}

func makeTwoTrees(problem Problem) (there, back flightTree) {
	// get the number of days
	var days Day = Day(problem.n)
	meetDay := days / 2
	for _, f := range problem.flights {
		if f.Day < meetDay {
			addFlight(&there, f.Day, f.From, f.To, f.Cost, problem.n)
		} else {
			addFlight(&back, days-1-f.Day, f.To, f.From, f.Cost, problem.n)
		}
	}
	return
}
