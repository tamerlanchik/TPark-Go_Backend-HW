package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type testpair struct {
	input   string
	params  SortParams
	ansFile string
}

var tests = []testpair{
	{
		directoryName + "/" + "test-0.txt",
		SortParams{},
		directoryName + "/" + "test-0-ans.txt",
	},
	{
		directoryName + "/" + "test-1.txt",
		SortParams{isReverse: true},
		directoryName + "/" + "test-1-ans.txt",
	},
	{
		directoryName + "/" + "test-2.txt",
		SortParams{isNumeral: true},
		directoryName + "/" + "test-2-ans.txt",
	},
	{
		directoryName + "/" + "test-3.txt",
		SortParams{isNumeral: true, isReverse: true},
		directoryName + "/" + "test-3-ans.txt",
	},
	{
		directoryName + "/" + "test-4.txt",
		SortParams{isRegisterIgnor: true},
		directoryName + "/" + "test-4-ans.txt",
	},
	{
		directoryName + "/" + "test-5.txt",
		SortParams{isNumeral: true, columnCount: 2},
		directoryName + "/" + "test-5-ans.txt",
	},
	{
		directoryName + "/" + "test-6.txt",
		SortParams{isDelEqual: true},
		directoryName + "/" + "test-6-ans.txt",
	},
	{
		directoryName + "/" + "test-7.txt",
		SortParams{isDelEqual: true, isRegisterIgnor: true},
		directoryName + "/" + "test-7-ans.txt",
	},
}

func TestInterface(t *testing.T) {
	for idx, pair := range tests {
		sortItems := getStringsArray(pair.input)
		ans := strings.Join(sort(sortItems, pair.params), "\r\n")

		file, err := ioutil.ReadFile(pair.ansFile)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		rightAns := string(file)

		if ans != rightAns {
			t.Errorf("Test %d failed.\n->Got\n %s\n->instead\n %s\n", idx+1, ans, rightAns)
		}
	}

}
