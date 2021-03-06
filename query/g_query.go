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
	OPEN_BRACKET
	CLOSE_BRACKET
	AND
	UNKNOWN

	ALL = "ALL"
)

type GQueryToken struct {
	token_type GQueryTokenType
	value      string
}

func (t GQueryToken) String() string {
	var token_string string

	switch t.token_type {
	case ID:
		token_string = "ID"
	case CLASS:
		token_string = "CLASS"
	case ATTRIBUTE:
		token_string = "ATTRIBUTE"
	case TAG:
		token_string = "TAG"
	case EQUALS:
		token_string = "EQUALS"
	case VALUE:
		token_string = "VALUE"
	case SPACE:
		token_string = "SPACE"
	case OPEN_BRACKET:
		token_string = "OPEN_BRACKET"
	case CLOSE_BRACKET:
		token_string = "CLOSE_BRACKET"
	case AND:
		token_string = "AND"
	default:
		token_string = "UNKNOWN"
	}

	return fmt.Sprintf("\"%s=%s\"", token_string, t.value)
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
		case '#', '.', '=', ' ', '[', ']', ',':
			finished = true
			tmp--
		}

		tmp += 1
	}

	return input[pos:tmp], tmp - 1

}

func tokenzier(input string) []GQueryToken {
	var value string
	var previous_token GQueryToken

	in_brackets := false

	tokens := []GQueryToken{}
	pos := 0

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
		case ',':
			tokens = append(tokens, GQueryToken{token_type: AND, value: ","})
		case ' ':
			tokens = append(tokens, GQueryToken{token_type: SPACE, value: " "})
		case '[':
			in_brackets = true
			tokens = append(tokens, GQueryToken{token_type: OPEN_BRACKET, value: "["})
		case ']':
			in_brackets = false
			tokens = append(tokens, GQueryToken{token_type: CLOSE_BRACKET, value: "]"})
		case '^':
			value, pos = consumeAttribute(input, pos+1)
			tokens = append(tokens, GQueryToken{token_type: TAG, value: value})
		default:
			value, pos = consumeAttribute(input, pos)
			if len(tokens) == 0 {
				tokens = append(tokens, GQueryToken{token_type: TAG, value: value})
			} else {
				previous_token = tokens[len(tokens)-1]

				if in_brackets && previous_token.token_type == EQUALS {
					tokens = append(tokens, GQueryToken{token_type: VALUE, value: value})
				} else if in_brackets {
					tokens = append(tokens, GQueryToken{token_type: ATTRIBUTE, value: value})
				} else if previous_token.token_type == SPACE || previous_token.token_type == AND {
					tokens = append(tokens, GQueryToken{token_type: TAG, value: value})
				} else {
					tokens = append(tokens, GQueryToken{token_type: VALUE, value: value})
				}
			}
		}

		pos += 1
	}

	return tokens
}

func GQuery(lookup string, node *parser.Node) *Query {
	var token GQueryToken
	var tokens []GQueryToken

	// In order to have the ability to chain queryies together
	// we need to keep track of two queries. One for the current query,
	// and the other for the returning query.
	// TODO(vishen): Possibly rethink the whole way of doing queries :s.
	query := NewQueryFromNode(node)
	return_query := &Query{}

	current_node := false

	if lookup == ALL {
		query.Nodes = node.AllChildren()
		return query
	}

	tokens = tokenzier(lookup)
	pos := 0

	for {

		if pos >= len(tokens) {
			break
		}

		token = tokens[pos]

		switch token.token_type {
		case ID:
			pos++
			token = tokens[pos]
			if current_node {
				query = query.FindByAttributeEquals("id", token.value)
			} else {
				query = query.FindChildrenByAttributeEquals("id", token.value)

			}
			current_node = true
		case CLASS:
			pos++
			token = tokens[pos]
			if current_node {
				query = query.FindByAttributeEquals("class", token.value)
			} else {
				query = query.FindChildrenByAttributeEquals("class", token.value)
			}
			current_node = true
		case TAG:
			query = query.FindByTag(token.value)
			current_node = true
		case SPACE, CLOSE_BRACKET:
			current_node = false
		// case OPEN_BRACKET:
		// 	current_node = true
		case AND:
			// If we get to the AND token, assume that that is the end of the previous
			// `Query` and add the query nodes to the return_query, then reset the `Query`
			// so that we can move onto the next query.
			return_query.AddNodes(query.Nodes)
			query.Reset()
			query.Nodes = append(query.Nodes, node)
		case ATTRIBUTE:
			attr_name := token.value
			if pos+2 < len(tokens) && tokens[pos+1].token_type == EQUALS && tokens[pos+2].token_type == VALUE {
				if current_node {
					query = query.FindByAttributeEquals(attr_name, tokens[pos+2].value)
				} else {
					query = query.FindChildrenByAttributeEquals(attr_name, tokens[pos+2].value)
				}
			} else {
				if current_node {
					query = query.FindByAttribute(attr_name)
				} else {
					query = query.FindChildrenByAttribute(attr_name)
				}
			}

			current_node = true
		}

		pos++
	}
	return_query.AddNodes(query.Nodes)

	return return_query

}
