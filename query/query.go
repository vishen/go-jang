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
	nodes := []*parser.Node{node}

	return &Query{Nodes: nodes}
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
		found_nodes = append(found_nodes, findAttribute(attr_name, node)...)

	}

	return &Query{Nodes: found_nodes}

}

func (self Query) FindByAttributeEquals(attr_name, equals string) *Query {

	found_nodes := []*parser.Node{}

	for _, node := range self.Nodes {
		found_nodes = append(found_nodes, findAttributeEquals(attr_name, equals, node)...)

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

func findAttribute(attr_name string, node *parser.Node) []*parser.Node {
	found_nodes := []*parser.Node{}

	if node.Attributes != nil {
		for _, attr := range node.Attributes {
			if attr.Name == attr_name {
				found_nodes = append(found_nodes, node)
				break
			}
		}
	}

	for _, child := range node.Children {
		found_nodes = append(found_nodes, findAttribute(attr_name, child)...)
	}

	return found_nodes
}

func findAttributeEquals(attr_name, equals string, node *parser.Node) []*parser.Node {
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

	for _, child := range node.Children {
		found_nodes = append(found_nodes, findAttributeEquals(attr_name, equals, child)...)
	}

	return found_nodes
}

/*
	Wrapper functions to allow for easier chaining of methods.
*/
func (self Query) T(tag string) *Query {
	return self.FindByTag(tag)
}

func (self Query) A(attr_name string) *Query {
	return self.FindByAttribute(attr_name)
}

func (self Query) Ae(attr_name, equals string) *Query {
	return self.FindByAttributeEquals(attr_name, equals)
}
