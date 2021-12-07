package main

//https://rosettacode.org/wiki/Cartesian_product_of_two_or_more_lists#Go
func cartN(a ...[]int) [][]int {
	c := 1
	for _, a := range a {
		c *= len(a)
	}
	if c == 0 {
		return nil
	}
	p := make([][]int, c)
	b := make([]int, c*len(a))
	n := make([]int, len(a))
	s := 0
	for i := range p {
		e := s + len(a)
		pi := b[s:e]
		p[i] = pi
		s = e
		for j, n := range n {
			pi[j] = a[j][n]
		}
		for j := len(n) - 1; j >= 0; j-- {
			n[j]++
			if n[j] < len(a[j]) {
				break
			}
			n[j] = 0
		}
	}
	return p
}
