package parser

import (
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

func TestIdentifier(t *testing.T) {
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
		t.Fatalf("expected *ast.Identifier. got %T", program.Statements[0])
	}
	if ident.Value != "a" {
		t.Errorf("ident.Value expected %q got=%q", "a", ident.Value)
	}
	if ident.TokenLiteral() != "a" {
		t.Errorf("ident.TokenLiteral expected 'a' instead got %q", ident.TokenLiteral())
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
