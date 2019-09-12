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
	maxLevel  int
}

//var expr Expr

type visitor struct {
	level int
	expr  *Expr
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if v.level < v.expr.maxLevel {
		v.expr.calculate()
	}
	v.expr.maxLevel = v.level
	switch n.(type) {
	case *ast.BinaryExpr:
		node := n.(*ast.BinaryExpr)
		v.expr.operators = append(v.expr.operators, node.Op)
	case *ast.BasicLit:
		node := n.(*ast.BasicLit)
		value, _ := strconv.ParseFloat(node.Value, 64)
		v.expr.values = append(v.expr.values, value)
	}
	v.level++
	return v
}

func CalcInterface(expression string) (ans float64) {
	var expr Expr
	var v visitor
	v.expr = &expr
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
