package query

import (
	// "fmt"
	"github.com/vishen/go-jang/parser"
	"github.com/vishen/go-jang/tokenizer"
	"testing"
)

func getTokens(html string) *parser.Node {
	return parser.Parser(tokenizer.Tokenizer(html))
}

func TestFindByAttribute(t *testing.T) {

	html := `
<html>
	<body>
		<div id="home" class="one">
			<p class="two" href> Text 1</p>
			<p class="three"> Text 2 </p>
		</div>
		<div href="p"></div>
	</body>
</html>

`

	// Using tokenzier because I am lazy and cbf building Node objects
	// at the moment.
	// TODO(vishen): Write out nodes by hand.
	root := getTokens(html)
	// fmt.Println(root)
	// parser.WalkNode(root)
	root_query := NewQueryFromNode(root)
	// fmt.Println(query.Nodes)
	// query = query.A("")
	// query = query.Ae("class", "three")

	query1 := root_query.FindByTag("div").FindChildrenByAttributeEquals("class", "two")

	if len(query1.Nodes) != 1 || query1.Nodes[0].Tag != "p" || query1.Nodes[0].Attributes[0].Name != "class" || query1.Nodes[0].Attributes[0].Values[0] != "two" {
		t.Error("Query didn't match properly")
	}

	query2 := root_query.FindByTag("div").FindChildrenByAttributeEquals("class", "one")

	if len(query2.Nodes) != 1 || query2.Nodes[0].Tag != "div" || query2.Nodes[0].Attributes[0].Name != "id" || query2.Nodes[0].Attributes[0].Values[0] != "home" || query2.Nodes[0].Attributes[1].Name != "class" || query2.Nodes[0].Attributes[1].Values[0] != "one" {
		t.Error("Query didn't match properly")
	}

}
