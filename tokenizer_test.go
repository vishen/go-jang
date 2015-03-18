package main

import "fmt"
import "testing"

func basicCheckArrayEquals(arr1 []Token, arr2 []Token) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i, token := range arr1 {
		if token.token_type != arr2[i].token_type || token.value != arr2[i].value {
			return false
		}
	}

	return true
}

func TestTokenizer(t *testing.T) {
	cases := []struct {
		in     string
		wanted []Token
	}{
		{`<h1 id="one" class="one two three"> Hello World </h1>`,
			[]Token{
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
			},
		},
	}

	for _, c := range cases {
		got := Tokenizer(c.in)
		//fmt.Printf("%v\n", c.wanted)
		fmt.Printf("%v\n", got)
		if !basicCheckArrayEquals(got, c.wanted) {
			t.Error("Failed")
		}
	}
}
