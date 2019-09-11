package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type testpair struct {
	input   []string
	ansFile string
}

var tests = []testpair{
	{[]string{directoryName + "/" + "test-0.txt"}, directoryName + "/" + "test-0-ans.txt"},
	{[]string{directoryName + "/" + "test-1.txt", "-r"}, directoryName + "/" + "test-1-ans.txt"},
	{[]string{directoryName + "/" + "test-2.txt", "-n"}, directoryName + "/" + "test-2-ans.txt"},
	{[]string{directoryName + "/" + "test-3.txt", "-n", "-r"}, directoryName + "/" + "test-3-ans.txt"},
	{[]string{directoryName + "/" + "test-4.txt", "-f"}, directoryName + "/" + "test-4-ans.txt"},
	{[]string{directoryName + "/" + "test-5.txt", "-k", "2", "-n"}, directoryName + "/" + "test-5-ans.txt"},
	{[]string{directoryName + "/" + "test-6.txt", "-u"}, directoryName + "/" + "test-6-ans.txt"},
	{[]string{directoryName + "/" + "test-7.txt", "-u", "-f"}, directoryName + "/" + "test-7-ans.txt"},
}

func TestInterface(t *testing.T) {
	for idx, pair := range tests {
		ans := strings.Join(Interface(pair.input), "\r\n")

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
