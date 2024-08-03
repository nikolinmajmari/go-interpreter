package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
		let a = 0;
		let b = 1;
		let c = 2;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram returned nil !")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements, ParseProgram() returned  %d! ", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"a"}, {"b"}, {"c"},
	}
	for i, tt := range tests {
		stm := program.Statements[i]
		if !testLetStatement(t, stm, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 123;
	return a;
	return b;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Errorf("Expected %d statements, got %d", 3, len(program.Statements))
	}

	for _, stm := range program.Statements {
		rStm, ok := stm.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement is not an ast.ReturnStatement, got %T", rStm)
		}
		if rStm.TokenLiteral() != "return" {
			t.Errorf("Expected %q token literal, got %q", "return", rStm.TokenLiteral())
		}
	}
}

func TestIdentifierLiteral(t *testing.T) {
	input := "a;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statments, 1 expected got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not an ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier. got %T", stmt.Expression)
	}
	if ident.Value != "a" {
		t.Errorf("ident.Value expected %q got=%q", "a", ident.Value)
	}
	if ident.TokenLiteral() != "a" {
		t.Errorf("ident.TokenLiteral expected 'a' instead got %q", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "1;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not an ast.ExpressionStateent, got %T", stmt)
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral got %T", stmt.Expression)
	}
	if literal.Value != 1 {
		t.Errorf("Expected literal value 1 got %d", literal.Value)
	}
	if literal.TokenLiteral() != "1" {
		t.Errorf("Expected token literal of 1 got %q", literal.TokenLiteral())
	}

}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors:", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error : %q", msg)
	}
	t.FailNow()
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!1;", "!", 1},
		{"-1;", "-", 1},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement, got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement, got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stm.Expression is expected to be *ast.PrefixExpression, got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s, got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5+5;", 5, "+", 5},
		{"5-5;", 5, "-", 5},
		{"5*5;", 5, "*", 5},
		{"5/5;", 5, "/", 5},
		{"5>5;", 5, ">", 5},
		{"5<5;", 5, "<", 5},
		{"5==5;", 5, "==", 5},
		{"5!=5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Errorf("Expected 1 statement, got %d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}
		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Expected statement %q got %q instead", tt.expected, program.String())
		}
	}
}

func testInfixExpression(
	t *testing.T,
	expression ast.Expression,
	left interface{}, operator interface{},
	right interface{}) bool {
	infExp, ok := expression.(ast.InfixExpression)
	if !ok {
		t.Errorf("Expected ast.InfixExpression type, got %T", expression)
		return false
	}
	if !testLiteralExpression(t, infExp.Left, left) {
		return false
	}
	if infExp.Operator != operator {
		t.Errorf("Expected %q as operator got %q instead", operator, infExp.Operator)
	}
	if !testLiteralExpression(t, infExp.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, string(v))
	}
	t.Errorf("Unhandeled expression type %T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is expected *ast.IntegerLiteral got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("Expected integ.Value to be %d, got=%d instead", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Integ token literal expected to be %q, got=%q",
			integ.TokenLiteral(), fmt.Sprintf("%d", value))
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected ast.Identifier got %T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("Expected ident.Value to be %q instead got %q", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("Expected ident.TokenLiteral() to be %q instead got %q", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLetStatement(t *testing.T, stm ast.Statement, name string) bool {
	if stm.TokenLiteral() != "let" {
		t.Errorf("Expected let token literal got %s", stm.TokenLiteral())
		return false
	}
	letStm, ok := stm.(*ast.LetStatement)
	if !ok {
		t.Errorf("Expected let statement, got %T", stm)
		return false
	}
	if letStm.Name.Value != name {
		t.Errorf("LetStatement.Name.Value expects %s, got %s", name, letStm.Name.Value)
		return false
	}
	if letStm.Name.TokenLiteral() != name {
		t.Errorf("LetStatement.Name.TokenLiteral expects %s got %s", name, letStm.Name.TokenLiteral())
		return false
	}
	return true
}
