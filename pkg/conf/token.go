package conf

type Token interface {
	String() string
}

type StringToken struct {
	data string
}

func (t *StringToken) String() string {
	return t.data
}

type WhitespaceToken struct {
	data byte
}

func (t *WhitespaceToken) String() string {
	return string([]byte{t.data})
}

type NewlineToken struct {
}

func (t *NewlineToken) String() string {
	return "\n"
}

type DoubleQuotaToken struct {
}

func (t *DoubleQuotaToken) String() string {
	return "\""
}

type SingleQuotaToken struct {
}

func (t *SingleQuotaToken) String() string {
	return "'"
}

type HashToken struct {
}

func (t *HashToken) String() string {
	return "#"
}

type SemicolonToken struct {
}

func (t *SemicolonToken) String() string {
	return ";"
}

type OpeningCurlyBracketToken struct {
}

func (t *OpeningCurlyBracketToken) String() string {
	return "{"
}

type ClosingCurlyBracketToken struct {
}

func (t *ClosingCurlyBracketToken) String() string {
	return "}"
}

type BackslashToken struct {
}

func (t *BackslashToken) String() string {
	return "\\"
}
