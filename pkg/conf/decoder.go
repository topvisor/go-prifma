package conf

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
)

const DefaultSplitChars = "\r\n\t\f\v #\"';{}\\"

var DefaultDecoder = &Decoder{
	SplitChars:          DefaultSplitChars,
	TokenFactory:        DefaultTokenFactory,
	TokenHandlerFactory: NewTokenHandler,
}

type Decoder struct {
	SplitChars          string
	TokenFactory        TokenFactoryFunc
	TokenHandlerFactory TokenHandlerFactory
}

func (t *Decoder) Decode(base Block, globs ...string) error {
	for _, glob := range globs {
		filenames, err := filepath.Glob(glob)
		if err != nil {
			return err
		}

		for _, filename := range filenames {
			if err := t.decode(base, filename); err != nil {
				return err
			}
		}
	}
}

func (t *Decoder) decode(base Block, filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(t.split)

	for scanner.Scan() {
		tokens := t.TokenFactory(scanner.Text())
		if tokens == nil {
			continue
		}
		for _, token := range tokens {

		}
	}
}

func (t *Decoder) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, t.SplitChars); i >= 0 {
		return i + 1, data[0 : i+1], nil
	}
	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
