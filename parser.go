package main

import (
	"fmt"
	"strings"
)

const (
	START_NODE_TAG = "html"
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
	var attribute_values []string

	nodes := []*Node{}

	var current_nodeattr NodeAttribute

	var current_token Token
	//var previous_token Token
	waiting := true
	current_depth := 0

	pos := 0

	tokens_length := len(tokens)
	for {
		if pos >= tokens_length {
			break
		}
		//previous_token = current_token
		current_token = tokens[pos]

		if waiting {
			if current_token.token_type == Tag && current_token.value == START_NODE_TAG {
				waiting = false
				pos -= 2
			}

		} else {

			switch current_token.token_type {
			case OpenTag:
				if tokens[pos+1].token_type == ForwardSlash {
					// Need to close of current node
					// pos += 3
					// current_depth -= 1
					// if current_node.parent != nil {
					// 	previous_node = current_node.parent
					// 	previous_node.children = append(previous_node.children, current_node)
					// 	// fmt.Printf("Adding child %s ###### %s\n", current_node.tag, previous_node.tag)
					// 	fmt.Printf("%v [%p]\n", previous_node, previous_node)
					// 	current_node = previous_node

					// }

					close_tag := tokens[pos+2].value
					next_children := []*Node{}
					i := len(nodes) - 1
					for ; i > 0; i-- {
						if nodes[i].tag == close_tag {
							nodes[i].children = append(nodes[i].children, next_children...)
							break
						}

						next_children = append(next_children, nodes[i])
					}
					nodes = nodes[0:i]
					pos += 3
					current_depth -= 1

				} else {
					current_node = &Node{attributes: []NodeAttribute{}, children: make([]*Node, 0)}
					root_node = current_node
					nodes = append(nodes, current_node)
				}
			case CloseTag:
				// previous_node = current_node
			case Tag:
				current_node.tag = current_token.value
			case Attribute:
				//fmt.Println(current_token.value)
				current_nodeattr = NodeAttribute{name: current_token.value, values: []string{}}
			case Value:
				attribute_values = strings.Split(current_token.value, " ")
				current_nodeattr.values = attribute_values

				current_node.attributes = append(current_node.attributes, current_nodeattr)
			case Text:
				current_node.text = current_token.value

			}
		}

		pos += 1

	}

	for i, node := range nodes {
		fmt.Printf("%d, %s\n", i, node)
	}

	return root_node
}
