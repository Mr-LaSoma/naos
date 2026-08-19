package tokens

const (
	EOFRune rune = 0
)

type Token struct {
	Kind   TokenKind
	Lexeme string

	Start Position
	End   Position
}

func (t Token) IsValid() bool {
	return t.Kind.IsValid() &&
		t.Start.IsValid() && t.End.IsValid() &&
		t.Start.FileName == t.End.FileName
}

func (t Token) String() string {
	if t.Lexeme == "" {
		return "[" + t.Kind.String() + "]"
	}
	return "[" + t.Kind.String() + "]{" + t.Lexeme + "}"
}
