package main

import (
	"testing"
	"github.com/graphql-go/graphql"
	"encoding/json"
)

func TestGql(t *testing.T) {
	cases := []struct {
		query        string
		expectedData string
	}{
		{"{hello}", `{"hello":"world"}`},
		{`{
			login(user: "james", password: "abc123") {
				token
			}
		`, `{"login":{"token":"fake-token"}}`,
		},
		{"{}", `null`},
		{
			query: `{ hello bar obj { name subobj { foo } } }`,
			expectedData: `{"bar":"baz","hello":"world","obj":{"name":"James","subobj":{"foo":"bar"}}}`,
		},
	}

	for _, c := range cases {
		params := graphql.Params{Schema: schema, RequestString: c.query}
		r := graphql.Do(params)
		rJSON, err := json.Marshal(r.Data)
		if err != nil {
			t.Fatal("request failed", err)
		}
		if string(rJSON) != c.expectedData {
			t.Fatalf("expected '%s', got '%s'\n", c.expectedData, string(rJSON))
		}
	}
}
