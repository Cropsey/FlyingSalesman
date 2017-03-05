package fsp

import (
	"fmt"
	"github.com/emef/bitfield"
)

type Mitm struct{} // meet in the middle

func (m Mitm) Name() string {
	return "MeetInTheMiddle"
}

func (m Mitm) Solve(_ <-chan struct{}, problem Problem) <-chan Solution {
	fmt.Println("==========================================")
	result := make(chan Solution)
	if problem.n < 2 {
		go func() {
			result <- Solution{}
		}()
		return result
	}
	there, back := makeTwoTrees(problem)
	fmt.Println("Problem:", problem)
	fmt.Println("There:", there, len(there))
	fmt.Println("Back:", back, len(back))
	var mps meetPlaces = make(map[City]meetPlace)
	fmt.Println("going there")
	visited1 := csInit(problem.n)
	visited1.add(problem.start)
	mp1 := halfDFS([]City{problem.start}, visited1, 0, Day(len(there)), &there)
	fmt.Println("MP1:", mp1)
	mps.add(true, mp1)
	fmt.Println("MPS", mps)
	fmt.Println("going back")
	visited2 := csInit(problem.n)
	visited2.add(problem.start)
	mp2 := halfDFS([]City{problem.start}, visited2, 0, Day(len(back)), &back)
	fmt.Println("MP2:", mp2)
	found := mps.add(false, mp2)
	fmt.Println("MPS", mps)
	fmt.Println("Found:", found)
	var solution Solution
	if found != nil {
		solution = problem.route2solution(*found)
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
	fmt.Println("All visited?", cs.String(), other.String())
	var bi uint32
	for i := 0; i < other.n; i++ {
		bi = uint32(i)
		ob := other.data.Test(bi)
		cb := cs.data.Test(bi)
		fmt.Println("iter:", cs.n, other.n, i, bi, ob, cb)
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

// returns true if full route can be costructed
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
		xxx := meetPlace{&l, &r}
		mps[city] = xxx
		mp = mps[city]
		fmt.Println("after adding")
		fmt.Println(mps[city].left)
		fmt.Println(mps[city].right)
	}
	hrsCurrent := mp.left
	hrsOther := mp.right
	if !left {
		hrsCurrent = mp.right
		hrsOther = mp.left
	}
	var found *halfRoute = nil
	fmt.Println("finding nemo")
	fmt.Println("needle:", hr.visited)
	fmt.Println("haystack:", hrsOther)
	for _, v := range *hrsOther {
		if v.visited.allVisited(hr.visited) {
			found = &v
		}
	}
	fmt.Println("found:", found)
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
			fmt.Println("hrroute", hr.route)
			for i := len((*hr).route) - 2; i >= 1; i-- {
				result = append(result, ((*hr).route)[i])
			}
		}
		return &result
	}
	return nil
}

func halfDFS(partial []City, visited citySet, day, endDay Day, ft *flightTree) *halfRoute {
	fmt.Println("halfDFS:", partial, visited, day, endDay)
	if day == endDay {
		// we have reached the meeting day
		return &halfRoute{visited, partial}
	}
	lastVisited := partial[len(partial)-1]
	//TODO not looking at cost at all
	for city, _ := range (*ft)[day][lastVisited] {
		if !visited.test(city) {
			solution := halfDFS(append(partial, city),
				visited.add(city),
				day+1, endDay, ft)
			if solution != nil {
				return solution
			}
		}
	}
	return nil
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
	fmt.Printf("Adding flight: %v %v %v %v\n", day, from, to, cost)
	(*ft)[day][from][to] = cost
}

func makeTwoTrees(problem Problem) (there, back flightTree) {
	// get the number of days
	var days Day = Day(problem.n)
	fmt.Println("Max day: ", days)
	meetDay := days / 2
	fmt.Println("Meet day: ", meetDay)
	for _, f := range problem.flights {
		if f.Day < meetDay {
			fmt.Println("Adding there..")
			addFlight(&there, f.Day, f.From, f.To, f.Cost)
		} else {
			fmt.Println("Adding back..")
			addFlight(&back, days-1-f.Day, f.To, f.From, f.Cost)
		}
	}
	return
}
