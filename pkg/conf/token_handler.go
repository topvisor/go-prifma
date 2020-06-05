package conf

import (
	"github.com/topvisor/go-prifma/pkg/utils"
)

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
	IsBackslashed       bool
	IsBackslashedNow    bool
}

func NewTokenHandler(base Block, decoder *Decoder) *DefaultTokenHandler {
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
	t.IsBackslashedNow = false

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

	if !t.IsBackslashedNow {
		t.IsBackslashed = false
	}

	return err
}

func (t *DefaultTokenHandler) HandleTokenAsString(token Token) error {
	return t.HandleStringToken(&StringToken{data: token.String()})
}

func (t *DefaultTokenHandler) HandleStringToken(token *StringToken) error {
	if t.IsBackslashed {
		t.IsBackslashed = false

		t.LastArg += new(BackslashToken).String()
	}

	t.LastArg += token.String()

	return nil
}

func (t *DefaultTokenHandler) HandleWhitespaceToken(token *WhitespaceToken) error {
	if t.IsQuotaOpened() {
		return t.HandleTokenAsString(token)
	}

	if t.IsBackslashed {
		t.IsBackslashed = false

		return t.HandleTokenAsString(token)
	}

	return t.CommitLastArg(false)
}

func (t *DefaultTokenHandler) HandleNewlineToken(token *NewlineToken) (err error) {
	if t.IsQuotaOpened() {
		err = t.HandleTokenAsString(token)
	} else if t.IsBackslashed {
		t.IsBackslashed = false

		err = t.HandleTokenAsString(token)
	} else if !t.IsCallCommitted() {
		err = NewErrParse(t.LineNumber, t.Line, "unexpected new line")
	} else {
		t.IsCommentLine = false
	}

	t.Line = ""
	t.LineNumber++

	return err
}

func (t *DefaultTokenHandler) HandleDoubleQuotaToken(token *DoubleQuotaToken) error {
	if t.Directive == "" {
		return NewErrParse(t.LineNumber, t.Line, "unexpected double quota")
	}

	if t.IsSingleQuotaOpened {
		return t.HandleTokenAsString(token)
	}

	if t.IsBackslashed {
		t.IsBackslashed = false

		return t.HandleTokenAsString(token)
	}

	t.IsDoubleQuotaOpened = !t.IsDoubleQuotaOpened

	if !t.IsDoubleQuotaOpened {
		return t.CommitLastArg(true)
	}

	return nil
}

func (t *DefaultTokenHandler) HandleSingleQuotaToken(token *SingleQuotaToken) error {
	if t.Directive == "" {
		return NewErrParse(t.LineNumber, t.Line, "unexpected single quota")
	}

	if t.IsDoubleQuotaOpened {
		return t.HandleTokenAsString(token)
	}

	if t.IsBackslashed {
		t.IsBackslashed = false

		return t.HandleTokenAsString(token)
	}

	t.IsSingleQuotaOpened = !t.IsSingleQuotaOpened

	if !t.IsSingleQuotaOpened {
		return t.CommitLastArg(true)
	}

	return nil
}

func (t *DefaultTokenHandler) HandleHashToken(token *HashToken) error {
	if t.IsQuotaOpened() {
		return t.HandleTokenAsString(token)
	}

	if t.IsBackslashed {
		t.IsBackslashed = false

		return t.HandleTokenAsString(token)
	}

	if !t.IsCallCommitted() {
		return NewErrParse(t.LineNumber, t.Line, "unexpected #")
	}

	t.IsCommentLine = true

	return nil
}

func (t *DefaultTokenHandler) HandleSemicolonToken(token *SemicolonToken) (err error) {
	if t.IsQuotaOpened() {
		return t.HandleTokenAsString(token)
	}

	if t.IsBackslashed {
		t.IsBackslashed = false

		return t.HandleTokenAsString(token)
	}

	if t.IsCallCommitted() {
		return NewErrParse(t.LineNumber, t.Line, "unexpected semicolon")
	}

	if err = t.CommitLastArg(false); err != nil {
		return err
	}

	if t.Directive == IncludeDirective {
		err = t.Decoder.Decode(t.BlockWrapper.Current, t.Args...)
	} else {
		err = t.BlockWrapper.Current.Call(NewCommand(t.LineNumber, t.Directive, t.Args...))
	}

	t.Directive = ""
	t.Args = make([]string, 0)

	return err
}

func (t *DefaultTokenHandler) HandleOpeningCurlyBracketToken(token *OpeningCurlyBracketToken) (err error) {
	if t.IsQuotaOpened() {
		return t.HandleTokenAsString(token)
	}

	if t.IsBackslashed {
		t.IsBackslashed = false

		return t.HandleTokenAsString(token)
	}

	if t.IsCallCommitted() {
		return NewErrParse(t.LineNumber, t.Line, "unexpected opening curly bracket")
	}

	if err = t.CommitLastArg(false); err != nil {
		return err
	}

	block, err := t.BlockWrapper.Current.CallBlock(NewCommand(t.LineNumber, t.Directive, t.Args...))

	t.BlockWrapper = &BlockWrapper{
		Parent:  t.BlockWrapper,
		Current: block,
	}

	t.Directive = ""
	t.Args = make([]string, 0)

	return err
}

func (t *DefaultTokenHandler) HandleClosingCurlyBracketToken(token *ClosingCurlyBracketToken) error {
	if t.IsQuotaOpened() {
		return t.HandleTokenAsString(token)
	}

	if t.IsBackslashed {
		t.IsBackslashed = false

		return t.HandleTokenAsString(token)
	}

	if !t.IsCallCommitted() || t.BlockWrapper.Parent == nil {
		return NewErrParse(t.LineNumber, t.Line, "unexpected closing curly bracket")
	}

	t.BlockWrapper = t.BlockWrapper.Parent

	return nil
}

func (t *DefaultTokenHandler) HandleBackslashToken(token *BackslashToken) error {
	if t.Directive == "" {
		return NewErrParse(t.LineNumber, t.Line, "unexpected backslash")
	}

	if t.IsBackslashed {
		t.IsBackslashed = false

		return t.HandleTokenAsString(token)
	}

	t.IsBackslashed = true
	t.IsBackslashedNow = true

	return nil
}

func (t *DefaultTokenHandler) IsQuotaOpened() bool {
	return t.IsDoubleQuotaOpened || t.IsSingleQuotaOpened
}

func (t *DefaultTokenHandler) IsCallCommitted() bool {
	return t.Directive == "" && t.LastArg == ""
}

func (t *DefaultTokenHandler) CommitLastArg(commitEmptyArg bool) (err error) {
	if t.Directive == "" {
		t.Directive = t.LastArg
	} else if t.LastArg != "" || commitEmptyArg {
		t.Args = append(t.Args, utils.UnescapeString(t.LastArg))
	}

	t.LastArg = ""

	return nil
}
