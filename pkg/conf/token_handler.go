package conf

const IncludeDirective = "include"

type TokenHandlerFactory func(base Block, decoder *Decoder) TokenHandler

type TokenHandler interface {
	Handle(token Token) error
}

type DefaultTokenHandler struct {
	BlockWrapper *BlockWrapper
	Decoder      *Decoder

	Directive           string
	Args                []string
	LastArg             string
	Line                string
	LineNumber          int
	IsCommentLine       bool
	IsDoubleQuotaOpened bool
	IsSingleQuotaOpened bool
	IsEscaped           bool
	IsEscapedNow        bool
}

func NewTokenHandler(base Block, decoder *Decoder) TokenHandler {
	t := &DefaultTokenHandler{
		BlockWrapper: &BlockWrapper{
			Current: base,
		},
		Decoder:    decoder,
		Args:       make([]string, 0),
		LineNumber: 1,
	}

	return t
}

func (t *DefaultTokenHandler) Handle(token Token) (err error) {
	t.IsEscapedNow = false

	if nlToken, ok := token.(*NewlineToken); ok {
		return t.HandleNewlineToken(nlToken)
	}

	t.Line += token.String()

	if t.IsCommentLine {
		return nil
	}

	switch tokenImpl := token.(type) {
	case *StringToken:
		err = t.HandleStringToken(tokenImpl)
	case *WhitespaceToken:
		err = t.HandleWhitespaceToken(tokenImpl)
	case *DoubleQuotaToken:
		err = t.HandleDoubleQuotaToken(tokenImpl)
	case *SingleQuotaToken:
		err = t.HandleSingleQuotaToken(tokenImpl)
	case *HashToken:
		err = t.HandleHashToken(tokenImpl)
	case *SemicolonToken:
		err = t.HandleSemicolonToken(tokenImpl)
	case *OpeningCurlyBracketToken:
		err = t.HandleOpeningCurlyBracketToken(tokenImpl)
	case *ClosingCurlyBracketToken:
		err = t.HandleClosingCurlyBracketToken(tokenImpl)
	case *BackslashToken:
		err = t.HandleBackslashToken(tokenImpl)
	}

	if !t.IsEscapedNow {
		t.IsEscaped = false
	}

	return err
}

func (t *DefaultTokenHandler) HandleStringToken(token *StringToken) error {
	t.LastArg += token.String()

	return nil
}

func (t *DefaultTokenHandler) HandleWhitespaceToken(token *WhitespaceToken) error {
	if t.IsStringEscaped() {
		t.LastArg += token.String()

		return nil
	}

	t.CommitLastArg()

	return nil
}

func (t *DefaultTokenHandler) HandleNewlineToken(token *NewlineToken) error {
	if !t.IsCallCommitted() || t.IsStringEscaped() {
		return NewErrParse(t.LineNumber, t.Line)
	}

	t.IsCommentLine = false
	t.Line = ""
	t.LineNumber++

	return nil
}

func (t *DefaultTokenHandler) HandleDoubleQuotaToken(token *DoubleQuotaToken) error {
	if t.Directive == "" {
		return NewErrParse(t.LineNumber, t.Line)
	}

	if t.IsSingleQuotaOpened || t.IsEscaped {
		t.LastArg += token.String()

		return nil
	}

	t.IsDoubleQuotaOpened = !t.IsDoubleQuotaOpened

	return nil
}

func (t *DefaultTokenHandler) HandleSingleQuotaToken(token *SingleQuotaToken) error {
	if t.Directive == "" {
		return NewErrParse(t.LineNumber, t.Line)
	}

	if t.IsDoubleQuotaOpened || t.IsEscaped {
		t.LastArg += token.String()

		return nil
	}

	t.IsSingleQuotaOpened = !t.IsSingleQuotaOpened

	return nil
}

func (t *DefaultTokenHandler) HandleHashToken(token *HashToken) error {
	if t.IsStringEscaped() {
		t.LastArg += token.String()

		return nil
	}

	if !t.IsCallCommitted() {
		return NewErrParse(t.LineNumber, t.Line)
	}

	t.IsCommentLine = true

	return nil
}

func (t *DefaultTokenHandler) HandleSemicolonToken(token *SemicolonToken) (err error) {
	if t.IsStringEscaped() {
		t.LastArg += token.String()

		return nil
	}

	if t.IsCallCommitted() {
		return NewErrParse(t.LineNumber, t.Line)
	}
	t.CommitLastArg()

	if t.Directive == IncludeDirective {
		err = t.Decoder.Decode(t.BlockWrapper.Current, t.Args...)
	} else {
		err = t.BlockWrapper.Current.Call(NewCommand(t.Directive, t.Args...))
	}

	t.Directive = ""
	t.Args = make([]string, 0)

	return err
}

func (t *DefaultTokenHandler) HandleOpeningCurlyBracketToken(token *OpeningCurlyBracketToken) error {
	if t.IsStringEscaped() {
		t.LastArg += token.String()

		return nil
	}

	if t.IsCallCommitted() {
		return NewErrParse(t.LineNumber, t.Line)
	}
	t.CommitLastArg()

	block, err := t.BlockWrapper.Current.CallBlock(NewCommand(t.Directive, t.Args...))

	t.BlockWrapper = &BlockWrapper{
		Parent:  t.BlockWrapper,
		Current: block,
	}
	t.Directive = ""
	t.Args = make([]string, 0)

	return err
}

func (t *DefaultTokenHandler) HandleClosingCurlyBracketToken(token *ClosingCurlyBracketToken) error {
	if t.IsStringEscaped() {
		t.LastArg += token.String()

		return nil
	}

	if !t.IsCallCommitted() || t.BlockWrapper.Parent == nil {
		return NewErrParse(t.LineNumber, t.Line)
	}

	t.BlockWrapper = t.BlockWrapper.Parent

	return nil
}

func (t *DefaultTokenHandler) HandleBackslashToken(token *BackslashToken) error {
	if t.Directive == "" {
		return NewErrParse(t.LineNumber, t.Line)
	}

	if t.IsEscaped {
		t.LastArg += token.String()
	} else {
		t.IsEscaped = true
		t.IsEscapedNow = true
	}

	return nil
}

func (t *DefaultTokenHandler) IsCallCommitted() bool {
	return t.Directive == "" && t.LastArg == ""
}

func (t *DefaultTokenHandler) IsStringEscaped() bool {
	return t.IsDoubleQuotaOpened || t.IsSingleQuotaOpened || t.IsEscaped
}

func (t *DefaultTokenHandler) CommitLastArg() {
	if t.LastArg != "" {
		if t.Directive == "" {
			t.Directive = t.LastArg
		} else {
			t.Args = append(t.Args, t.LastArg)
		}

		t.LastArg = ""
	}
}
