package fsp

import (
	//"fmt"
	"github.com/emef/bitfield"
)

type Mitm struct{} // meet in the middle

func (m Mitm) Name() string {
	return "MeetInTheMiddle"
}

func (m Mitm) Solve(_ <-chan struct{}, problem Problem) <-chan Solution {
	result := make(chan Solution)
	if problem.n < 2 {
		go func() {
			result <- Solution{}
		}()
		return result
	}
	// processing Problem into two trees
	there, back := makeTwoTrees(problem)
	var mps meetPlaces = make(map[City]meetPlace)

	// we 
	left := make(chan halfRoute)
	right := make(chan halfRoute)

	// run, Forrest!
	go startHalfDFS(left, problem, &there)
	go startHalfDFS(right, problem, &back)

	var found *[]City = nil
	var hr halfRoute
	var ok bool
	for {
		select {
		case hr, ok = <-left:
			//if hr.visited.n == 0 {
			if ! ok {
				left = nil
			} else {
				found = mps.add(true, &hr)
			}
		case hr, ok = <-right:
			if !ok  {
				right = nil
			} else {
				found = mps.add(false, &hr)
			}
		}
		if found != nil || (left == nil && right == nil ) { break }
	}

	var solution Solution
	if found != nil {
		solution = problem.route2solution(*found)
	} else {
	}
	go func() {
		result <- solution
	}()
	return result
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
func (cs citySet) test(c City) bool {
	return cs.data.Test(uint32(c))
}

//TODO this is terrible name, make something better
func (cs citySet) allVisited(other citySet) bool {
	var bi uint32
	for i := 0; i < other.n; i++ {
		bi = uint32(i)
		ob := other.data.Test(bi)
		cb := cs.data.Test(bi)
		if !(cb || ob) {
			return false
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
type flightTree map[Day]map[City]map[City]Money

type halfRoute struct {
	visited citySet
	route   []City
}
type meetPlaces map[City]meetPlace

type meetPlace struct {
	left, right *[]halfRoute
}

// returns route if full route can be costructed, otherwise nil
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
		mp = mps[city]
	}
	hrsCurrent := mp.left
	hrsOther := mp.right
	if !left {
		hrsCurrent = mp.right
		hrsOther = mp.left
	}
	var found *halfRoute = nil
	for _, v := range *hrsOther {
		if v.visited.allVisited(hr.visited) {
			found = &v
		}
	}
	hrsNew := append(*hrsCurrent, *hr)
	hrsCurrent = &hrsNew
	if found != nil {
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
func startHalfDFS(output chan halfRoute, problem Problem, ft *flightTree) {
	defer close(output)

	visited := csInit(problem.n)
	visited.add(problem.start)
	halfDFS(output, []City{problem.start}, visited, 0, Day(len(*ft)), ft)
}

func halfDFS(output chan halfRoute, partial []City, visited citySet, day, endDay Day, ft *flightTree) {
	if day == endDay {
		// we have reached the meeting day
		output <- halfRoute{visited, partial}
		return
	}
	lastVisited := partial[len(partial)-1]
	//TODO not looking at cost at all
	for city, _ := range (*ft)[day][lastVisited] {
		if !visited.test(city) {
			halfDFS(output, append(partial, city),
				visited.add(city),
				day+1, endDay, ft)
		}
	}
	return
}

func addFlight(ft *flightTree, day Day, from, to City, cost Money) {
	if (*ft) == nil {
		(*ft) = make(map[Day]map[City]map[City]Money)
	}
	if (*ft)[day] == nil {
		(*ft)[day] = make(map[City]map[City]Money)
	}
	if (*ft)[day][from] == nil {
		(*ft)[day][from] = make(map[City]Money)
	}
	(*ft)[day][from][to] = cost
}

func makeTwoTrees(problem Problem) (there, back flightTree) {
	// get the number of days
	var days Day = Day(problem.n)
	meetDay := days / 2
	for _, f := range problem.flights {
		if f.Day < meetDay {
			addFlight(&there, f.Day, f.From, f.To, f.Cost)
		} else {
			addFlight(&back, days-1-f.Day, f.To, f.From, f.Cost)
		}
	}
	return
}
