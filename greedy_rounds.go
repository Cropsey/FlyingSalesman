package fsp

import (
	"container/heap"
	"math"
	"sort"
	"time"
)

type GreedyRounds struct {
	graph       Graph
	currentBest Money
}

func (d GreedyRounds) Name() string {
	return "GreedyRounds"
}

func NewGreedyRounds(g Graph) GreedyRounds {
	return GreedyRounds{graph, Money(math.MaxInt32)}
}

func initStart(g Graph, problem Problem) []fd {
	h := fdHeap(make([]fd, 0, 10))

	for _, fromList := range g.data {
		for _, flights := range fromList {
			for _, f := range flights {
				stat := problem.stats.ByDest[f.From][f.To]
				discount := stat.AvgPrice - float32(f.Cost)
				if len(h) < cap(h) {
					h = append(h, fd{f, discount})
					if len(h) == cap(h) {
						heap.Init(&h)
					}
				} else {
					if h[0].d < discount {
						heap.Pop(&h)
						heap.Push(&h, fd{f, discount})
					}
				}
			}
		}
	}
	return h
}

func (d GreedyRounds) Solve(comm comm, problem Problem) {
	flights := make([]Flight, 0, problem.n)
	visited := make(map[City]bool)
	partial := partial{visited, flights, problem.n, 0}

	for i, f := range initStart(d.graph, problem) {
		printInfo("GreedyRounds start", i, f)
		partial.fly(f.f)
		d.dfs(comm, &partial, time.After(3*time.Second))
		partial.backtrack()
	}
}

func (d *GreedyRounds) dfs(comm comm, partial *partial, timeout <-chan time.Time) bool {
	if expired(timeout) {
		return true
	}
	if partial.cost > d.currentBest {
		return false
	}
	if partial.roundtrip() {
		sf := make([]Flight, partial.n)
		copy(sf, partial.flights)
		sort.Sort(ByDay(sf))
		d.currentBest = comm.sendSolution(NewSolution(sf))
	}

	lf := partial.lastFlight()
	if partial.hasVisited(lf.To) {
		return false
	}

	dst := d.graph.fromDaySortedCost[lf.To][int(lf.Day+1)%d.graph.size]
	for _, f := range dst {
		partial.fly(f)
		expired := d.dfs(comm, partial, timeout)
		partial.backtrack()
		if expired {
			return true
		}
	}
	return false
}

type fd struct {
	f Flight
	d float32
}

type fdHeap []fd

func (h fdHeap) Len() int           { return len(h) }
func (h fdHeap) Less(i, j int) bool { return h[i].d < h[j].d }
func (h fdHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *fdHeap) Push(x interface{}) {
	*h = append(*h, x.(fd))
}
func (h *fdHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
