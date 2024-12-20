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
		{"(1 > 2) == false", true},
		{"(1 > 2) == true", false},
		{"(1 < 2) == false", false},
		{"(1 < 2) == true", true},
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

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10+2;", 12},
		{"return 4+2;9;", 6},
		{"return 10*2+2;", 22},
		{"11;return 11-5*2;10;", 1},
		{"if(false){return 1;}return 2;", 2},
		{`if (10 > 1) {
					if (10 > 1) {
						return 10;
					}
				return 1;
				}`, 10},
	}
	for _, test := range tests {
		eval := testEval(test.input)
		returnValue, ok := eval.(*object.Integer)
		if !ok {
			t.Fatalf("Expected evaluated object of type *object.ReturnValue got %T instead", eval)
		}
		testIntegerObject(t, returnValue, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"1 + true;1;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"1;true + false;5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if(true){true + false;}", "unknown operator: BOOLEAN + BOOLEAN"},
	}
	for _, test := range tests {
		evaluated := testEval(test.input)
		errorObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("Expected evaluated object of type *object.Error got %T instead", evaluated)
		}
		if test.expected != errorObj.Message {
			t.Fatalf("Expected %q error message got %q instead", test.expected, errorObj.Message)
		}
	}

}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Object is not null, got %T (%v)", obj, obj)
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
