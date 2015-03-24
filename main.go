package main

import "fmt"

func findHrefs(node *Node) {
	for _, attr := range node.attributes {
		if attr.name == "href" {
			fmt.Printf("%s=%s\n", attr.name, attr.values)
		}
	}

	for _, child := range node.children {
		findHrefs(child)
	}
}

func main() {
	url := "http://google.com"
	root, _ := Get(url)
	WalkNode(root)

	fmt.Println("###############################")
	findHrefs(root)

}
