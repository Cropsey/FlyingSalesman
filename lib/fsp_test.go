package fsp

import "testing"

//import "github.com/Cropsey/fsp/lib"

func TestEmpty(t *testing.T) {
	flights := []Flight{}
	stops := []string{}
	p := Problem{&flights, &stops}
	var d dunno
	s := d.Solve(p)
	if len(s) != 0 {
		t.Error("Empty solution (%v)", s)
	}
}
