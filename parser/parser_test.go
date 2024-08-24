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

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		assertProgramLength(t, program, 1)
		stmt := assertExpressionStatement(t, program.Statements[0])
		testBooleanLiteral(t, stmt.Expression, tt.value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		Value    interface{}
	}{
		{"!1;", "!", 1},
		{"-1;", "-", 1},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		assertProgramLength(t, program, 1)
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
		if !testLiteralExpression(t, exp.Right, tt.Value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5+5;", 5, "+", 5},
		{"5-5;", 5, "-", 5},
		{"5*5;", 5, "*", 5},
		{"5/5;", 5, "/", 5},
		{"5>5;", 5, ">", 5},
		{"5<5;", 5, "<", 5},
		{"5==5;", 5, "==", 5},
		{"5!=5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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
		{
			"true", "true",
		},
		{
			"false", "false",
		},
		{
			"true == false",
			"(true == false)",
		},
		{
			"1 > 2 == true",
			"((1 > 2) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(1 + 2) * 3",
			"((1 + 2) * 3)",
		},
		{
			"-2/(1 + 2)",
			"((-2) / (1 + 2))",
		},
		{
			"!( true == false )",
			"(!(true == false))",
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

func TestIfExpression(t *testing.T) {
	input := `if (a > b) { a; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	assertProgramLength(t, program, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])

	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expected expression of type *ast.IfExpression got %T instead", stmt.Expression)
	}
	if !testInfixExpression(t, ifExp.Condition, "a", ">", "b") {
		return
	}
	if len(ifExp.Consequence.Statements) != 1 {
		t.Errorf("Expected 1 statement on consequence, got %d instead", len(ifExp.Consequence.Statements))
	}
	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] from consequence is not an *ast.ExpressionStatement got %T instead",
			ifExp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "a") {
		return
	}
	if ifExp.Alternative != nil {
		t.Fatalf("Expected nil alternative got %v", ifExp.Alternative)
	}
}
func TestIfElseExpression(t *testing.T) {
	input := `if (a<b) {a;} else {b;}`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	assertProgramLength(t, program, 1)
	stmt := assertExpressionStatement(t, program.Statements[0])
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expected expression to be of type *ast.IfExpression got %T instead", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "a", "<", "b") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Expected 1 statement as consequence, got %d", len(exp.Consequence.Statements))
	}
	if exp.Alternative == nil {
		t.Fatalf("Expected alternative on statement")
	}
	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("Expected 1 alternative statement, got %d", len(exp.Alternative.Statements))
	}
	consequence := assertExpressionStatement(t, exp.Consequence.Statements[0])
	if !testIdentifier(t, consequence.Expression, "a") {
		return
	}
	alternative := assertExpressionStatement(t, exp.Alternative.Statements[0])
	if !testIdentifier(t, alternative.Expression, "b") {
		return
	}
}

func testInfixExpression(
	t *testing.T,
	expression ast.Expression,
	left interface{}, operator interface{},
	right interface{}) bool {
	infExp, ok := expression.(*ast.InfixExpression)
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
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, string(v))
	case bool:
		return testBooleanLiteral(t, exp, bool(v))
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

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("Expected *ast.Boolean, got %T instead", exp)
		return false
	}
	if boolean.Value != value {
		t.Errorf("Expected expression boolean value to be %t got %t instead", value, boolean.Value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("Expected expression literal to be %q, got %q instead",
			boolean.TokenLiteral(), fmt.Sprintf("%t", value))
		return false
	}
	return true
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(a, b){ a + b }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	assertProgramLength(t, program, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected program statement to be *ast.ExpressionStatement got %T instead", program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expected expression to be *ast.FunctionLiteral got %T instead", stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("Expected function to have 2 parameters, got %d parameters instead", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "a")
	testLiteralExpression(t, function.Parameters[1], "b")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in function body got %d instead", len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected body statement to be *ast.ExpressionStatement, got %T instead", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "a", "+", "b")
	return
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(a) {};", expectedParams: []string{"a"}},
		{input: "fn(a, b, c) {};", expectedParams: []string{"a", "b", "c"}},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		assertProgramLength(t, program, 1)
		stm, ok := program.Statements[0].(*ast.ExpressionStatement)
		function, ok := stm.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("Expected expression to be of type Function go %T instead", stm.Expression)
		}
		if len(function.Parameters) != len(test.expectedParams) {
			t.Fatalf("Expected %d params, got %d instead", len(test.expectedParams), len(function.Parameters))
		}
		for i, literal := range test.expectedParams {
			testLiteralExpression(t, function.Parameters[i], literal)
		}
	}
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

func assertProgramLength(t *testing.T, program *ast.Program, expectedLength int) {
	t.Helper()
	if len(program.Statements) != expectedLength {
		t.Fatalf("Expected program length to be %d, instead got %d", expectedLength, len(program.Statements))
	}
}

func assertExpressionStatement(t *testing.T, statement ast.Statement) *ast.ExpressionStatement {
	t.Helper()
	stmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not an ast.ExpressionStateent, got %T", stmt)
	}
	return stmt
}
