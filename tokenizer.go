package main

import (
	"fmt"
	"strings"
)

type TokenType int

type Token struct {
	token_type TokenType
	value      string
	column     int
	line       int
}

func (t Token) String() string {
	var token_type_string string
	switch t.token_type {
	case OpenTag:
		token_type_string = "OpenTag"
	case CloseTag:
		token_type_string = "CloseTag"
	case Attribute:
		token_type_string = "Attribute"
	case Value:
		token_type_string = "Value"
	case Text:
		token_type_string = "Text"
	case Assign:
		token_type_string = "Assign"
	case ForwardSlash:
		token_type_string = "ForwardSlash"
	case Tag:
		token_type_string = "Tag"
	default:
		token_type_string = "Unknown"
	}
	return fmt.Sprintf("[%s] %s %d:%d", token_type_string, t.value, t.line, t.column)
}

const (
	OpenTag TokenType = iota
	CloseTag
	Attribute
	Value
	Text
	Assign
	ForwardSlash
	Tag
)

func consumeValue(input string, pos int) (string, int) {
	// Start consuming an attributes value
	// Starts at possibly the following
	// eg: href=^"/hello" OR href=^hello

	var end_char uint8 = ' '
	if input[pos] == '"' || input[pos] == '\'' {
		end_char = input[pos]
		pos++
	}

	current := pos

	for {
		if current >= len(input) {
			break
		}

		switch input[current] {
		case end_char:
			return strings.TrimSpace(input[pos:current]), current
		case '>':
			if end_char == ' ' {
				return strings.TrimSpace(input[pos:current]), current - 1
			}
		}
		current += 1
	}
	return strings.TrimSpace(input[pos:current]), current
}

func consumeAttribute(input string, pos int) (string, int) {

	current := pos
	for {
		if current >= len(input) {
			break
		}

		switch input[current] {
		case '>', '/', '=', ' ':
			return strings.TrimSpace(input[pos:current]), current - 1
		}
		current += 1
	}
	return strings.TrimSpace(input[pos:current]), current
}

func consumeText(input string, pos int) (string, int) {
	current := pos
	for {
		if current >= len(input) {
			break
		}

		switch input[current] {
		case '<':
			if input[current+1] == '/' {
				return strings.TrimSpace(input[pos:current]), current - 1
			}
		}
		current += 1
	}
	return strings.TrimSpace(input[pos:current]), current
}

func Tokenizer(input string) []Token {
	length := len(input)
	pos := 0
	line := 0
	column := 0
	tokens := []Token{}

	var last_token_type TokenType
	var last_token Token
	var tmp_token_type TokenType
	in_declaration := false

	for {

		if pos >= length {
			break
		}
		switch input[pos] {
		case '\n':
			line += 1
			column = 0
		case '<':
			in_declaration = true
			tokens = append(tokens, Token{token_type: OpenTag, value: "<", column: column, line: line})
		case '>':
			in_declaration = false
			tokens = append(tokens, Token{token_type: CloseTag, value: ">", column: column, line: line})
		case '/':
			tokens = append(tokens, Token{token_type: ForwardSlash, value: "/", column: column, line: line})
		case '=':
			tokens = append(tokens, Token{token_type: Assign, value: "=", column: column, line: line})
		case ' ', '\t':
		// pos += 1
		// column += 1
		// continue
		default:
			last_token = tokens[len(tokens)-1]
			last_token_type = last_token.token_type
			if last_token_type == CloseTag {
				value, count := consumeText(input, pos)
				pos = count
				tokens = append(tokens, Token{token_type: Text, value: value, column: column, line: line})
			} else if last_token_type == Assign {
				value, count := consumeValue(input, pos)
				pos = count
				tokens = append(tokens, Token{token_type: Value, value: value, column: column, line: line})
			} else if in_declaration {
				value, count := consumeAttribute(input, pos)
				pos = count

				if last_token_type == OpenTag || last_token_type == ForwardSlash {
					tmp_token_type = Tag
				} else {
					tmp_token_type = Attribute
				}

				tokens = append(tokens, Token{token_type: tmp_token_type, value: value, column: column, line: line})
			}
		}
		column += 1
		pos += 1

	}

	return tokens
}
