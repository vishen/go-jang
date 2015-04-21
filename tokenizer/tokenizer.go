package tokenizer

import (
	"fmt"
	"strings"
)

type TokenType int

type Token struct {
	Token_type TokenType
	Value      string
	Column     int
	Line       int
}

func (t Token) String() string {
	var token_type_string string
	switch t.Token_type {
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
	case Comment:
		token_type_string = "Comment"
	default:
		token_type_string = "Unknown"
	}
	return fmt.Sprintf("[%s] %s %d:%d", token_type_string, t.Value, t.Line, t.Column)
}

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

	var end_char uint8 = ' '
	in_quotes := false
	if input[pos] == '"' || input[pos] == '\'' {
		end_char = input[pos]
		pos++

		in_quotes = true
	}

	current := pos
	for {
		if current >= len(input) {
			break
		}
		if !in_quotes {
			switch input[current] {
			case '>', '/', '=', ' ':
				return strings.TrimSpace(input[pos:current]), current - 1
			}
		} else {
			switch input[current] {
			case end_char:
				return strings.TrimSpace(input[pos:current]), current
			}
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
			// Ignore comments for now
			// TODO(vishen): Make some more tokens for comments - probably easier if
			// this is done in the tokenizer and just add a new Comment token.
			if input[pos+1] == '!' && input[pos+2] == '-' {
				tmp := pos + 3
				start := tmp + 1
				for {
					if input[tmp] == '>' && input[tmp-1] == '-' {
						pos = tmp
						tokens = append(tokens, Token{Token_type: Comment, Value: strings.TrimSpace(input[start : tmp-2]), Column: column, Line: line})
						break
					}

					tmp += 1
				}
			} else {
				in_declaration = true
				tokens = append(tokens, Token{Token_type: OpenTag, Value: "<", Column: column, Line: line})
			}
		case '>':
			in_declaration = false
			// Sometimes an Attribute isn't followed by an Assign and Value.
			// In this case just pretend that there is blank ones for the parser
			// to pick up the Attribute. Not a huge fan of this.
			// TODO(vishen): Find a better way of doing this.
			if last_token_type == Attribute {
				tokens = append(tokens,
					Token{Token_type: Assign, Value: "=", Column: column, Line: line},
					Token{Token_type: Value, Value: "", Column: column, Line: line},
				)
			}
			tokens = append(tokens, Token{Token_type: CloseTag, Value: ">", Column: column, Line: line})
		case '/':
			// Because of javascript comments we just assume if we see a '/'
			// and we are not in the declaration just comsime the Text

			//TODO(vishen): Find a better way to handle this
			if !in_declaration {
				value, count := consumeText(input, pos)
				pos = count
				tokens = append(tokens, Token{Token_type: Text, Value: value, Column: column, Line: line})
			} else {
				tokens = append(tokens, Token{Token_type: ForwardSlash, Value: "/", Column: column, Line: line})
			}
		case '=':
			tokens = append(tokens, Token{Token_type: Assign, Value: "=", Column: column, Line: line})
		case ' ', '\t', '\r':
			// Leave empty
		default:
			if len(tokens) > 0 {
				last_token = tokens[len(tokens)-1]
				last_token_type = last_token.Token_type
				if last_token_type == CloseTag {
					value, count := consumeText(input, pos)
					pos = count
					tokens = append(tokens, Token{Token_type: Text, Value: value, Column: column, Line: line})
				} else if last_token_type == Assign {
					value, count := consumeValue(input, pos)
					pos = count
					tokens = append(tokens, Token{Token_type: Value, Value: value, Column: column, Line: line})
				} else if in_declaration {

					// Sometimes an Attribute isn't followed by an Assign and Value.
					// In this case just pretend that there is blank ones for the parser
					// to pick up the Attribute. Not a huge fan of this.
					// TODO(vishen): Find a better way of doing this.
					if last_token_type == Attribute {
						tokens = append(tokens,
							Token{Token_type: Assign, Value: "=", Column: column, Line: line},
							Token{Token_type: Value, Value: "", Column: column, Line: line},
						)
					}

					value, count := consumeAttribute(input, pos)
					pos = count

					if last_token_type == OpenTag || last_token_type == ForwardSlash {
						tmp_token_type = Tag
					} else {
						tmp_token_type = Attribute
					}

					tokens = append(tokens, Token{Token_type: tmp_token_type, Value: value, Column: column, Line: line})
				}
			}
		}
		column += 1
		pos += 1

	}

	return tokens
}
