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
	if len(os.Args) < 2 {
		fmt.Println("Empty expression")
		os.Exit(1)
	}
	fmt.Println(CalcExpression(os.Args[1]))
}

// С этой инвариантной структурой ходим по дереву
type TreeWalkStorage struct {
	operators []token.Token
	values    []float64
	// храним самый глубокий уровень дерева,
	// чтобы определить начало возврата вверх при обходе дерева
	maxLevel int
}

type visitor struct {
	level int // текущий уровень дерева
	expr  *TreeWalkStorage
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	// Вычисляем поддерево, из которого вышли
	for i := 0; i < v.expr.maxLevel-v.level; i++ {
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

func CalcExpression(expression string) (ans float64) {
	var expr TreeWalkStorage
	var v visitor
	v.expr = &expr
	// Парсинг строки как исходного кода Golang.
	// Построение дерева выражений.
	f, _ := parser.ParseExpr(expression)
	ast.Walk(v, f) // Обход "in depth-first" :(

	// Досчитываем результат после окончания обхода дерева
	for len(expr.operators) > 0 {
		expr.calculate()
	}
	ans = expr.values[0]
	return
}

func (e *TreeWalkStorage) calculate() {
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
