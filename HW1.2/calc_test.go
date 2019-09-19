package main

import "testing"

type testPair struct {
	input  string
	answer float64
}

var tests = []testPair{
	{"(1+2)-3", 0},
	{"( 1 + 2)*3", 9},
	{"4", 4},
	{"2/4", 0.5},
	{"-1", -1},
	{"14 - 5*5 - 3", -14},
	{"1 + 2*(1 + 2*(1+2*(1-0.25)))", 13},
	{"(((1-0.25)*2+1)*2+1)*2+1", 13},
	{"1/3", 1. / 3},
	{"3.2/2 - (3*(7-2*3)/6)/( (1+2+3) / (9-6/2))", 1.1},
}

func TestCalcExpression(t *testing.T) {
	for idx, val := range tests {
		res := CalcExpression(val.input)
		if res != val.answer {
			t.Errorf("Test %d failed.\n->Got\n %#f\n->instead\n %#f\n", idx+1, res, val.answer)
		}
	}
}
