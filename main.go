package main

import (
	"fmt"
	"github.com/vishen/go-jang/parser"
	"github.com/vishen/go-jang/request"
	"os"
)

func printNodes(nodes []*parser.Node) {
	for _, node := range nodes {
		fmt.Println(node)
	}
}

func main() {
	url := os.Args[1]
	fmt.Println(url)
	root, _ := request.Get(url)
	// parser.WalkNode(root)

	// root.FindAndPrintAttributes("href")

	printNodes(root.FindTag("div"))

}
