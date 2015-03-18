package main

import "fmt"
import "testing"

func walkNode(root *Node) {
	if root.children == nil {
		return
	}
	for child := range root.children {
		fmt.Printf("%s, %s, %v\n", child.tag, child.text, child.attributes)
		walkNode(child)
	}
}

func TestParser(t *testing.T) {
	test_case := []Token{
		Token{token_type: OpenTag, value: "<"},

		Token{token_type: Attribute, value: "h1"},

		Token{token_type: Attribute, value: "id"},
		Token{token_type: Assign, value: "="},
		Token{token_type: Quote, value: "\""},
		Token{token_type: Value, value: "one"},
		Token{token_type: Quote, value: "\""},

		Token{token_type: Attribute, value: "class"},
		Token{token_type: Assign, value: "="},
		Token{token_type: Quote, value: "\""},
		Token{token_type: Value, value: "one two three"},
		Token{token_type: Quote, value: "\""},

		Token{token_type: CloseTag, value: ">"},

		Token{token_type: Text, value: "Hello World"},

		Token{token_type: OpenTag, value: "<"},
		Token{token_type: ForwardSlash, value: "/"},
		Token{token_type: Attribute, value: "h1"},
		Token{token_type: CloseTag, value: ">"},
	}

	root := Parser(test_case)

	walkNode(root)
}
