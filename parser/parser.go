package parser

import (
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

// TODO(vishen): Maybe add the depth of a node?
type Node struct {
	Id         int
	Tag        string
	Attributes []NodeAttribute

	Text string

	Parent   *Node
	Children []*Node
}

// TODO(vishen): Maybe need to move these into a seperate directory
// and allow them to take a list of nodes, so we can potentially chain them.
func (self Node) FindTag(tag string) []*Node {
	found_nodes := []*Node{}

	if self.Tag == tag {
		found_nodes = append(found_nodes, &self)
	}

	for _, child := range self.Children {
		found_nodes = append(found_nodes, child.FindTag(tag)...)
	}

	return found_nodes
}

func (self Node) FindAttribute(attr_name string) []*Node {
	found_nodes := []*Node{}

	for _, attr := range self.Attributes {
		if attr.Name == attr_name {
			found_nodes = append(found_nodes, &self)
			break
		}
	}

	for _, child := range self.Children {
		found_nodes = append(found_nodes, child.FindAttribute(attr_name)...)
	}

	return found_nodes
}

func (self Node) FindAttributeEquals(attr_name, equals string) []*Node {
	found_nodes := []*Node{}

	for _, attr := range self.Attributes {
		if attr.Name == attr_name {
			for _, value := range attr.Values {
				if value == equals {
					found_nodes = append(found_nodes, &self)
					break
				}
			}

		}
	}

	for _, child := range self.Children {
		found_nodes = append(found_nodes, child.FindAttributeEquals(attr_name, equals)...)
	}

	return found_nodes
}

func WalkNode(root *Node) {
	fmt.Printf("New Node: %s, %s, %v, %v, self: %p parent: %p\n", root.Tag, root.Text, root.Attributes, root.Children, &root, root.Parent)
	//fmt.Printf("New Node: %s, self: %p, parent: %p\n",root.Tag, root, root.parent)
	// fmt.Println(len(root.Children))
	for _, child := range root.Children {
		//fmt.Printf("New Node: %s, %s, %v\n", child.Tag, child.text, child.attributes)
		WalkNode(child)
	}
	fmt.Println("End Tag: ", root.Tag)
}

func Parser(tokens []tokenizer.Token) *Node {
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
					current_node = &Node{Id: current_id, Attributes: []NodeAttribute{}, Children: make([]*Node, 0)}
					current_id++
					if root_node == nil {
						root_node = current_node
					}
					nodes = append(nodes, current_node)
				}
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
