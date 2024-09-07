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
		t.Errorf("Expected object of *object.Integer type got %T instead", obj)
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

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}
