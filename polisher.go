package fsp

import (
	"math/rand"
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
	if len(u.solution.flights) < 5 {
		return
	}
	p.update <- u
}

func (p Polisher) Solve(comm comm, problem Problem) {
	for u := range p.update {
		go p.run2(comm, u)
		go p.run3(comm, u)
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

a->d   d->c   c->b   b->a
giPrev gi     gjPrev gj
*/
func swap(comm comm, g Graph, u update, i, j int) {
	flights := u.solution.flights
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
			comm.send(NewSolution(swapped), u.originalEngine)
		}
	}
}

/*
0 ---- 1 ---- 2 ---- 3 ---- 4 ---- 5 ---- 6
A      B      C      D      E      F      G
a->b   b->c   c->d   d->e   e->f   f->g   g->a
fiPrev fi     fjPrev fj     fkPrev fk

a->d   d->c   c->f   f->e   e->b   b->g   g->a
giPrev gi     gjPrev gj     gkPrev gk
*/
func swap3a(comm comm, g Graph, u update, i, j, k int) {
	flights := u.solution.flights
	prevI := i - 1
	prevJ := j - 1
	prevK := k - 1
	fiPrev := flights[prevI]
	fjPrev := flights[prevJ]
	fkPrev := flights[prevK]
	giPrev := g.get(fiPrev.From, fiPrev.Day, fjPrev.To)
	gjPrev := g.get(fjPrev.From, fjPrev.Day, fkPrev.To)
	gkPrev := g.get(fkPrev.From, fkPrev.Day, fiPrev.To)
	fi := flights[i]
	fj := flights[j]
	fk := flights[k]
	gi := g.get(fj.From, fi.Day, fi.To)
	gj := g.get(fk.From, fj.Day, fj.To)
	gk := g.get(fi.From, fk.Day, fk.To)
	if exists(giPrev) && exists(gjPrev) && exists(gkPrev) && exists(gi) && exists(gj) && exists(gk) {
		oldF := []Flight{fiPrev, fi, fjPrev, fj, fkPrev, fk}
		newF := []Flight{*giPrev, *gi, *gjPrev, *gj, *gkPrev, *gk}
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
			for x := j + 1; x < prevK; x++ {
				swapped[x] = flights[x]
			}
			swapped[prevK] = *gkPrev
			swapped[k] = *gk
			for x := k + 1; x < len(flights); x++ {
				swapped[x] = flights[x]
			}
			comm.send(NewSolution(swapped), u.originalEngine)
		}
	}
}
func swap3b(comm comm, g Graph, u update, i, j, k int) {
	flights := u.solution.flights
	prevI := i - 1
	prevJ := j - 1
	prevK := k - 1
	fiPrev := flights[prevI]
	fjPrev := flights[prevJ]
	fkPrev := flights[prevK]
	giPrev := g.get(fiPrev.From, fiPrev.Day, fkPrev.To)
	gjPrev := g.get(fjPrev.From, fjPrev.Day, fiPrev.To)
	gkPrev := g.get(fkPrev.From, fkPrev.Day, fjPrev.To)
	fi := flights[i]
	fj := flights[j]
	fk := flights[k]
	gi := g.get(fk.From, fi.Day, fi.To)
	gj := g.get(fi.From, fj.Day, fj.To)
	gk := g.get(fj.From, fk.Day, fk.To)
	if exists(giPrev) && exists(gjPrev) && exists(gkPrev) && exists(gi) && exists(gj) && exists(gk) {
		oldF := []Flight{fiPrev, fi, fjPrev, fj, fkPrev, fk}
		newF := []Flight{*giPrev, *gi, *gjPrev, *gj, *gkPrev, *gk}
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
			for x := j + 1; x < prevK; x++ {
				swapped[x] = flights[x]
			}
			swapped[prevK] = *gkPrev
			swapped[k] = *gk
			for x := k + 1; x < len(flights); x++ {
				swapped[x] = flights[x]
			}
			comm.send(NewSolution(swapped), u.originalEngine)
		}
	}
}

func (p Polisher) run2(comm comm, u update) {
	start := time.Now()
	n := len(u.solution.flights)
	max := n - 1
	for i := 1; i < max; i++ {
		for j := i + 1; j <= max; j++ {
			swap(comm, graph, u, i, j)
		}
	}
	printInfo("polisher2 done in", time.Since(start))
}

func order(i, j int) (int, int) {
	if i < j {
		return i, j
	}
	return j, i
}

func order3(i, j, k int) (int, int, int) {
	i, j = order(i, j)
	i, k = order(i, k)
	j, k = order(j, k)
	return i, j, k
}

func (p Polisher) run3(comm comm, u update) {
	start := time.Now()
	n := len(u.solution.flights)
	if n < 130 {
		maxi := n - 2
		maxj := n - 1
		for i := 1; i < maxi; i++ {
			for j := i + 1; j < maxj; j++ {
				for k := j + 1; k <= maxj; k++ {
					swap3a(comm, graph, u, i, j, k)
					swap3b(comm, graph, u, i, j, k)
				}
			}
		}
	} else {
		timeout := time.After(3 * time.Second)
		seed := rand.New(rand.NewSource(time.Now().UnixNano()))

		for !expired(timeout) {
			i := seed.Intn(n-1) + 1
			j := seed.Intn(n-1) + 1
			k := seed.Intn(n-1) + 1
			i, j, k = order3(i, j, k)
			if i > j || j > k || i > k {
				//should not happen but in case there is
				//bug again in order3, this prevents bullshit
				continue
			}
			if i == j || j == k {
				continue
			}
			swap3a(comm, graph, u, i, j, k)
			swap3b(comm, graph, u, i, j, k)
		}
	}

	printInfo("polisher3 done in", time.Since(start))
}
