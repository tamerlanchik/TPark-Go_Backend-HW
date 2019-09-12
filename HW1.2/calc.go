// Source: https://medium.com/justforfunc/understanding-go-programs-with-go-parser-c4e88a6edb87

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
)

var level visitor

func main() {
	/*
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

type visitor struct {
	level int
	expr  *Expr
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if v.level < v.expr.maxLevel {
		for i := 0; i < v.expr.maxLevel-v.level; i++ {
			v.expr.calculate()
		}
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

	case *ast.UnaryExpr:
		node := n.(*ast.UnaryExpr)
		v.expr.operators = append(v.expr.operators, node.Op)
		// для унарного минуса вычитаем
		// реальное положительное число из нуля
		v.expr.values = append(v.expr.values, 0.0)

	// В случае скобочного выражения <ParentExpr>
	// уровень итогового дерева не меняется
	case *ast.ParenExpr:
		v.level--
	}
	//fmt.Printf("%s%T\n", strings.Repeat("\t", int(v.level)), n)
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

func (e *Expr) calculate() {
	oper := e.operators[len(e.operators)-1]
	e.operators = e.operators[:len(e.operators)-1]
	a := e.values[len(e.values)-1]
	b := e.values[len(e.values)-2]
	e.values = e.values[:len(e.values)-2]
	var result float64
	// Order of operands is reverces
	switch oper {
	case token.ADD:
		result = b + a
	case token.SUB:
		result = b - a
	case token.MUL:
		result = b * a
	case token.QUO:
		result = b / a
	default:
		fmt.Println("?")
		os.Exit(1)
	}
	e.values = append(e.values, result)
}
