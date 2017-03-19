package fsp

import (
	//"container/heap"
	"math"
	"sync"
)

type penalty struct {
	init Money
	m    *sync.Mutex
}

func (p *penalty) save(s partial) {
	p.m.Lock()
	if p.init == 0 {
		p.init = s.cost
	}
	normalized := float64(s.cost) / float64(p.init)
	fraction := normalized / 1000
	for _, f := range s.flights {
		f.Penalty += fraction
	}
	p.m.Unlock()
}

type heuristics func(*Flight) float64

type MetaEngine struct {
	weight []float64
	q      float64
	name   string
	graph  Graph
	h      heuristics
	p      *penalty
}

func (m MetaEngine) Name() string {
	return m.name
}

func (m MetaEngine) Solve(comm comm, problem Problem) {
	flights := make([]*Flight, 0, problem.n)
	visited := make(map[City]bool)
	partial := partial{visited, flights, problem.n, 0}
	for {
		f := nextFlight(m.graph.fromDaySortedCost[0][0], &partial, m.weight[0], m.h)
		partial.fly(f)
		partial.visited[0] = false
		if ok := m.run(&partial); ok {
			comm.sendSolution(NewSolution(partial.solution()))
		}
		m.p.save(partial)
		partial.flights = partial.flights[0:0]
		for k, _ := range partial.visited {
			partial.visited[k] = false
		}
		partial.cost = 0
	}
}

func (m *MetaEngine) run(partial *partial) bool {
	for {
		if partial.roundtrip() {
			return true
		}
		lf := partial.lastFlight()
		d := lf.Day + 1
		dst := m.graph.fromDaySortedCost[lf.To][d]
		nextFlight := nextFlight(dst, partial, m.weight[d], m.h)
		if nextFlight == nil {
			return false
		}
		partial.fly(nextFlight)
	}
}

func nextFlight(flights []*Flight, partial *partial, weight float64, h heuristics) *Flight {
	var cMax Money
	var pMax float64
	var hMax float64
	valid := false
	for _, f := range flights {
		if partial.hasVisited(f.To) {
			continue
		}
		valid = true
		if cMax < f.Cost {
			cMax = f.Cost
		}
		if pMax < f.Penalty {
			pMax = f.Penalty
		}
		hVal := h(f)
		if hMax < hVal {
			hMax = hVal
		}
	}
	if !valid {
		return nil
	}
	if pMax == 0 {
		pMax = 1
	}
	best := float64(math.MaxFloat32)
	var bestFlight *Flight
	for _, f := range flights {
		if partial.hasVisited(f.To) {
			continue
		}
		ncost := float64(f.Cost) / float64(cMax)
		npen := f.Penalty / pMax
		nheur := h(f) / hMax
		val := (1-weight)*ncost + weight*(npen+nheur)
		if best > val {
			best = val
			bestFlight = f
		}
	}

	return bestFlight
}

func initWeight(size int, max float64) []float64 {
	last := size - 2
	diff := max / float64(last)
	weight := make([]float64, size)
	weight[0] = max
	weight[last] = 0
	weight[last+1] = 0
	for i := 1; i < last; i++ {
		weight[i] = weight[i-1] - diff
	}
	return weight
}
