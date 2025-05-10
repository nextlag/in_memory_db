package compute

import (
	"errors"
	"fmt"
	"strings"
)

const (
	startState = iota
	letterOrPunctuationState
	whiteSpaceState
)

func (p *Compute) Parse(query string) ([]string, error) {
	if len(query) == 0 {
		return nil, errors.New("empty query")
	}

	var (
		tokens []string
		sb     strings.Builder
		state  = startState
	)

	for i := 0; i < len(query); i++ {
		ch := query[i]

		switch state {
		case startState:
			if !isLetterOrPunctuationSymbol(ch) {
				return nil, fmt.Errorf("invalid symbol: '%c'", ch)
			}

			sb.WriteByte(ch)
			state = letterOrPunctuationState

		case letterOrPunctuationState:
			if isSpaceSymbol(ch) {
				tokens = append(tokens, sb.String())
				sb.Reset()
				state = whiteSpaceState

				continue
			}

			if !isLetterOrPunctuationSymbol(ch) {
				return nil, fmt.Errorf("invalid symbol: '%c'", ch)
			}

			sb.WriteByte(ch)

		case whiteSpaceState:
			if isSpaceSymbol(ch) {
				continue
			}

			if !isLetterOrPunctuationSymbol(ch) {
				return nil, fmt.Errorf("invalid symbol: '%c'", ch)
			}

			sb.WriteByte(ch)
			state = letterOrPunctuationState
		}
	}

	if state == letterOrPunctuationState {
		tokens = append(tokens, sb.String())
	}

	return tokens, nil
}

func isSpaceSymbol(ch byte) bool {
	return ch == '\t' || ch == '\n' || ch == ' '
}

func isLetterOrPunctuationSymbol(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '*' || ch == '/' ||
		ch == '_' || ch == '.'
}
