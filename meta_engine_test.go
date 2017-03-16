package fsp

import "testing"

func TestWeight100(t *testing.T){
    t.Log(initWeight(10, 0.8))
}

func check(e, f *Flight, t *testing.T) {
    if !equal(*e, *f) {
        t.Error(e, "!=", f)
    }
}
func checkF32(e, f float32, t *testing.T) {
    if e != f {
        t.Error(e, "!=", f)
    }
}

func TestOrderNoH(t *testing.T){
    flights := []*Flight{
        {0, 1, 1, 1, 0},
        {0, 2, 1, 3, 0},
        {0, 3, 1, 5, 0},
    }
    ordered := order(flights, nil, 0.4)
    check(ordered[0], flights[0], t)
    check(ordered[1], flights[1], t)
    check(ordered[2], flights[2], t)
}

func TestNormalize(t *testing.T){
    a := make([]float32, 3)
    a[0], a[1], a[2] = 5, 1, 10
    normalize(a)
    checkF32(a[0], 0.5, t)
    checkF32(a[1], 0.1, t)
    checkF32(a[2], 1, t)
}

func TestOrderHeurIgnore(t *testing.T){
    flights := []*Flight{
        {0, 1, 1, 1, 0},
        {0, 2, 1, 10, 0},
        {0, 3, 1, 100, 0},
    }
    h := func(fts []*Flight) []float32 {
        x := make([]float32, 0, len(fts))
        for _, f := range fts {
            x = append(x, float32(100-f.Cost))
        }
        return x
    }
    ordered := order(flights, h, 0)
    check(ordered[0], flights[0], t)
    check(ordered[1], flights[1], t)
    check(ordered[2], flights[2], t)
}

func TestOrderHeurReverse(t *testing.T){
    flights := []*Flight{
        {0, 1, 1, 1, 0},
        {0, 2, 1, 10, 0},
        {0, 3, 1, 100, 0},
    }
    h := func(fts []*Flight) []float32 {
        x := make([]float32, 0, len(fts))
        for _, f := range fts {
            x = append(x, float32(101-f.Cost))
        }
        return x
    }
    ordered := order(flights, h, 1)
    check(ordered[0], flights[2], t)
    check(ordered[1], flights[1], t)
    check(ordered[2], flights[0], t)
}

func TestOrderHeurRevWeaker(t *testing.T){
    flights := []*Flight{
        {0, 1, 1, 1, 0},
        {0, 2, 1, 10, 0},
        {0, 3, 1, 100, 0},
    }
    h := func(fts []*Flight) []float32 {
        x := make([]float32, 0, len(fts))
        for _, f := range fts {
            x = append(x, float32(101-f.Cost))
        }
        return x
    }
    //f 0.01,  0.1,   1
    //h 1,     0.91,  0.01
    //s 0.406, 0.424, 0.604
    ordered := order(flights, h, 0.4)
    check(ordered[0], flights[0], t)
    check(ordered[1], flights[1], t)
    check(ordered[2], flights[2], t)
}

func TestOrderHeurRevStronger(t *testing.T){
    flights := []*Flight{
        {0, 1, 1, 1, 0},
        {0, 2, 1, 10, 0},
        {0, 3, 1, 100, 0},
    }
    h := func(fts []*Flight) []float32 {
        x := make([]float32, 0, len(fts))
        for _, f := range fts {
            x = append(x, float32(101-f.Cost))
        }
        return x
    }
    ordered := order(flights, h, 0.6)
    check(ordered[0], flights[2], t)
    check(ordered[1], flights[1], t)
    check(ordered[2], flights[0], t)
}
