package main

import (
	"io/ioutil"
	"net/http"
)

func Get(url string) (*Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	tokens := Tokenizer(string(contents))

	return Parser(tokens), nil
}
