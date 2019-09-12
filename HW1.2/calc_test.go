package main

import "testing"

type testPair struct {
	input  string
	answer float64
}

var tests = []testPair{
	{"(1+2)-3", 0},
	{"(1+2)*3", 9},
}

func TestCalcInterface(t *testing.T) {
	for idx, val := range tests {
		res := CalcInterface(val.input)
		if res != val.answer {
			t.Errorf("Test %d failed.\n->Got\n %#f\n->instead\n %#f\n", idx+1, res, val.answer)
		}
	}
}
