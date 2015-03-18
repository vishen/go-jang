package main

import (
	"fmt"
	"strings"
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
	childern []Node
}

func Parser(tokens []Token) *Node {
	var root_node *Node
	var previous_node *Node
	var current_node *Node

	var attribute_values []string

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
			if tokens[pos+1] == ForwardSlash {
				// Need to close of current node
				current_node = previous_node
				previous_node = previous_node.parent
				pos += 3
			} else {
				current_node = &Node{attributes: []NodeAttribute{}, parent: &current_node, children: []Node{}}
				if root_node == nil {
					root_node = current_node
				} else {
					current_depth += 1
					max_depth += 1
				}
			}
		case CloseTag:
			if previous_node != nil {
				previous_node.children = append(previous_node.children, current_node)
			}

			previous_node = &current_node
		case Attribute:
			if current_node.tag == "" {
				current_node.tag = current_token.value
			} else {
				current_nodeattr = NodeAttr{name: current_token.value, values: []string{}}
			}
		case Value:
			attribute_values = strings.Split(current_token.value, " ", -1)
			current_nodeattr.values = attribute_values

			current_node.attributes = append(current_node.attributes, current_nodeattr)
		case Text:
			current_node.text = current_token.value

		}

	}

	return root_node
}
