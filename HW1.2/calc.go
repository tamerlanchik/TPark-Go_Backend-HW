package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

var level visitor

func main() {

	/*f, _ := parser.ParseExpr("(1+2)-3")
	var v visitor
	ast.Walk(v, f)
	for len(opers) > 0 {
		calculate(&numbers, &opers)
	}

	if len(os.Args) < 2 {
		fmt.Println("Empty expression")
		os.Exit(1)
	}
	fmt.Println(CalcInterface(os.Args[1]))*/
	fmt.Println(CalcInterface("(1+2)-3"))
}

type Expr struct {
	operators []token.Token
	values    []float64
	maxLevel  visitor
}

var expr Expr

type visitor int

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if v < expr.maxLevel {
		expr.calculate()
	}
	expr.maxLevel = v
	switch n.(type) {
	case *ast.BinaryExpr:
		node := n.(*ast.BinaryExpr)
		expr.operators = append(expr.operators, node.Op)
	case *ast.BasicLit:
		node := n.(*ast.BasicLit)
		value, _ := strconv.ParseFloat(node.Value, 64)
		expr.values = append(expr.values, value)
	}
	//fmt.Printf("%s%T\n", strings.Repeat("\t", int(v)), n)
	return v + 1
}

func CalcInterface(expression string) (ans float64) {
	var v visitor
	f, _ := parser.ParseExpr(expression)
	ast.Walk(v, f)

	for len(expr.operators) > 0 {
		expr.calculate()
	}
	ans = expr.values[0]
	return
}

func generatePolandNotation(source string) (result []string) {

	return
}

func (e *Expr) calculate() {
	oper := e.operators[len(e.operators)-1]
	e.operators = e.operators[:len(e.operators)-1]
	a := e.values[len(e.values)-1]
	b := e.values[len(e.values)-2]
	e.values = e.values[:len(e.values)-2]
	var result float64
	switch oper {
	case token.ADD:
		result = a + b
	case token.SUB:
		result = a - b
	case token.MUL:
		result = a * b
	case token.QUO:
		result = a / b
	default:
		fmt.Println("?")
	}
	e.values = append(e.values, result)
}
