package evaluator

import (
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1", 1},
		{"2", 2},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		if evaluated == nil {
			t.Fatalf("nil value from eval from input %q", test.input)
		}
		testIntegerObject(t, evaluated, test.expected)
	}
}
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Expected object of *object.Integer type got = %T (%v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected *object.Integer object value to be %d got %d instead", expected, result.Value)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true", true},
		{"false", false},
		{"1<2", true},
		{"1>2", false},
		{"1>1", false},
		{"1<1", false},
		{"1==1", true},
		{"1==2", false},
		{"1!=2", true},
		{"1!=1", false},
	}
	for _, test := range tests {
		evaluated := testEval(test.input)
		if evaluated == nil {
			t.Fatalf("nil value evaluated frm %q input", test.input)
		}
		testBooleanObject(t, evaluated, test.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Expected object of *object.Boolean, got %T instead", obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected *object.Boolean object value to be %t got %t instead", expected, result.Value)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!!true", true},
		{"!5", false},
		{"!false", true},
		{"!!false", false},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestEvalIngerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-1", -1},
		{"1", 1},
		{"-10", -10},
		{"5+5", 10},
		{"2 * 2 + 2 * 3", 10},
		{"-1*3 + 2 * (1 - 5 ) + 13", 2},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, test := range tests {
		evaluated := testEval(test.input)
		if evaluated == nil {
			t.Fatalf("Input %q evaluated as nil", test.input)
		}
		testIntegerObject(t, evaluated, test.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}
