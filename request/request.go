package request

import (
	"errors"
	_ "fmt"
	"io/ioutil"
	"net/http"

	"github.com/vishen/go-jang/parser"
	"github.com/vishen/go-jang/tokenizer"
)

func Get(url string) (*parser.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	// fmt.Println(string(contents))
	tokens := tokenizer.Tokenizer(string(contents))

	if len(tokens) == 0 {
		return nil, errors.New("Response returned no Tokens.")
	}

	// for i, token := range tokens {
	// 	fmt.Println(i, token)
	// for _, ch := range token.Value {
	// 	fmt.Println(ch)
	// }
	// fmt.Println(token.Value)

	// if i == 50 {
	// 	break
	// }
	// }

	root := parser.Parser(tokens)

	if root == nil {
		return nil, errors.New("Tokens resulted in no Nodes")
	}

	return root, nil
}
