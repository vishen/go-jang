package query

import (
	"fmt"
	"github.com/vishen/go-jang/parser"
	"github.com/vishen/go-jang/request"
)

type GQueryTokenType int

const (
	ID GQueryTokenType = iota
	CLASS
	ATTRIBUTE
	TAG
	EQUALS
	VALUE
	SPACE
	UNKNOWN
)

type GQueryToken struct {
	token_type GQueryTokenType
	value      string
}

func GetGQueryFromUrl(url, lookup string) (*Query, error) {
	root, err := request.Get(url)

	if err != nil {
		return nil, err
	}

	return GQuery(lookup, root), nil
}

func consumeAttribute(input string, pos int) (string, int) {
	tmp := pos
	finished := false

	for {
		if finished || tmp >= len(input) {
			break
		}
		switch input[tmp] {
		case '#', '.', '=', ' ':
			finished = true
			tmp--
		}

		tmp += 1
	}

	return input[pos:tmp], tmp - 1

}

func tokenzier(input string) []GQueryToken {
	tokens := []GQueryToken{}
	pos := 0
	var value string

	for {

		if pos >= len(input) {
			break
		}

		switch input[pos] {
		case '#':
			tokens = append(tokens, GQueryToken{token_type: ID, value: "#"})
		case '.':
			tokens = append(tokens, GQueryToken{token_type: CLASS, value: "."})
		case '=':
			tokens = append(tokens, GQueryToken{token_type: EQUALS, value: "="})
		case ' ':
			tokens = append(tokens, GQueryToken{token_type: SPACE, value: " "})
		case '^':
			value, pos = consumeAttribute(input, pos+1)
			tokens = append(tokens, GQueryToken{token_type: TAG, value: value})
		default:
			value, pos = consumeAttribute(input, pos)

			if len(tokens) > 0 && tokens[len(tokens)-1].token_type == SPACE {
				tokens = append(tokens, GQueryToken{token_type: ATTRIBUTE, value: value})
			} else {
				tokens = append(tokens, GQueryToken{token_type: VALUE, value: value})
			}
		}

		pos += 1
	}

	return tokens
}

func GQuery(lookup string, node *parser.Node) *Query {

	// Handles the `#`, `.`, `<attr>` and `<attr>=<valiue>`

	tokens := tokenzier(lookup)

	var token GQueryToken

	query := NewQueryFromNode(node)

	pos := 0

	for {

		if pos >= len(tokens) {
			break
		}

		token = tokens[pos]

		fmt.Println(tokens)

		switch token.token_type {
		case ID:
			pos++
			token = tokens[pos]
			query = query.FindByAttributeEquals("id", token.value)
		case CLASS:
			pos++
			token = tokens[pos]
			query = query.FindByAttributeEquals("class", token.value)
		case TAG:
			query = query.FindByTag(token.value)
		case ATTRIBUTE:
			attr_name := token.value
			if pos+2 < len(tokens) && tokens[pos+1].token_type == EQUALS && tokens[pos+2].token_type == VALUE {
				query = query.FindByAttributeEquals(attr_name, tokens[pos+2].value)
			} else {
				query = query.FindByAttribute(attr_name)
			}
		}

		pos++
	}

	return query

}
