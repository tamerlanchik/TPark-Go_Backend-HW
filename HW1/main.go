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
	columnCount     int64
}

func main() {
	Interface(os.Args[1:])
}

func Interface(args []string) []string {
	var sourceFile, destFile string
	var params = SortParams{
		isReverse:       false,
		isDelEqual:      false,
		isRegisterIgnor: false,
		isNumeral:       false,
		columnCount:     0,
	}

	for idx := 0; idx < len(args); idx++ {
		val := args[idx]
		if val[0] == '-' {
			switch val {
			case "-f":
				params.isRegisterIgnor = true
			case "-u":
				params.isDelEqual = true
			case "-r":
				params.isReverse = true
			case "-o":
				// TODO: проверка отсутствия имени файла
				if idx+1 < len(args) {
					destFile = args[idx+1]
					idx++
				} else {
					fmt.Println("No destination filename passed")
					os.Exit(1)
				}
			case "-n":
				params.isNumeral = true
			case "-k":
				if idx+1 < len(args) {
					colCount, _ := strconv.ParseInt(args[idx+1], 10, 8)
					params.columnCount = colCount
					idx++
				} else {
					fmt.Println("No column number passed")
					os.Exit(1)
				}
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
			result = a < b
			if params.isReverse {
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
