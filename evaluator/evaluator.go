package evaluator

import (
	"interpreter/ast"
	"interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.Boolean:
		return boolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(left, node.Operator, right)
	}
	return nil
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	l, rok := left.(*object.Integer)
	if !rok {
		return nil
	}
	r, lok := right.(*object.Integer)
	if !lok {
		return nil
	}
	switch operator {
	case "+":
		return evalPlusOperatorInfixExpression(l, r)
	case "-":
		return evalMinusOperatorInfixExpression(l, r)
	case "*":
		return evalMultiplyOperatorInfixExpression(l, r)
	case "/":
		return evalDivisionOperatorInfixExpression(l, r)
	default:
		return nil
	}
}

func evalDivisionOperatorInfixExpression(l *object.Integer, r *object.Integer) *object.Integer {
	return &object.Integer{
		Value: l.Value / r.Value,
	}
}

func evalMultiplyOperatorInfixExpression(l *object.Integer, r *object.Integer) *object.Integer {
	return &object.Integer{
		Value: l.Value * r.Value,
	}
}

func evalMinusOperatorInfixExpression(l *object.Integer, r *object.Integer) *object.Integer {
	return &object.Integer{
		Value: l.Value - r.Value,
	}
}

func evalPlusOperatorInfixExpression(left *object.Integer, right *object.Integer) *object.Integer {
	return &object.Integer{
		Value: left.Value + right.Value,
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	}
	return nil
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	value, ok := right.(*object.Integer)
	if !ok || right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	return &object.Integer{
		Value: -value.Value,
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func boolToBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}
