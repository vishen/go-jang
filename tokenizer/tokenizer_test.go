package tokenizer

import "fmt"
import "testing"

func basicCheckArrayEquals(arr1 []Token, arr2 []Token) bool {
	if len(arr1) != len(arr2) {
		fmt.Printf("[Error] Mismatched number of tokens; %d & %d", len(arr1), len(arr2))
		return false
	}

	for i, token := range arr1 {
		if token.Token_type != arr2[i].Token_type || token.Value != arr2[i].Value {
			fmt.Printf("[Error] %d %s != %s\n", i, token, arr2[i])
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
		{`<h1 id=one class="one two three"> Hello World </h1>`,
			[]Token{
				Token{Token_type: OpenTag, Value: "<"},

				Token{Token_type: Tag, Value: "h1"},

				Token{Token_type: Attribute, Value: "id"},
				Token{Token_type: Assign, Value: "="},
				Token{Token_type: Value, Value: "one"},

				Token{Token_type: Attribute, Value: "class"},
				Token{Token_type: Assign, Value: "="},
				Token{Token_type: Value, Value: "one two three"},

				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: Text, Value: "Hello World"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: ForwardSlash, Value: "/"},
				Token{Token_type: Tag, Value: "h1"},
				Token{Token_type: CloseTag, Value: ">"},
			},
		},
		{
			`<h1 id=one class="one two three"> 
Hello World
</h1>`,
			[]Token{
				Token{Token_type: OpenTag, Value: "<"},

				Token{Token_type: Tag, Value: "h1"},

				Token{Token_type: Attribute, Value: "id"},
				Token{Token_type: Assign, Value: "="},
				Token{Token_type: Value, Value: "one"},

				Token{Token_type: Attribute, Value: "class"},
				Token{Token_type: Assign, Value: "="},
				Token{Token_type: Value, Value: "one two three"},

				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: Text, Value: "Hello World"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: ForwardSlash, Value: "/"},
				Token{Token_type: Tag, Value: "h1"},
				Token{Token_type: CloseTag, Value: ">"},
			},
		},
		{
			`<!doctype html>
<html>
	<head>
		<meta content="hello">
		<meta content="one"/>
	</head>
</html>`,
			[]Token{
				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: Tag, Value: "!doctype"},
				Token{Token_type: Attribute, Value: "html"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: Tag, Value: "html"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: Tag, Value: "head"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: Tag, Value: "meta"},
				Token{Token_type: Attribute, Value: "content"},
				Token{Token_type: Assign, Value: "="},
				Token{Token_type: Value, Value: "hello"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: Tag, Value: "meta"},
				Token{Token_type: Attribute, Value: "content"},
				Token{Token_type: Assign, Value: "="},
				Token{Token_type: Value, Value: "one"},
				Token{Token_type: ForwardSlash, Value: "/"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: ForwardSlash, Value: "/"},
				Token{Token_type: Tag, Value: "head"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: ForwardSlash, Value: "/"},
				Token{Token_type: Tag, Value: "html"},
				Token{Token_type: CloseTag, Value: ">"},
			},
		},
		{
			`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"> 
		<html xmlns="http://www.w3.org/1999/xhtml">Hello World`,
			[]Token{
				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: Tag, Value: "!DOCTYPE"},
				Token{Token_type: Attribute, Value: "html"},
				Token{Token_type: Attribute, Value: "PUBLIC"},
				Token{Token_type: Attribute, Value: "-//W3C//DTD XHTML 1.0 Transitional//EN"},
				Token{Token_type: Attribute, Value: "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: Tag, Value: "html"},
				Token{Token_type: Attribute, Value: "xmlns"},
				Token{Token_type: Assign, Value: "="},
				Token{Token_type: Value, Value: "http://www.w3.org/1999/xhtml"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: Text, Value: "Hello World"},
			},
		},
		{
			`<html><!-- This is a comment --></html>`,
			[]Token{
				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: Tag, Value: "html"},
				Token{Token_type: CloseTag, Value: ">"},

				Token{Token_type: OpenTag, Value: "<"},
				Token{Token_type: ForwardSlash, Value: "/"},
				Token{Token_type: Tag, Value: "html"},
				Token{Token_type: CloseTag, Value: ">"},
			},
		},
	}

	for _, c := range cases {
		got := Tokenizer(c.in)

		if !basicCheckArrayEquals(got, c.wanted) {
			t.Error("Failed")
		}
	}

}
