package token

// constants
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	/// Operators
	ASSIGN = "="
	PLUS   = "+"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"

	LBRACKET = "{"
	RBRACKET = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) Type {
	tokenType, ok := keywords[ident]
	if ok {
		return tokenType
	}
	return IDENT
}

type Type string
