package main

import (
	//` "encoding/json"
	"fmt"
	"github.com/vishen/go-jang/parser"
	"github.com/vishen/go-jang/query"
	"os"
)

func printNodes(nodes []*parser.Node) {
	for _, node := range nodes {
		fmt.Println(node)
		//bytes, _ := json.MarshalIndent(&node, "", " ")
		//bytes, _ := json.Marshal(&node)
		//fmt.Println(string(bytes))
	}
}

func main() {
	url := os.Args[1]
	lookup := os.Args[2]
	fmt.Println(url, lookup)

	// q, err := query.GetQueryFromUrl(url)

	// if err != nil {
	// 	panic(err)
	// }

	// printNodes(q.T("div").Nodes)

	q, err := query.GetGQueryFromUrl(url, lookup)

	if err != nil {
		panic(err)
	}

	printNodes(q.Nodes)

	//_ = json.NewEncoder(os.Stdout).Encode(&q.Nodes)

}
