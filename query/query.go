package query

import (
	// "fmt"
	"github.com/vishen/go-jang/parser"
	"github.com/vishen/go-jang/request"
)

func GetQueryFromUrl(url string) (*Query, error) {
	root, err := request.Get(url)

	if err != nil {
		return nil, err
	}

	return NewQueryFromNode(root), nil
}

type Query struct {
	Nodes []*parser.Node
}

func NewQueryFromNode(node *parser.Node) *Query {
	nodes := []*parser.Node{}
	nodes = append(nodes, node)
	return &Query{Nodes: nodes}
}

func (self *Query) Reset() {
	self.Nodes = []*parser.Node{}
}

func (self *Query) AddNodes(nodes []*parser.Node) {
	self.Nodes = append(self.Nodes, nodes...)
	self.clean()
}

/*
	It is possible that there could be the same node
	more than once in a `Query` - like when combining
	2 queries together. Just remove all duplicate nodes
*/
func (self *Query) clean() {
	cleaned := []*parser.Node{}
	var can_add bool
	for _, node := range self.Nodes {
		can_add = true
		for _, cleaned_node := range cleaned {
			if cleaned_node.Id == node.Id {
				can_add = false
				break
			}
		}

		if can_add {
			cleaned = append(cleaned, node)
		}
	}

	self.Nodes = cleaned

}

func (self Query) FindByTag(tag string) *Query {

	found_nodes := []*parser.Node{}

	for _, node := range self.Nodes {
		found_nodes = append(found_nodes, findTag(tag, node)...)

	}

	return &Query{Nodes: found_nodes}

}

func (self Query) FindByAttribute(attr_name string) *Query {

	found_nodes := []*parser.Node{}

	for _, node := range self.Nodes {
		found_nodes = append(found_nodes, findAttribute(attr_name, node, false)...)

	}

	return &Query{Nodes: found_nodes}

}

func (self Query) FindChildrenByAttribute(attr_name string) *Query {

	found_nodes := []*parser.Node{}

	for _, node := range self.Nodes {
		found_nodes = append(found_nodes, findAttribute(attr_name, node, true)...)

	}

	return &Query{Nodes: found_nodes}

}

func (self Query) FindByAttributeEquals(attr_name, equals string) *Query {

	found_nodes := []*parser.Node{}

	for _, node := range self.Nodes {
		found_nodes = append(found_nodes, findAttributeEquals(attr_name, equals, node, false)...)

	}

	return &Query{Nodes: found_nodes}

}

func (self Query) FindChildrenByAttributeEquals(attr_name, equals string) *Query {

	found_nodes := []*parser.Node{}

	for _, node := range self.Nodes {
		found_nodes = append(found_nodes, findAttributeEquals(attr_name, equals, node, true)...)

	}

	return &Query{Nodes: found_nodes}

}

/*
	private helper functions
*/
func findTag(tag string, node *parser.Node) []*parser.Node {
	found_nodes := []*parser.Node{}

	if node.Tag == tag {
		found_nodes = append(found_nodes, node)
	}

	for _, child := range node.Children {
		found_nodes = append(found_nodes, findTag(tag, child)...)
	}

	return found_nodes
}

func findAttribute(attr_name string, node *parser.Node, search_children bool) []*parser.Node {
	found_nodes := []*parser.Node{}

	if node.Attributes != nil {
		for _, attr := range node.Attributes {
			if attr.Name == attr_name {
				found_nodes = append(found_nodes, node)
				break
			}
		}
	}

	if search_children {
		for _, child := range node.Children {
			found_nodes = append(found_nodes, findAttribute(attr_name, child, search_children)...)
		}
	}
	return found_nodes
}

func findAttributeEquals(attr_name, equals string, node *parser.Node, search_children bool) []*parser.Node {
	found_nodes := []*parser.Node{}

	for _, attr := range node.Attributes {
		if attr.Name == attr_name {
			for _, value := range attr.Values {
				if value == equals {
					found_nodes = append(found_nodes, node)
					break
				}
			}

		}
	}

	if search_children {
		for _, child := range node.Children {
			found_nodes = append(found_nodes, findAttributeEquals(attr_name, equals, child, search_children)...)
		}
	}

	return found_nodes
}

/*
	Wrapper functions to allow for easier chaining of methods.
*/
// func (self Query) T(tag string) *Query {
// 	return self.FindByTag(tag)
// }

// func (self Query) A(attr_name string) *Query {
// 	return self.FindChildrenByAttribute(attr_name)
// }

// func (self Query) Ae(attr_name, equals string) *Query {
// 	return self.FindChildrenByAttributeEquals(attr_name, equals)
// }
