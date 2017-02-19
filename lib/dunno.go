package fsp

type dunno struct{}

func (d dunno) Solve(_ Problem) Solution {
	return []int{}
}
