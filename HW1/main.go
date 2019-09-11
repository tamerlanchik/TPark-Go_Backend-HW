package main

import (
	"bufio"
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
}

func main() {
	fmt.Println(os.Args[1:])
	Interface(os.Args[1:])
}

func Interface(args []string) []string {
	var sourceFile, destFile string
	var params = map[string]int64{
		"isReverse":  0,
		"isDelEqual": 0,
		"isRegIgnor": 0,
		"isNum":      0,
		"colCount":   0,
	}

	for idx := 0; idx < len(args); idx++ {
		val := args[idx]
		if val[0] == '-' {
			switch val {
			case "-f":
				params["isRegIgnor"] = 1
			case "-u":
				params["isDelEqual"] = 1
			case "-r":
				params["isReverse"] = 1
			case "-o":
				// TODO: проверка отсутствия имени файла
				destFile = args[idx+1]
				idx++
			case "-n":
				params["isNum"] = 1
			case "-k":
				colCount, _ := strconv.ParseInt(args[idx+1], 10, 8)
				params["colCount"] = colCount
				idx++
			default:
				fmt.Println("Wrong argument")
				os.Exit(1)
			}
		} else {
			if sourceFile == "" {
				sourceFile = val
			} else {
				fmt.Println("Source file already written")
				os.Exit(1)
			}
		}
	}

	fmt.Println(params, destFile, sourceFile)

	sortPack := getStringsArray(sourceFile)

	sortPack = sort(sortPack, params)
	if destFile != "" {
		file, err := os.Create(destFile)
		if err != nil {
			fmt.Println("Cannot create destination file")
			os.Exit(1)
		}
		defer file.Close()

		file.WriteString(strings.Join(sortPack, "\r\n"))

	} else {
		fmt.Println(strings.Join(sortPack, "\r\n"))
	}
	return sortPack
}

func sort(elems []string, params map[string]int64) []string {
	SortLib.Slice(elems, func(i, j int) (result bool) {
		a := strings.Split(elems[i], " ")[params["colCount"]]
		b := strings.Split(elems[j], " ")[params["colCount"]]

		if params["isRegIgnor"] == 1 {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}

		if params["isNum"] == 1 {
			aNum, _ := strconv.ParseFloat(a, 64)
			bNum, _ := strconv.ParseFloat(b, 64)
			result = aNum < bNum
			if params["isReverse"] != 0 {
				result = aNum > bNum
			}
		} else {
			result = a < b
			if params["isReverse"] != 0 {
				result = a > b
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
