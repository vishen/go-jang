package main

import (
	"fmt"
	"github.com/vishen/go-jang/parser"
	"github.com/vishen/go-jang/request"
	"os"
)

func findHrefs(node *parser.Node) {
	for _, attr := range node.Attributes {
		if attr.Name == "href" {
			fmt.Printf("%s=%s\n", attr.Name, attr.Values)
		}
	}

	for _, child := range node.Children {
		findHrefs(child)
	}
}

func main() {
	url := os.Args[1]
	fmt.Println(url)
	root, _ := request.Get(url)
	parser.WalkNode(root)

	fmt.Println("###############################")
	findHrefs(root)

}
