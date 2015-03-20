package main

import (
	"strings"
	"fmt"
)

type NodeAttribute struct {
	name   string
	values []string
}

type Node struct {
	tag        string
	attributes []NodeAttribute

	text string

	parent   *Node
	children []*Node
}


func WalkNode(root *Node) {
	fmt.Printf("New Node: %s, %s, %v, %v, self: %p parent: %p\n", root.tag, root.text, root.attributes, root.children, &root, root.parent)
	//fmt.Printf("New Node: %s, self: %p, parent: %p\n",root.tag, root, root.parent) 
	fmt.Println(len(root.children))
	for _, child := range root.children {
		//fmt.Printf("New Node: %s, %s, %v\n", child.tag, child.text, child.attributes)
		WalkNode(child)
	}
	fmt.Println("End Tag: ", root.tag)
}

func Parser(tokens []Token) *Node {
	var root_node *Node
	var current_node *Node
	var previous_node *Node
	var attribute_values []string
	
	//all_nodes := []*Node{}

	var current_nodeattr NodeAttribute
	var current_token Token

	max_depth := 0
	current_depth := 0

	pos := 0

	tokens_length := len(tokens)
	for {
		if pos >= tokens_length {
			break
		}
		current_token = tokens[pos]
		switch current_token.token_type {
		case OpenTag:
			if tokens[pos+1].token_type == ForwardSlash {
				// Need to close of current node
				pos += 3
				current_depth -= 1
				if current_node.parent != nil {
					previous_node = current_node.parent
					previous_node.children = append(previous_node.children, current_node)
					current_node = previous_node
				}
			} else {
				current_node = &Node{attributes: []NodeAttribute{}, parent: current_node, children: make([]*Node, 0)}
				//all_nodes = append(all_nodes, current_node)
				if root_node == nil {
					root_node = current_node
				} else {
					current_depth += 1
					max_depth += 1
				}
			}
		case CloseTag:
			previous_node = current_node
		case Attribute:
			if current_node.tag == "" {
				current_node.tag = current_token.value
			} else {
				current_nodeattr = NodeAttribute{name: current_token.value, values: []string{}}
			}
		case Value:
			attribute_values = strings.Split(current_token.value, " ")
			current_nodeattr.values = attribute_values

			current_node.attributes = append(current_node.attributes, current_nodeattr)
		case Text:
			previous_node.text = current_token.value

		}

		pos += 1

	}
	return root_node
}
