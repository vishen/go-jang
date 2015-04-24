package parser

import (
	"encoding/json"
	"fmt"
	"github.com/vishen/go-jang/tokenizer"
	"strings"
)

const (
	START_NODE_TAG = "html"
)

type NodeAttribute struct {
	Name   string
	Values []string
}

func (self NodeAttribute) String() string {
	return fmt.Sprintf("%s=\"%s\"", self.Name, strings.Join(self.Values, " "))
}

// TODO(vishen): Maybe add the depth of a node?
type Node struct {
	Id         int
	Tag        string
	Attributes []NodeAttribute

	Column int
	Line   int

	Text string

	Parent   *Node
	Children []*Node
}

func (self Node) String() string {
	children_ids := []int{}

	for _, child := range self.Children {
		children_ids = append(children_ids, child.Id)
	}

	return fmt.Sprintf("[%d] %s (%d:%d), %d children%v, parent[%d], %s, %s", self.Id, self.Tag, self.Line,
		self.Column, len(children_ids), children_ids, self.Parent.Id, self.Attributes, self.Text)
}

func (self Node) AllChildren() []*Node {
	children := []*Node{}

	children = append(children, self.Children...)

	for _, child := range self.Children {
		children = append(children, child.AllChildren()...)

	}

	return children
}

func (self Node) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marhsalling %q\n", self)
	children := []int{}
	for _, child := range self.Children {
		children = append(children, child.Id)
	}

	return json.Marshal(&struct {
		Id         int
		Tag        string
		Attributes []NodeAttribute

		Column int
		Line   int

		Text     string
		Children []int
	}{
		Id:         self.Id,
		Tag:        self.Tag,
		Attributes: self.Attributes,
		Column:     self.Column,
		Line:       self.Line,
		Text:       self.Text,
		Children:   children,
	})
}

func WalkNode(root *Node) {
	// Only used for debugging at the moment.
	// Provides information about the structure of a node.
	fmt.Printf("New Node: %s, %s, %v, %v, self: %p parent: %p\n", root.Tag, root.Text, root.Attributes, root.Children, &root, root.Parent)
	for _, child := range root.Children {
		WalkNode(child)
	}
	fmt.Println("End Tag: ", root.Tag)
}

func Parser(tokens []tokenizer.Token) *Node {

	/*for i, token := range tokens {
		fmt.Println(i, token)

		if i > 10 {
			return &Node{}
		}
	}*/
	var root_node *Node
	var current_node *Node
	var attribute_values []string

	nodes := []*Node{}

	var current_nodeattr NodeAttribute

	var current_token tokenizer.Token
	//var previous_token Token
	waiting := true
	current_depth := 0

	pos := 0
	current_id := 1

	tokens_length := len(tokens)

	for {
		if pos >= tokens_length {
			break
		}
		//previous_token = current_token
		current_token = tokens[pos]

		if waiting {
			if current_token.Token_type == tokenizer.Tag && current_token.Value == START_NODE_TAG {
				waiting = false
				pos -= 2
			}

		} else {
			switch current_token.Token_type {
			case tokenizer.OpenTag:
				if tokens[pos+1].Token_type == tokenizer.ForwardSlash {
					close_tag := tokens[pos+2].Value
					next_children := []*Node{}
					found := false
					i := len(nodes) - 1
					for ; i >= 0; i-- {
						if nodes[i].Tag == close_tag {
							// Need to add the children in reverse order
							for j := len(next_children) - 1; j >= 0; j-- {
								next_children[j].Parent = nodes[i]
								nodes[i].Children = append(nodes[i].Children, next_children[j])
							}
							found = true
							break
						}

						next_children = append(next_children, nodes[i])
					}
					if found {
						nodes = nodes[0 : i+1]
					}
					pos += 3
					current_depth -= 1

				} else {
					current_node = &Node{Id: current_id,
						Attributes: []NodeAttribute{},
						Children:   make([]*Node, 0),
						Column:     current_token.Column,
						Line:       current_token.Line,
					}
					current_id++
					if root_node == nil {
						root_node = current_node
					}
					nodes = append(nodes, current_node)
				}
			case tokenizer.Comment:
				nodes = append(nodes, &Node{Id: current_id, Tag: "comment", Text: current_token.Value, Column: current_token.Column, Line: current_token.Line})
				current_id++
			case tokenizer.CloseTag:
				// previous_node = current_node
			case tokenizer.Tag:
				current_node.Tag = current_token.Value
			case tokenizer.Attribute:
				//fmt.Println(current_token.value)
				current_nodeattr = NodeAttribute{Name: current_token.Value, Values: []string{}}
			case tokenizer.Value:
				attribute_values = strings.Split(current_token.Value, " ")
				current_nodeattr.Values = attribute_values
				current_node.Attributes = append(current_node.Attributes, current_nodeattr)

			case tokenizer.Text:
				current_node.Text = current_token.Value

			}
		}

		pos += 1

	}

	return root_node
}
