package conf

type TokenFactoryFunc func(data string) []Token

func DefaultTokenFactory(data string) []Token {
	dataLen := len(data)
	if dataLen == 0 {
		return nil
	}

	var token Token
	if dataLen > 1 {
		token = &StringToken{data: data[:dataLen-1]}
	}

	var lastCharToken Token
	lastChar := data[dataLen-1]
	switch lastChar {
	case '\r':
		fallthrough
	case '\t':
		fallthrough
	case '\f':
		fallthrough
	case '\v':
		fallthrough
	case ' ':
		lastCharToken = &WhitespaceToken{data: lastChar}
	case '\n':
		lastCharToken = new(NewlineToken)
	case '#':
		lastCharToken = new(HashToken)
	case '"':
		lastCharToken = new(DoubleQuotaToken)
	case '\'':
		lastCharToken = new(SingleQuotaToken)
	case ';':
		lastCharToken = new(SemicolonToken)
	case '{':
		lastCharToken = new(OpeningCurlyBracketToken)
	case '}':
		lastCharToken = new(ClosingCurlyBracketToken)
	case '\\':
		lastCharToken = new(BackslashToken)
	default:
		token = &StringToken{data: data}
	}

	tokens := make([]Token, 0, 2)

	if token != nil {
		tokens = append(tokens, token)
	}
	if lastCharToken != nil {
		tokens = append(tokens, lastCharToken)
	}

	return tokens
}
