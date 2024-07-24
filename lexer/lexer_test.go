package lexer

import (
	"fmt"
	"interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let a = 6  ;
	let b = 10;
	let sum = fn(x,y){
		x + y;
	}
	let result = sum(a,b);
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "b"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "sum"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACKET, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACKET, "}"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "sum"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		fmt.Printf("%q\n", tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q on %v",
				i, tt.expectedType, tok.Type, []string{string(tok.Type), tok.Literal})
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected %q, got %q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
