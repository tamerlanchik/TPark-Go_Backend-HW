package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	SortLib "sort"
	"strconv"
	"strings"
)

const (
	directoryName = "tests"
)

type SortParams struct {
	isReverse       bool
	isDelEqual      bool
	isRegisterIgnor bool
	isNumeral       bool
	columnCount     int64
}

func main() {
	//Interface(os.Args[1:])
	var sourceFile, destFile string
	var params = SortParams{}

	sourceFile = os.Args[1]
	os.Args = os.Args[1:] // Remove source filename from the parser view

	flag.BoolVar(&params.isReverse, "r", false, "Sort in reversed order")
	flag.BoolVar(&params.isDelEqual, "u", false, "Sort deleting coincident elements")
	flag.BoolVar(&params.isRegisterIgnor, "f", false, "Sort no taking register in account")
	flag.BoolVar(&params.isNumeral, "n", false, "Sort as numeral values")
	flag.Int64Var(&params.columnCount, "k", 0, "Sort by column")
	flag.StringVar(&destFile, "o", "", "Srite result into the file")

	flag.Parse()

	sortPack := getStringsArray(sourceFile)

	sortPack = sort(sortPack, params)
	if destFile != "" {
		file, err := os.Create(destFile)
		if err != nil {
			fmt.Println("Cannot create the destination file")
			os.Exit(1)
		}
		defer file.Close()

		file.WriteString(strings.Join(sortPack, "\r\n"))

	} else {
		fmt.Println(strings.Join(sortPack, "\r\n"))
	}
	fmt.Println(sortPack)
}

func sort(elems []string, params SortParams) []string {

	if params.isDelEqual {
		deleteDublicates(&elems, params.isRegisterIgnor)
	}

	SortLib.Slice(elems, func(i, j int) (result bool) {
		a := strings.Split(elems[i], " ")[params.columnCount]
		b := strings.Split(elems[j], " ")[params.columnCount]

		if params.isRegisterIgnor {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}

		if params.isNumeral {
			aNum, _ := strconv.ParseFloat(a, 64)
			bNum, _ := strconv.ParseFloat(b, 64)
			result = aNum < bNum
			if params.isReverse {
				result = aNum > bNum
			}
		} else {
			if params.isReverse {
				result = a > b
			} else {
				result = a < b
			}
		}
		return
	})

	return elems
}

// Читает файл построчно и возвращает массив его строк
func getStringsArray(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	ans := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ans = append(ans, scanner.Text())
	}

	return ans
}

func deleteDublicates(sl *[]string, isRegisterIgnored bool) {
	set := make(map[string]bool)
	count := 0
	slice := *sl

	for idx, val := range slice {
		if isRegisterIgnored {
			val = strings.ToLower(val)
		}
		if _, ex := set[val]; ex {
			count++
			slice[idx] = slice[len(slice)-count]
		} else {
			set[val] = true
		}
	}

	*sl = slice[:len(slice)-count]

}
