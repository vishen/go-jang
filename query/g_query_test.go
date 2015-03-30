package query

import (
	"fmt"
	"testing"
)

func createToken(token_type GQueryTokenType, value string) GQueryToken {
	return GQueryToken{token_type: token_type, value: value}
}

func TestTokenizer(t *testing.T) {

	test_cases := []struct {
		lookup string
		wanted []GQueryToken
	}{
		{
			"^h1",
			[]GQueryToken{
				createToken(TAG, "h1"),
			},
		},
		{
			"#hello .value data-href=hello ^div",
			[]GQueryToken{
				createToken(ID, "#"),
				createToken(VALUE, "hello"),
				createToken(SPACE, " "),
				createToken(CLASS, "."),
				createToken(VALUE, "value"),
				createToken(SPACE, " "),
				createToken(ATTRIBUTE, "data-href"),
				createToken(EQUALS, "="),
				createToken(VALUE, "hello"),
				createToken(SPACE, " "),
				createToken(TAG, "div"),
			},
		},
	}

	for _, test_case := range test_cases {
		tokens := tokenzier(test_case.lookup)

		actual := test_case.wanted

		for i, token := range tokens {
			if token.token_type != actual[i].token_type || token.value != actual[i].value {
				fmt.Println("Tokens don't match.")
				fmt.Println("Got:", token)
				fmt.Println("Wanted:", actual[i])

				t.Error("Tokens Mismatch.")
				break
			}
		}
	}

}
