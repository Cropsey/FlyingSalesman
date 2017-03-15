package fsp

import (
	"time"
)

type Polisher struct {
	graph  Graph
	update chan update
}

func (p Polisher) Name() string {
	return "Polisher"
}

func NewPolisher(graph Graph) Polisher {
	return Polisher{
		graph,
		make(chan update, 100),
	}
}

func (p Polisher) try(u update) {
	if len(u.s.flights) < 5 {
		return
	}
	p.update <- u
}

func (p Polisher) Solve(comm comm, problem Problem) {
	for u := range p.update {
		go p.run(comm, u, time.After(3*time.Second))
	}
}

func exists(f *Flight) bool {
	return f != nil
}

/*
0 ---- 1 ---- 2 ---- 3 ---- 4
A      B      C      D      A
a->b   b->c   c->d   d->a
fiPrev fi     fjPrev fj

swap 1 and 3
a->d   d->c   c->b   b->a
giPrev gi     gjPrev gj
*/
func swap(comm comm, g Graph, flights []Flight, i, j int) bool {
	prevI := i - 1
	prevJ := j - 1
	fiPrev := flights[prevI]
	fjPrev := flights[prevJ]
	giPrev := g.get(fiPrev.From, fiPrev.Day, fjPrev.To)
	gjPrev := g.get(fjPrev.From, fjPrev.Day, fiPrev.To)
	fi := flights[i]
	fj := flights[j]
	gi := g.get(fj.From, fi.Day, fi.To)
	gj := g.get(fi.From, fj.Day, fj.To)
	if exists(giPrev) && exists(gjPrev) && exists(gi) && exists(gj) {
		oldF := []Flight{fiPrev, fi, fjPrev, fj}
		newF := []Flight{*giPrev, *gi, *gjPrev, *gj}
		if Cost(oldF) > Cost(newF) {
			swapped := make([]Flight, len(flights))
			for x := 0; x < prevI; x++ {
				swapped[x] = flights[x]
			}
			swapped[prevI] = *giPrev
			swapped[i] = *gi
			for x := i + 1; x < prevJ; x++ {
				swapped[x] = flights[x]
			}
			swapped[prevJ] = *gjPrev
			swapped[j] = *gj
			for x := j + 1; x < len(flights); x++ {
				swapped[x] = flights[x]
			}
			comm.sendSolution(NewSolution(swapped))
			return true
		}
	}
	return false
}

/*
0 ---- 1 ---- 2 ---- 3
A      B      C      A
a->b   b->c   c->a
fiPrev fi     fj

swap 1 and 2
a->c   c->b   b->a
giPrev gi     gj
*/
func swapAdj(comm comm, g Graph, flights []Flight, i, j int) {
	prevI := i - 1
	fiPrev := flights[prevI]
	fi := flights[i]
	fj := flights[j]
	giPrev := g.get(fiPrev.From, fiPrev.Day, fi.To)
	gi := g.get(fj.From, fi.Day, fiPrev.To)
	gj := g.get(fi.From, fj.Day, fj.To)
	if exists(giPrev) && exists(gi) && exists(gj) {
		oldF := []Flight{fiPrev, fi, fj}
		newF := []Flight{*giPrev, *gi, *gj}
		if Cost(oldF) > Cost(newF) {
			swapped := make([]Flight, len(flights))
			for x := 0; x < prevI; x++ {
				swapped[x] = flights[x]
			}
			swapped[prevI] = *giPrev
			swapped[i] = *gi
			swapped[j] = *gj
			for x := j + 1; x < len(flights); x++ {
				swapped[x] = flights[x]
			}
			comm.sendSolution(NewSolution(swapped))
		}
	}
}

func (p Polisher) run(comm comm, u update, timeout <-chan time.Time) {
	n := len(u.s.flights)
	max := n - 1
	for i := 1; i < max; i++ {
		for j := i + 1; j < max; j++ {
			diff := i - j
			if i > j {
				i, j = j, i
			}

			if diff == 1 || diff == -1 {
				swapAdj(comm, graph, u.s.flights, i, j)
			} else {
				swap(comm, graph, u.s.flights, i, j)
			}
		}
	}
	printInfo("polisher done")
}
