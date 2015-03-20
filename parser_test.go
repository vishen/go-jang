package main

//import "fmt"
import "testing"

/*
New Node: div, self: 0x2082543c0, parent: 0x0
New Node: div2, self: 0x208254420, parent: 0x2082543c0
New Node: h3, self: 0x208254480, parent: 0x208254420
New Node: span, self: 0x208254540, parent: 0x2082544e0
New Node: h2, self: 0x2082545a0, parent: 0x208254540
End Tag:  h2
End Tag:  span
End Tag:  h3
End Tag:  div2
New Node: h1, self: 0x2082544e0, parent: 0x208254480
End Tag:  h1
End Tag:  div
*/

func TestParser(t *testing.T) {
	test_case := []Token{
		Token{token_type: OpenTag, value: "<"},
		Token{token_type: Attribute, value: "div"},
		Token{token_type: CloseTag, value: ">"},
		
		Token{token_type: OpenTag, value: "<"},
		Token{token_type: Attribute, value: "div2"},
		Token{token_type: CloseTag, value: ">"},
		
		
		Token{token_type: OpenTag, value: "<"},
		Token{token_type: Attribute, value: "h3"},
		Token{token_type: CloseTag, value: ">"},
		Token{token_type: Text, value: "Parse Bofucker"},
		Token{token_type: OpenTag, value: "<"},
		Token{token_type: ForwardSlash, value: "/"},
		Token{token_type: Attribute, value: "h3"},
		Token{token_type: CloseTag, value: ">"},

		Token{token_type: OpenTag, value: "<"},
		Token{token_type: ForwardSlash, value: "/"},
		Token{token_type: Attribute, value: "div2"},
		Token{token_type: CloseTag, value: ">"},


		Token{token_type: OpenTag, value: "<"},
		Token{token_type: Attribute, value: "h1"},
		Token{token_type: CloseTag, value: ">"},
		Token{token_type: Text, value: "Parse Mofucker"},
		Token{token_type: OpenTag, value: "<"},
		Token{token_type: ForwardSlash, value: "/"},
		Token{token_type: Attribute, value: "h1"},
		Token{token_type: CloseTag, value: ">"},


		Token{token_type: OpenTag, value: "<"},
		Token{token_type: Attribute, value: "span"},	
		Token{token_type: Attribute, value: "id"},
		Token{token_type: Assign, value: "="},
		Token{token_type: Quote, value: "\""},
		Token{token_type: Value, value: "this-is-an-id"},
		Token{token_type: Quote, value: "\""},
		Token{token_type: CloseTag, value: ">"},
		
		Token{token_type: OpenTag, value: "<"},
		Token{token_type: Attribute, value: "h2"},
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
		Token{token_type: Attribute, value: "h2"},
		Token{token_type: CloseTag, value: ">"},

		Token{token_type: OpenTag, value: "<"},
		Token{token_type: ForwardSlash, value: "/"},
		Token{token_type: Attribute, value: "span"},
		Token{token_type: CloseTag, value: ">"},
		
		Token{token_type: OpenTag, value: "<"},
		Token{token_type: ForwardSlash, value: "/"},
		Token{token_type: Attribute, value: "div"},
		Token{token_type: CloseTag, value: ">"},
	}

	_  = Node{
		tag: "div",
		children: []*Node{
			&Node{
				tag: "span",
				attributes: []NodeAttribute{
					NodeAttribute{
						name: "id",
						values: []string{"this-is-an-id",},
					},
				},
				children: []*Node{
					&Node{
						tag: "h1",
						attributes: []NodeAttribute{
							NodeAttribute{
								name: "id",
								values: []string{"one",},
							},
							NodeAttribute{
								name: "class",
								values: []string{"one", "two", "three",},
							},
						},
					},
				},
			},
		},
	}
	
	//WalkNode(&node)

	root := Parser(test_case)
	//fmt.Println(root)
	WalkNode(root)
}
