package fsp

import (
	"math"
	"sort"
    "time"
)

type Bottleneck struct {
	graph       Graph
	currentBest Money
}

func NewBottleneck(g Graph) Bottleneck {
	return Bottleneck{
		graph,
		Money(math.MaxInt32),
	}

}

func (d Bottleneck) Name() string {
	return "Bottleneck"
}

type byCost2 []Flight

func (f byCost2) Len() int {
	return len(f)
}
func (f byCost2) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byCost2) Less(i, j int) bool {
	return f[i].Cost < f[j].Cost
}

func (d Bottleneck) Solve(comm comm, problem Problem) {
	flights := make([]*Flight, 0, problem.n)
	visited := make(map[City]bool)
	partial := partial{visited, flights, problem.n, 0}
    btn := d.findBottlenecks(problem)
    t := 30000.0 / float32(len(btn))
    printInfo("Found",len(btn),"bottlenecks")
	for _, b := range btn {
        timePerBtn := t / float32(min(len(btn),3))
        printInfo("Testing",min(len(btn),3),"flights in bottleneck")
		sort.Sort(byCost2(b))
		for _, f := range b {
			partial.fly(&f)
            tb := time.Duration(timePerBtn)*time.Millisecond
            printInfo("running dfs from bottleneck for",tb )
			d.dfs(comm, &partial, time.After(tb))
			partial.backtrack()
		}
	}
}

type byCount [][]Flight

func (f byCount) Len() int {
	return len(f)
}
func (f byCount) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f byCount) Less(i, j int) bool {
	return len(f[i]) < len(f[j])
}

type btnStat struct {
	from    [][]Flight
	to      [][]Flight
	noBFrom []bool
	noBTo   []bool
	cutoff  int
}

func initB(n int) btnStat {
	b := btnStat{}
	b.from = make([][]Flight, n)
	b.to = make([][]Flight, n)
	b.cutoff = n/4
	for i := range b.from {
		b.from[i] = make([]Flight, 0, b.cutoff)
		b.to[i] = make([]Flight, 0, b.cutoff)
	}
	b.noBFrom = make([]bool, n)
	b.noBTo = make([]bool, n)
	return b
}

func (b *btnStat) add(f Flight) {
	if !b.noBFrom[f.From] {
		b.from[f.From] = append(b.from[f.From], f)
		if len(b.from[f.From]) > b.cutoff {
			b.from[f.From] = nil
			b.noBFrom[f.From] = true
		}
	}
	if !b.noBTo[f.To] {
		b.to[f.To] = append(b.to[f.To], f)
		if len(b.to[f.To]) > b.cutoff {
			b.to[f.To] = nil
			b.noBTo[f.To] = true
		}
	}
}

func (b btnStat) get() [][]Flight {
	all := make([][]Flight, 0, len(b.from)+len(b.to))
	for _, f := range b.from {
		if f != nil && len(f) > 0 {
			all = append(all, f)
		}
	}
	for _, f := range b.to {
		if f != nil && len(f) > 0 {
			all = append(all, f)
		}
	}
	sort.Sort(byCount(all))
	return all
}

func (b *Bottleneck) findBottlenecks(p Problem) [][]Flight {
	bs := initB(p.n)
	for _, f := range p.flights {
		if f.From == 0 || f.To == 0 {
			continue
		}
		bs.add(f)
	}
	return bs.get()
}

func (b *Bottleneck) dfs(comm comm, partial *partial, timeout <-chan time.Time) bool {
    if expired(timeout) {
        return true
    }
	if partial.cost > b.currentBest {
		return false
	}
	if partial.roundtrip() {
		b.currentBest = comm.sendSolution(NewSolution(partial.solution()))
	}

	lf := partial.lastFlight()
	if partial.hasVisited(lf.To) {
		return false
	}

	dst := b.graph.fromDaySortedCost[lf.To][int(lf.Day+1)%b.graph.size]
	for _, f := range dst {
		partial.fly(f)
        expired := b.dfs(comm, partial, timeout)
        if expired {
            return true
        }
		partial.backtrack()
	}
    return false
}
