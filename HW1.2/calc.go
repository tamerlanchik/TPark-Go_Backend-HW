package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
)

var operations = map[byte]int{
	'(': 0,
	'+': 1,
	'*': 2,
	')': 5,
}

var numbers = make([]float64, 0)
var opers = make([]token.Token, 0)
var level visitor

func main() {

	f, _ := parser.ParseExpr("(1+2)-3")
	/*ast.Inspect(f, func(n ast.Node) bool {
		fmt.Println(n)
		ast.V
		return true
	})*/
	var v visitor
	//numbers := make([]float64, 0)
	ast.Walk(v, f)
	fmt.Println(numbers)
	fmt.Println(opers)
	for len(opers) > 0 {
		calculate(&numbers, &opers)
	}

	if len(os.Args) < 2 {
		fmt.Println("Empty expression")
		os.Exit(1)
	}
	fmt.Println(CalcInterface(os.Args[1]))
}

type visitor int

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if v < level {
		fmt.Println(opers)
		fmt.Println(numbers)
		calculate(&numbers, &opers)
	}
	level = v
	switch n.(type) {
	case *ast.BinaryExpr:
		node := n.(*ast.BinaryExpr)
		opers = append(opers, node.Op)
		//fmt.Println("X= " + string(node.X.(type)) + " Y = " + n.Y)
		/*a := numbers[len(numbers)-1]
		b := numbers[len(numbers)-2]
		numbers = numbers[:len(numbers)-2]
		var result float64*/
		/*switch node.Op {
		case token.ADD:
			fmt.Println("+")
			//result = a + b
		case token.SUB:
			fmt.Println("-")
			//result = a - b
		case token.MUL:
			fmt.Println("*")
			//result = a * b
		case token.QUO:
			fmt.Println("/")
			//result = a / b
		default:
			fmt.Println("?")
		}*/
		//numbers = append(numbers, result)
		//fmt.Print(" " + string(node.Op))
	case *ast.BasicLit:
		node := n.(*ast.BasicLit)
		value, _ := strconv.ParseFloat(node.Value, 64)
		numbers = append(numbers, value)
		fmt.Println(node.Value)
	}
	fmt.Printf("%s%T\n", strings.Repeat("\t", int(v)), n)
	return v + 1
}

func CalcInterface(expression string) (ans float64) {

	return
}

func generatePolandNotation(source string) (result []string) {

	return
}

func calculate(values *[]float64, operators *[]token.Token) {
	//vs := *values
	//opers := *operators

	oper := opers[len(opers)-1]
	opers = opers[:len(opers)-1]
	a := numbers[len(numbers)-1]
	b := numbers[len(numbers)-2]
	numbers = numbers[:len(numbers)-2]
	var result float64
	switch oper {
	case token.ADD:
		fmt.Println("+")
		result = a + b
	case token.SUB:
		fmt.Println("-")
		result = a - b
	case token.MUL:
		fmt.Println("*")
		result = a * b
	case token.QUO:
		fmt.Println("/")
		result = a / b
	default:
		fmt.Println("?")
	}
	numbers = append(numbers, result)

	//*values = vs
	//*operators = opers
}

/*
func calculateExpression(expression []string) (answer float64, resultErr error) {
	var stack = make([]string, 0, len(expression))

	for _, val := range expression {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			stack = append(stack, val)
			continue
		}

		a, _ := strconv.ParseFloat(stack[len(stack)-1], 64)
		b, _ := strconv.ParseFloat(stack[len(stack)-2], 64)
		stack = stack[:len(stack)-2]

		var result float64
		switch val {
		case "+":
			result = a + b
		case "-":
			result = a - b
		case "*":
			result = a * b
		case "/":
			result = a / b
		default:
			resultErr = fmt.Errorf("Wrong operator")
			return
		}
		fmt.Println(result)
		//stack = append(stack, string(result))

	}
	return
}
*/
