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
			t.Fatalf("nil value from eval")
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
