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

const (
	OpenTag TokenType = iota
	CloseTag
	Attribute
	Value
	Text
	Assign
	ForwardSlash
	Quote
)

func consumeValue(input string, pos int) (string, int) {

	current := pos
	for {
		if current >= len(input) {
			break
		}

		switch input[current] {
		case '"', '\'':
			return strings.TrimSpace(input[pos:current]), current - 1
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
		case '<', '>', '/':
			return strings.TrimSpace(input[pos:current]), current - 1
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
	in_quotes := false

	for {

		if pos >= length {
			break
		}
		switch input[pos] {
		case '\n':
			line += 1
			column = 0
		case '<':
			tokens = append(tokens, Token{token_type: OpenTag, value: "<", column: column, line: line})
		case '>':
			tokens = append(tokens, Token{token_type: CloseTag, value: ">", column: column, line: line})
		case '/':
			tokens = append(tokens, Token{token_type: ForwardSlash, value: "/", column: column, line: line})
		case '=':
			tokens = append(tokens, Token{token_type: Assign, value: "=", column: column, line: line})
		case '"':
			in_quotes = !in_quotes
			tokens = append(tokens, Token{token_type: Quote, value: "\"", column: column, line: line})
		case '\'':
			in_quotes = !in_quotes
			tokens = append(tokens, Token{token_type: Quote, value: "'", column: column, line: line})
		case ' ':
			pos += 1
			column += 1
			continue
		default:
			last_token_type = tokens[len(tokens)-1].token_type
			if last_token_type == CloseTag {
				value, count := consumeText(input, pos)
				pos = count
				tokens = append(tokens, Token{token_type: Text, value: value, column: column, line: line})
			} else if last_token_type == Quote && in_quotes {
				value, count := consumeValue(input, pos)
				pos = count
				tokens = append(tokens, Token{token_type: Value, value: value, column: column, line: line})
			} else {
				value, count := consumeAttribute(input, pos)
				pos = count
				tokens = append(tokens, Token{token_type: Attribute, value: value, column: column, line: line})
			}
		}
		column += 1
		pos += 1

	}

	return tokens
}

func main() {
	fmt.Println("Hello, Mac!")
}
