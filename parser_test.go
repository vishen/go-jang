package main

import "fmt"
import "testing"

func createNode(nt TokenType, value string) Token {
	return Token{token_type: nt, value: value}
}

func TestParserSingleNode(t *testing.T) {
	test_case := []Token{
		createNode(OpenTag, "<"),
		createNode(Tag, "html"),
		createNode(CloseTag, ">"),

		createNode(Text, "Hello World"),

		createNode(OpenTag, "<"),
		createNode(ForwardSlash, "/"),
		createNode(Tag, "html"),
		createNode(CloseTag, ">"),
	}

	wanted := Node{
		tag:  "html",
		text: "Hello World",
	}

	actual := Parser(test_case)

	if wanted.tag != actual.tag || wanted.text != actual.text {
		fmt.Printf("Expected - %s, actual - %s\n", wanted, actual)
		t.Error("Failed")
	}
}

func TestParserChildrenNode(t *testing.T) {
	test_case := []Token{
		createNode(OpenTag, "<"),
		createNode(Tag, "html"),
		createNode(CloseTag, ">"),

		createNode(OpenTag, "<"),
		createNode(Tag, "h1"),
		createNode(CloseTag, ">"),

		createNode(Text, "Hello World"),

		createNode(OpenTag, "<"),
		createNode(ForwardSlash, "/"),
		createNode(Tag, "h1"),
		createNode(CloseTag, ">"),

		createNode(OpenTag, "<"),
		createNode(ForwardSlash, "/"),
		createNode(Tag, "html"),
		createNode(CloseTag, ">"),
	}

	wanted := Node{
		tag: "html",

		children: []*Node{
			&Node{
				tag:  "h1",
				text: "Hello World",
			},
		},
	}

	actual := Parser(test_case)

	if wanted.tag != actual.tag || wanted.children[0].tag != actual.children[0].tag || wanted.children[0].text != actual.children[0].text {
		fmt.Printf("Expected  - %s, actual - %s\n", wanted, actual)
		t.Error("Failed")
	}
}
