package fsp

import (
	//"container/heap"
	//"math"
	"sort"
    "time"
)

type result struct {
    i   int
    val float32
}

type heuristics func([]*Flight) []float32

type MetaEngine struct {
    weight []float32
    q          float32
    name       string
    timeout    time.Duration
    graph      Graph
    h          heuristics
}

type byVal []result
func (x byVal) Len() int { return len(x) }
func (x byVal) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byVal) Less(i, j int) bool { return x[i].val < x[j].val }

func (m MetaEngine) Name() string {
    return m.name
}

func (m MetaEngine) Solve(comm comm, problem Problem) {
	flights := make([]Flight, 0, problem.n)
	visited := make(map[City]bool)
	partial := partial{visited, flights, problem.n, 0}

    dst := m.graph.fromDaySortedCost[0][0]
	for i, f := range dst {
		printInfo(m.Name(), "starts", i, f)
        var t <-chan time.Time
        if m.timeout == 0 {
            t = nil
        } else {
            t = time.After(m.timeout*time.Second)
        }
		partial.fly(*f)
		m.run(comm, &partial, t)
		partial.backtrack()
	}
}

func normalize(x []float32) {
    max := x[0]
    for _, v := range x[1:] {
        if max < v {
            max = v
        }
    }
    for i, v := range x {
        x[i] = v/max
    }
}

func order(flights []*Flight, h heuristics, weight float32) []*Flight {
    if len(flights) == 0 {
        return flights
    }
    if h == nil {
        return flights
    }
    fts := make([]float32, 0, len(flights))
    for _, f := range flights {
        fts = append(fts, float32(f.Cost))
    }
    heur := h(flights)
    normalize(heur)
    normalize(fts)
    s := make([]result, 0, len(flights))
    for i, _ := range flights {
        v := (1-weight)*fts[i] + weight*heur[i]
        s = append(s, result{i, v})
    }

    sort.Sort(byVal(s))
    ordered := make([]*Flight, 0, len(flights))
    for _, x := range s {
        ordered = append(ordered, flights[x.i])
    }
    return ordered
}

func (m *MetaEngine) run(comm comm, partial *partial, timeout <-chan time.Time) bool {
    if expired(timeout) {
        return true
    }
	if partial.cost > best.totalCost {
		return false
	}
	if partial.roundtrip() {
		sf := make([]Flight, partial.n)
		copy(sf, partial.flights)
		sort.Sort(ByDay(sf))
        comm.sendSolution(NewSolution(sf))
	}
	lf := partial.lastFlight()
	if partial.hasVisited(lf.To) {
		return false
	}

    d := int(lf.Day+1)%m.graph.size
    flights := m.graph.fromDaySortedCost[lf.To][d]
    dst := order(flights, m.h, m.weight[d])
	for _, f := range dst {
		partial.fly(*f)
		expired := m.run(comm, partial, timeout)
		partial.backtrack()
		if expired {
			return true
		}
	}
    return false
}

func initWeight(size int, max float32) []float32{
    last := size - 2
    diff := max/float32(last)
    weight := make([]float32, size)
    weight[0] = max
    weight[last] = 0
    weight[last+1] = 0
    for i := 1; i < last; i++ {
        weight[i] = weight[i-1] - diff
    }
    return weight
}
